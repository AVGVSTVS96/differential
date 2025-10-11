package styles

import (
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/lipgloss/v2/compat"
	"github.com/lucasb-eyer/go-colorful"
)

// Style is a wrapper around lipgloss.Style for convenience
type Style = lipgloss.Style

// NewStyle creates a new lipgloss style
func NewStyle() Style {
	return lipgloss.NewStyle()
}

// Terminal holds information about the terminal environment
// This is exposed for backward compatibility but should use theme.Terminal
var Terminal struct {
	BackgroundIsDark bool
}

// AdaptiveColorToString converts a compat.AdaptiveColor to the appropriate
// hex color string based on the current terminal background
func AdaptiveColorToString(color compat.AdaptiveColor) *string {
	// Use the global Terminal variable which is set by the application
	if Terminal.BackgroundIsDark {
		if _, ok := color.Dark.(lipgloss.NoColor); ok {
			return nil
		}
		c1, _ := colorful.MakeColor(color.Dark)
		hex := c1.Hex()
		return &hex
	}
	if _, ok := color.Light.(lipgloss.NoColor); ok {
		return nil
	}
	c1, _ := colorful.MakeColor(color.Light)
	hex := c1.Hex()
	return &hex
}

// Helper functions for style pointers
func BoolPtr(b bool) *bool       { return &b }
func StringPtr(s string) *string { return &s }
func UintPtr(u uint) *uint       { return &u }
