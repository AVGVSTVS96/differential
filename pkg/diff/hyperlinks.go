package diff

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CreateFileHyperlink creates an OSC 8 hyperlink for a file path
// Format: ESC]8;;file://absolute/path ESC\text ESC]8;; ESC\
func CreateFileHyperlink(path string, text string, enabled bool) string {
	if !enabled {
		return text
	}

	// Convert to absolute path for file:// URLs
	absPath := path
	if !filepath.IsAbs(path) {
		cwd, err := os.Getwd()
		if err == nil {
			absPath = filepath.Join(cwd, path)
		}
	}

	// Clean the path
	absPath = filepath.Clean(absPath)

	// Create OSC 8 hyperlink
	// Format: \x1b]8;;URL\x1b\\TEXT\x1b]8;;\x1b\\
	url := fmt.Sprintf("file://%s", absPath)
	return fmt.Sprintf("\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\", url, text)
}

// SupportsHyperlinks attempts to detect if the terminal supports OSC 8 hyperlinks
func SupportsHyperlinks() bool {
	// Check for known terminals that support hyperlinks
	termProgram := os.Getenv("TERM_PROGRAM")
	
	// iTerm2 supports hyperlinks
	if termProgram == "iTerm.app" {
		return true
	}

	// VS Code terminal supports hyperlinks
	if termProgram == "vscode" {
		return true
	}

	// WezTerm supports hyperlinks
	if termProgram == "WezTerm" {
		return true
	}

	// Kitty terminal supports hyperlinks
	if strings.Contains(os.Getenv("TERM"), "kitty") {
		return true
	}

	// VTE-based terminals (GNOME Terminal, etc.) version 0.50+
	vteVersion := os.Getenv("VTE_VERSION")
	if vteVersion != "" {
		// VTE_VERSION is an integer like "5000" for 0.50.0
		// Versions >= 5000 support hyperlinks
		if len(vteVersion) >= 4 && vteVersion >= "5000" {
			return true
		}
	}

	// Alacritty supports hyperlinks
	if strings.Contains(os.Getenv("TERM"), "alacritty") {
		return true
	}

	// Hyper terminal supports hyperlinks
	if termProgram == "Hyper" {
		return true
	}

	// Windows Terminal supports hyperlinks
	if os.Getenv("WT_SESSION") != "" {
		return true
	}

	// Default to false if we can't detect support
	return false
}
