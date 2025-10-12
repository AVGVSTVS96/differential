// Package config provides configuration management for differential.
//
// Configuration is loaded from multiple sources with the following precedence:
//   1. Built-in defaults (lowest priority)
//   2. Git config ([differential] and [delta] sections)
//   3. TOML config file (~/.config/differential/config.toml)
//   4. Custom config file (via --config flag)
//   5. CLI flags (highest priority)
//
// Usage:
//   cfg, err := config.Load(customPath, disableGitConfig)
//   if err != nil { ... }
//
// The package supports full delta compatibility by reading [delta] sections
// from git config, making differential a drop-in replacement.
package config

// Config represents the complete configuration for differential
type Config struct {
	Core     CoreConfig     `toml:"core"`
	Display  DisplayConfig  `toml:"display"`
	Colors   ColorsConfig   `toml:"colors"`
	Style    StyleConfig    `toml:"style"`
	Features FeaturesConfig `toml:"features"`
}

// CoreConfig contains core application settings
type CoreConfig struct {
	Theme       string `toml:"theme"`
	SideBySide  bool   `toml:"side-by-side"`
	LineNumbers bool   `toml:"line-numbers"`
	Pager       string `toml:"pager"`
	NoPager     bool   `toml:"no-pager"`
}

// DisplayConfig contains display and formatting settings
type DisplayConfig struct {
	Width                int    `toml:"width"`
	Tabs                 int    `toml:"tabs"`
	MaxLineLength        int    `toml:"max-line-length"`
	Unified              int    `toml:"unified"`
	WrapMaxLines         string `toml:"wrap-max-lines"`
	WrapLeftSymbol       string `toml:"wrap-left-symbol"`
	WrapRightSymbol      string `toml:"wrap-right-symbol"`
	RelativePaths        bool   `toml:"relative-paths"`
	KeepPlusMinusMarkers bool   `toml:"keep-plus-minus-markers"`
	IgnoreAllSpace       bool   `toml:"ignore-all-space"`
	IgnoreSpaceChange    bool   `toml:"ignore-space-change"`
	IgnoreBlankLines     bool   `toml:"ignore-blank-lines"`
	ShowFileHeaders      bool   `toml:"show-file-headers"`  // Show file status headers
	ShowStats            bool   `toml:"show-stats"`         // Show diff statistics
	GraphicalStats       bool   `toml:"graphical-stats"`    // Show graphical statistics bars
	LineNumbersStyle     string `toml:"line-numbers-style"` // Line number style: standard, compact, minimal
	FileSeparator        string `toml:"file-separator"`     // Separator style between files
	HunkSeparator        string `toml:"hunk-separator"`     // Separator style between hunks
}

// ColorsConfig contains color and syntax highlighting settings
type ColorsConfig struct {
	SyntaxTheme string `toml:"syntax-theme"`
	ColorMode   string `toml:"color"` // auto, always, never
}

// StyleConfig contains styling options for different elements
type StyleConfig struct {
	MinusStyle         string `toml:"minus-style"`
	MinusEmphStyle     string `toml:"minus-emph-style"`
	PlusStyle          string `toml:"plus-style"`
	PlusEmphStyle      string `toml:"plus-emph-style"`
	ZeroStyle          string `toml:"zero-style"`
	CommitStyle        string `toml:"commit-style"`
	FileStyle          string `toml:"file-style"`
	HunkHeaderStyle    string `toml:"hunk-header-style"`
}

// FeaturesConfig contains feature flags for navigation and interactivity
type FeaturesConfig struct {
	Navigate   bool `toml:"navigate"`   // Enable hunk navigation (interactive mode)
	Hyperlinks bool `toml:"hyperlinks"` // Enable clickable file path hyperlinks
}
