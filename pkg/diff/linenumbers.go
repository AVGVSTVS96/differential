package diff

import (
	"fmt"
	"strconv"
	"strings"
)

// LineNumberStyle defines how line numbers are formatted
type LineNumberStyle int

const (
	// LineNumberStandard shows full line numbers with padding (default)
	// Format: "    4     5" (12 chars width)
	LineNumberStandard LineNumberStyle = iota
	
	// LineNumberCompact shows condensed line numbers
	// Format: "4 5" (minimal width)
	LineNumberCompact
	
	// LineNumberMinimal shows no line numbers, only markers
	// Format: "" (0 width, just the marker)
	LineNumberMinimal
	
	// LineNumberLeft shows only left (old) line numbers
	// Format: "    4" (6 chars width)
	LineNumberLeft
	
	// LineNumberRight shows only right (new) line numbers
	// Format: "    5" (6 chars width)
	LineNumberRight
)

// String returns the string representation of the line number style
func (s LineNumberStyle) String() string {
	switch s {
	case LineNumberStandard:
		return "standard"
	case LineNumberCompact:
		return "compact"
	case LineNumberMinimal:
		return "minimal"
	case LineNumberLeft:
		return "left"
	case LineNumberRight:
		return "right"
	default:
		return "standard"
	}
}

// ParseLineNumberStyle parses a string into a LineNumberStyle
func ParseLineNumberStyle(s string) LineNumberStyle {
	switch strings.ToLower(s) {
	case "compact":
		return LineNumberCompact
	case "minimal", "none":
		return LineNumberMinimal
	case "left", "old":
		return LineNumberLeft
	case "right", "new":
		return LineNumberRight
	default:
		return LineNumberStandard
	}
}

// FormatLineNumbers formats line numbers according to the specified style
func FormatLineNumbers(oldLine, newLine int, style LineNumberStyle, lineType LineType) string {
	switch style {
	case LineNumberCompact:
		return formatCompact(oldLine, newLine, lineType)
	case LineNumberMinimal:
		return ""
	case LineNumberLeft:
		return formatLeft(oldLine, lineType)
	case LineNumberRight:
		return formatRight(newLine, lineType)
	default: // LineNumberStandard
		return formatStandard(oldLine, newLine, lineType)
	}
}

// formatStandard formats line numbers in standard style (12 chars total)
func formatStandard(oldLine, newLine int, lineType LineType) string {
	switch lineType {
	case LineRemoved:
		if oldLine > 0 {
			return fmt.Sprintf("%6d       ", oldLine)
		}
		return "            "
	case LineAdded:
		if newLine > 0 {
			return fmt.Sprintf("      %6d", newLine)
		}
		return "            "
	case LineContext:
		if oldLine > 0 && newLine > 0 {
			return fmt.Sprintf("%6d %6d", oldLine, newLine)
		}
		return "            "
	default:
		return "            "
	}
}

// formatCompact formats line numbers in compact style (minimal width)
func formatCompact(oldLine, newLine int, lineType LineType) string {
	switch lineType {
	case LineRemoved:
		if oldLine > 0 {
			return fmt.Sprintf("%d ", oldLine)
		}
		return ""
	case LineAdded:
		if newLine > 0 {
			return fmt.Sprintf("%d", newLine)
		}
		return ""
	case LineContext:
		if oldLine > 0 && newLine > 0 {
			// Use consistent width based on line numbers
			oldWidth := len(strconv.Itoa(oldLine))
			newWidth := len(strconv.Itoa(newLine))
			maxWidth := max(oldWidth, newWidth)
			return fmt.Sprintf("%*d %*d", maxWidth, oldLine, maxWidth, newLine)
		}
		return ""
	default:
		return ""
	}
}

// formatLeft formats only old (left) line numbers (6 chars width)
func formatLeft(oldLine int, lineType LineType) string {
	if lineType == LineAdded {
		return "      " // Empty space for added lines
	}
	if oldLine > 0 {
		return fmt.Sprintf("%6d", oldLine)
	}
	return "      "
}

// formatRight formats only new (right) line numbers (6 chars width)
func formatRight(newLine int, lineType LineType) string {
	if lineType == LineRemoved {
		return "      " // Empty space for removed lines
	}
	if newLine > 0 {
		return fmt.Sprintf("%6d", newLine)
	}
	return "      "
}

// GetLineNumberWidth returns the display width for the given style
func GetLineNumberWidth(style LineNumberStyle) int {
	switch style {
	case LineNumberStandard:
		return 12 // "    4     5 "
	case LineNumberCompact:
		return 0 // Variable width, handled dynamically
	case LineNumberMinimal:
		return 0 // No line numbers
	case LineNumberLeft, LineNumberRight:
		return 6 // "    4 " or "    5 "
	default:
		return 12
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
