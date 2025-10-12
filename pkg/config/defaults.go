package config

import (
	"os"
	"path/filepath"
)

// DefaultConfig returns a Config with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		Core: CoreConfig{
			Theme:       "opencode",
			SideBySide:  false,
			LineNumbers: false,
			Pager:       "",
			NoPager:     false,
		},
		Display: DisplayConfig{
			Width:                0, // auto-detect
			Tabs:                 4,
			MaxLineLength:        3000,
			Unified:              3, // default context lines
			WrapMaxLines:         "2",
			WrapLeftSymbol:       "↵",
			WrapRightSymbol:      "↵",
			RelativePaths:        false,
			KeepPlusMinusMarkers: false,
			ShowFileHeaders:      true,
			ShowStats:            true,
			GraphicalStats:       false,
			LineNumbersStyle:     "standard",
			FileSeparator:        "none",
			HunkSeparator:        "none",
		},
		Colors: ColorsConfig{
			SyntaxTheme: "auto",
			ColorMode:   "auto",
		},
		Style: StyleConfig{
			MinusStyle:      "normal auto",
			MinusEmphStyle:  "normal auto",
			PlusStyle:       "syntax auto",
			PlusEmphStyle:   "syntax auto",
			ZeroStyle:       "normal",
			CommitStyle:     "raw",
			FileStyle:       "blue",
			HunkHeaderStyle: "line-number syntax",
		},
		Features: FeaturesConfig{
			Navigate:   false, // Interactive navigation (Phase 4B)
			Hyperlinks: true,  // Enable clickable file paths by default
		},
	}
}

// DefaultConfigPath returns the default configuration file path
func DefaultConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	
	// Check XDG_CONFIG_HOME first
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		return filepath.Join(xdgConfig, "differential", "config.toml")
	}
	
	// Fallback to ~/.config/differential/config.toml
	return filepath.Join(homeDir, ".config", "differential", "config.toml")
}

// GitConfigPath returns the git config file path
func GitConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(homeDir, ".gitconfig")
}

// DeltaToConfigKeyMap maps delta config keys to differential config structure
var DeltaToConfigKeyMap = map[string]string{
	"theme":                    "core.theme",
	"side-by-side":             "core.side-by-side",
	"line-numbers":             "core.line-numbers",
	"pager":                    "core.pager",
	"width":                    "display.width",
	"tabs":                     "display.tabs",
	"max-line-length":          "display.max-line-length",
	"syntax-theme":             "colors.syntax-theme",
	"dark":                     "colors.dark",
	"light":                    "colors.light",
	"navigate":                 "features.navigate",
	"hyperlinks":               "features.hyperlinks",
	"relative-paths":           "display.relative-paths",
	"keep-plus-minus-markers":  "display.keep-plus-minus-markers",
	"minus-style":              "style.minus-style",
	"minus-emph-style":         "style.minus-emph-style",
	"plus-style":               "style.plus-style",
	"plus-emph-style":          "style.plus-emph-style",
	"zero-style":               "style.zero-style",
	"commit-style":             "style.commit-style",
	"file-style":               "style.file-style",
	"hunk-header-style":        "style.hunk-header-style",
}
