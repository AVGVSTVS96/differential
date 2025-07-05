package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	UI          UIConfig          `toml:"ui"`
	Git         GitConfig         `toml:"git"`
	Keybindings KeybindingsConfig `toml:"keybindings"`
}

type UIConfig struct {
	Theme        string `toml:"theme"`
	DefaultView  string `toml:"default_view"`
	TabWidth     int    `toml:"tab_width"`
	LineNumbers  bool   `toml:"line_numbers"`
	SyntaxHighlight bool `toml:"syntax_highlight"`
	WrapLines    bool   `toml:"wrap_lines"`
}

type GitConfig struct {
	DefaultContext   int  `toml:"default_context"`
	IgnoreWhitespace bool `toml:"ignore_whitespace"`
	ShowStats        bool `toml:"show_stats"`
}

type KeybindingsConfig struct {
	Quit           string `toml:"quit"`
	Help           string `toml:"help"`
	ToggleView     string `toml:"toggle_view"`
	NextHunk       string `toml:"next_hunk"`
	PrevHunk       string `toml:"prev_hunk"`
	ScrollUp       string `toml:"scroll_up"`
	ScrollDown     string `toml:"scroll_down"`
	PageUp         string `toml:"page_up"`
	PageDown       string `toml:"page_down"`
	Search         string `toml:"search"`
	StageHunk      string `toml:"stage_hunk"`
	RefreshDiff    string `toml:"refresh_diff"`
	ToggleNumbers  string `toml:"toggle_numbers"`
}

func NewConfig() *Config {
	return &Config{
		UI: UIConfig{
			Theme:           "dracula",
			DefaultView:     "unified",
			TabWidth:        4,
			LineNumbers:     true,
			SyntaxHighlight: true,
			WrapLines:       false,
		},
		Git: GitConfig{
			DefaultContext:   3,
			IgnoreWhitespace: false,
			ShowStats:        true,
		},
		Keybindings: KeybindingsConfig{
			Quit:          "q",
			Help:          "?",
			ToggleView:    "tab",
			NextHunk:      "}",
			PrevHunk:      "{",
			ScrollUp:      "k",
			ScrollDown:    "j",
			PageUp:        "ctrl+b",
			PageDown:      "ctrl+f",
			Search:        "/",
			StageHunk:     "s",
			RefreshDiff:   "r",
			ToggleNumbers: "n",
		},
	}
}

func (c *Config) ConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "differential", "config.toml")
}