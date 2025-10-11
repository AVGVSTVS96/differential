package config

// Config represents the complete configuration for differential
type Config struct {
	Core    CoreConfig    `toml:"core"`
	Display DisplayConfig `toml:"display"`
	Colors  ColorsConfig  `toml:"colors"`
	Style   StyleConfig   `toml:"style"`
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
	Width            int    `toml:"width"`
	Tabs             int    `toml:"tabs"`
	MaxLineLength    int    `toml:"max-line-length"`
	Unified          int    `toml:"unified"`
	WrapMaxLines     string `toml:"wrap-max-lines"`
	WrapLeftSymbol   string `toml:"wrap-left-symbol"`
	WrapRightSymbol  string `toml:"wrap-right-symbol"`
	RelativePaths    bool   `toml:"relative-paths"`
	KeepPlusMinusMarkers bool `toml:"keep-plus-minus-markers"`
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
