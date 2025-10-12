package config

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

// Loader handles configuration loading and merging
type Loader struct {
	config        *Config
	disableGitConfig bool
}

// NewLoader creates a new configuration loader
func NewLoader() *Loader {
	return &Loader{
		config: DefaultConfig(),
	}
}

// DisableGitConfig disables loading from git config
func (l *Loader) DisableGitConfig() *Loader {
	l.disableGitConfig = true
	return l
}

// LoadFromFile loads configuration from a TOML file
func (l *Loader) LoadFromFile(path string) error {
	if path == "" {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Config file doesn't exist, that's okay
			return nil
		}
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var fileConfig Config
	if err := toml.Unmarshal(data, &fileConfig); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	// Merge file config into default config
	l.mergeConfig(&fileConfig)
	return nil
}

// LoadDefault loads configuration from default locations
func (l *Loader) LoadDefault() error {
	// First load from git config (if not disabled)
	if !l.disableGitConfig {
		// Try to load from [differential] section first
		if err := l.LoadFromGitConfig("differential"); err != nil {
			// Silently ignore errors (git config might not exist)
		}
		
		// Also try to load from [delta] section for compatibility
		if err := l.LoadFromGitConfig("delta"); err != nil {
			// Silently ignore errors
		}
	}
	
	// Then try default differential config file (overrides git config)
	if err := l.LoadFromFile(DefaultConfigPath()); err != nil {
		return err
	}
	
	return nil
}

// LoadFromGitConfig loads configuration from git config
func (l *Loader) LoadFromGitConfig(section string) error {
	// Get all config values for the section
	cmd := exec.Command("git", "config", "--get-regexp", fmt.Sprintf("^%s\\.", section))
	output, err := cmd.Output()
	if err != nil {
		// git config returns exit code 1 if section doesn't exist
		return nil
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			continue
		}
		
		// Remove section prefix (e.g., "delta." or "differential.")
		key := strings.TrimPrefix(parts[0], section+".")
		value := parts[1]
		
		l.applyGitConfigValue(key, value)
	}
	
	return nil
}

// applyGitConfigValue applies a single git config key-value pair
func (l *Loader) applyGitConfigValue(key, value string) {
	// Normalize key (replace underscores with hyphens)
	key = strings.ReplaceAll(key, "_", "-")
	
	// Apply value based on key
	switch key {
	// Core settings
	case "theme":
		l.config.Core.Theme = value
	case "side-by-side":
		l.config.Core.SideBySide = parseBool(value)
	case "line-numbers":
		l.config.Core.LineNumbers = parseBool(value)
	case "pager":
		l.config.Core.Pager = value
	case "no-pager":
		l.config.Core.NoPager = parseBool(value)
		
	// Display settings
	case "width":
		if v, err := parseInt(value); err == nil {
			l.config.Display.Width = v
		}
	case "tabs":
		if v, err := parseInt(value); err == nil {
			l.config.Display.Tabs = v
		}
	case "max-line-length":
		if v, err := parseInt(value); err == nil {
			l.config.Display.MaxLineLength = v
		}
	case "unified":
		if v, err := parseInt(value); err == nil {
			l.config.Display.Unified = v
		}
	case "wrap-max-lines":
		l.config.Display.WrapMaxLines = value
	case "wrap-left-symbol":
		l.config.Display.WrapLeftSymbol = value
	case "wrap-right-symbol":
		l.config.Display.WrapRightSymbol = value
	case "relative-paths":
		l.config.Display.RelativePaths = parseBool(value)
	case "keep-plus-minus-markers":
		l.config.Display.KeepPlusMinusMarkers = parseBool(value)
		
	// Color settings
	case "syntax-theme":
		l.config.Colors.SyntaxTheme = value
	case "color":
		l.config.Colors.ColorMode = value
	case "dark":
		if parseBool(value) {
			l.config.Colors.ColorMode = "dark"
		}
	case "light":
		if parseBool(value) {
			l.config.Colors.ColorMode = "light"
		}
		
	// Style settings
	case "minus-style":
		l.config.Style.MinusStyle = value
	case "minus-emph-style":
		l.config.Style.MinusEmphStyle = value
	case "plus-style":
		l.config.Style.PlusStyle = value
	case "plus-emph-style":
		l.config.Style.PlusEmphStyle = value
	case "zero-style":
		l.config.Style.ZeroStyle = value
	case "commit-style":
		l.config.Style.CommitStyle = value
	case "file-style":
		l.config.Style.FileStyle = value
	case "hunk-header-style":
		l.config.Style.HunkHeaderStyle = value
	}
}

// parseBool parses a git config boolean value
func parseBool(s string) bool {
	s = strings.ToLower(s)
	return s == "true" || s == "yes" || s == "on" || s == "1"
}

// parseInt parses a git config integer value, returning 0 if invalid
// Returns the parsed integer and an error if parsing fails
func parseInt(s string) (int, error) {
	if s == "" {
		return 0, nil
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid integer value %q: %w", s, err)
	}
	return i, nil
}

// mergeConfig merges a loaded config into the current config
// Non-zero/non-empty values from src override defaults in dst
// 
// Note: This implementation treats zero values (false, 0, "") as "not set" which means
// explicit false/0 values in TOML config files cannot override true/non-zero defaults.
// This is acceptable for current use cases since:
// - Git config values are parsed from strings, so "false" becomes parseBool(false) 
// - TOML unmarshalling only sets fields that are explicitly present
// - Default config typically uses false/0 for most options
// 
// If future use cases require distinguishing "not set" from "explicitly false/0",
// consider using pointer types (*bool, *int) or a separate "set" tracking mechanism.
func (l *Loader) mergeConfig(src *Config) {
	// Merge Core
	if src.Core.Theme != "" {
		l.config.Core.Theme = src.Core.Theme
	}
	if src.Core.SideBySide {
		l.config.Core.SideBySide = src.Core.SideBySide
	}
	if src.Core.LineNumbers {
		l.config.Core.LineNumbers = src.Core.LineNumbers
	}
	if src.Core.Pager != "" {
		l.config.Core.Pager = src.Core.Pager
	}
	if src.Core.NoPager {
		l.config.Core.NoPager = src.Core.NoPager
	}

	// Merge Display
	if src.Display.Width != 0 {
		l.config.Display.Width = src.Display.Width
	}
	if src.Display.Tabs != 0 {
		l.config.Display.Tabs = src.Display.Tabs
	}
	if src.Display.MaxLineLength != 0 {
		l.config.Display.MaxLineLength = src.Display.MaxLineLength
	}
	if src.Display.Unified != 0 {
		l.config.Display.Unified = src.Display.Unified
	}
	if src.Display.WrapMaxLines != "" {
		l.config.Display.WrapMaxLines = src.Display.WrapMaxLines
	}
	if src.Display.WrapLeftSymbol != "" {
		l.config.Display.WrapLeftSymbol = src.Display.WrapLeftSymbol
	}
	if src.Display.WrapRightSymbol != "" {
		l.config.Display.WrapRightSymbol = src.Display.WrapRightSymbol
	}
	if src.Display.RelativePaths {
		l.config.Display.RelativePaths = src.Display.RelativePaths
	}
	if src.Display.KeepPlusMinusMarkers {
		l.config.Display.KeepPlusMinusMarkers = src.Display.KeepPlusMinusMarkers
	}

	// Merge Colors
	if src.Colors.SyntaxTheme != "" {
		l.config.Colors.SyntaxTheme = src.Colors.SyntaxTheme
	}
	if src.Colors.ColorMode != "" {
		l.config.Colors.ColorMode = src.Colors.ColorMode
	}

	// Merge Style
	if src.Style.MinusStyle != "" {
		l.config.Style.MinusStyle = src.Style.MinusStyle
	}
	if src.Style.MinusEmphStyle != "" {
		l.config.Style.MinusEmphStyle = src.Style.MinusEmphStyle
	}
	if src.Style.PlusStyle != "" {
		l.config.Style.PlusStyle = src.Style.PlusStyle
	}
	if src.Style.PlusEmphStyle != "" {
		l.config.Style.PlusEmphStyle = src.Style.PlusEmphStyle
	}
	if src.Style.ZeroStyle != "" {
		l.config.Style.ZeroStyle = src.Style.ZeroStyle
	}
	if src.Style.CommitStyle != "" {
		l.config.Style.CommitStyle = src.Style.CommitStyle
	}
	if src.Style.FileStyle != "" {
		l.config.Style.FileStyle = src.Style.FileStyle
	}
	if src.Style.HunkHeaderStyle != "" {
		l.config.Style.HunkHeaderStyle = src.Style.HunkHeaderStyle
	}
}

// Config returns the loaded configuration
func (l *Loader) Config() *Config {
	return l.config
}

// Load loads configuration with the following precedence (lowest to highest):
// 1. Built-in defaults (lowest priority)
// 2. Git config (if not disabled)
// 3. Default config file (~/.config/differential/config.toml)
// 4. Custom config file (if specified)
// 5. CLI flags (highest priority, applied later by caller)
func Load(customConfigPath string, disableGitConfig bool) (*Config, error) {
	loader := NewLoader()
	
	if disableGitConfig {
		loader.DisableGitConfig()
	}

	// Load from custom config path if specified
	if customConfigPath != "" {
		if err := loader.LoadFromFile(customConfigPath); err != nil {
			return nil, err
		}
	} else {
		// Load from default locations
		if err := loader.LoadDefault(); err != nil {
			return nil, err
		}
	}

	return loader.Config(), nil
}
