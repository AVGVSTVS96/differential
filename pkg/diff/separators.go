package diff

import (
	"strings"

	"github.com/avgvstvs96/differential/pkg/styles"
	"github.com/avgvstvs96/differential/pkg/theme"
)

// SeparatorStyle defines the visual style of separators
type SeparatorStyle int

const (
	// SeparatorNone shows no separator
	SeparatorNone SeparatorStyle = iota
	
	// SeparatorBlank shows a blank line
	SeparatorBlank
	
	// SeparatorLine shows a single line (─)
	SeparatorLine
	
	// SeparatorDoubleLine shows a double line (═)
	SeparatorDoubleLine
	
	// SeparatorDashed shows a dashed line (┄)
	SeparatorDashed
	
	// SeparatorDotted shows a dotted line (⋯)
	SeparatorDotted
)

// String returns the string representation of the separator style
func (s SeparatorStyle) String() string {
	switch s {
	case SeparatorNone:
		return "none"
	case SeparatorBlank:
		return "blank"
	case SeparatorLine:
		return "line"
	case SeparatorDoubleLine:
		return "double-line"
	case SeparatorDashed:
		return "dashed"
	case SeparatorDotted:
		return "dotted"
	default:
		return "none"
	}
}

// ParseSeparatorStyle parses a string into a SeparatorStyle
func ParseSeparatorStyle(s string) SeparatorStyle {
	switch strings.ToLower(s) {
	case "blank", "space":
		return SeparatorBlank
	case "line", "single":
		return SeparatorLine
	case "double-line", "double", "thick":
		return SeparatorDoubleLine
	case "dashed", "dash":
		return SeparatorDashed
	case "dotted", "dots":
		return SeparatorDotted
	default:
		return SeparatorNone
	}
}

// RenderSeparator renders a separator line with the specified style and width
func RenderSeparator(style SeparatorStyle, width int, t theme.Theme) string {
	if style == SeparatorNone {
		return ""
	}
	
	if style == SeparatorBlank {
		return "\n"
	}
	
	// Get the separator character
	char := getSeparatorChar(style)
	
	// Create the separator line
	line := strings.Repeat(char, width)
	
	// Style with muted color
	separatorStyle := styles.NewStyle().Foreground(t.TextMuted())
	
	return separatorStyle.Render(line) + "\n"
}

// RenderHunkSeparator renders a separator between hunks
func RenderHunkSeparator(style SeparatorStyle, width int, t theme.Theme) string {
	return RenderSeparator(style, width, t)
}

// RenderFileSeparator renders a separator between files
func RenderFileSeparator(style SeparatorStyle, width int, t theme.Theme) string {
	return RenderSeparator(style, width, t)
}

// getSeparatorChar returns the character to use for the given style
func getSeparatorChar(style SeparatorStyle) string {
	switch style {
	case SeparatorLine:
		return "─"
	case SeparatorDoubleLine:
		return "═"
	case SeparatorDashed:
		return "┄"
	case SeparatorDotted:
		return "⋯"
	default:
		return "─"
	}
}
