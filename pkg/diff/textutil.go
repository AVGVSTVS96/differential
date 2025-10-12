package diff

import (
	"strings"
	"unicode/utf8"

	"github.com/avgvstvs96/differential/pkg/config"
	"github.com/charmbracelet/x/ansi"
)

// ExpandTabs replaces tab characters with spaces
// Properly accounts for wide characters (CJK, emojis) that occupy 2 columns
func ExpandTabs(s string, tabWidth int) string {
	if tabWidth == 0 || !strings.Contains(s, "\t") {
		return s
	}

	var result strings.Builder
	result.Grow(len(s) * 2) // Conservative estimate for wide chars
	col := 0

	for _, r := range s {
		if r == '\t' {
			// Calculate spaces needed to reach next tab stop
			spaces := tabWidth - (col % tabWidth)
			result.WriteString(strings.Repeat(" ", spaces))
			col += spaces
		} else {
			result.WriteRune(r)
			// Account for wide characters (CJK, emojis, etc.)
			// Simple heuristic: Unicode >= 0x1100 is typically wide
			// For more accuracy, consider using github.com/mattn/go-runewidth
			if r >= 0x1100 && (r <= 0x115F || // Hangul Jamo
				r >= 0x2E80 && r <= 0x9FFF || // CJK
				r >= 0xAC00 && r <= 0xD7A3 || // Hangul Syllables
				r >= 0xF900 && r <= 0xFAFF || // CJK Compatibility
				r >= 0xFE10 && r <= 0xFE19 || // Vertical forms
				r >= 0xFE30 && r <= 0xFE6F || // CJK Compatibility Forms
				r >= 0xFF00 && r <= 0xFF60 || // Fullwidth Forms
				r >= 0xFFE0 && r <= 0xFFE6 || // Fullwidth Forms
				r >= 0x1F300 && r <= 0x1F64F || // Emoticons
				r >= 0x1F680 && r <= 0x1F6FF) { // Transport symbols
				col += 2
			} else {
				col++
			}
		}
	}

	return result.String()
}

// TruncateWithEllipsis truncates a string to maxLen, adding ellipsis if needed
func TruncateWithEllipsis(s string, maxLen int, ellipsis string) string {
	if maxLen <= 0 {
		return s
	}

	// Use ansi.StringWidth for proper handling of ANSI sequences
	width := ansi.StringWidth(s)
	if width <= maxLen {
		return s
	}

	// Truncate using ansi.Truncate which handles ANSI properly
	return ansi.Truncate(s, maxLen, ellipsis)
}

// WrapLine wraps a long line into multiple lines
func WrapLine(line string, maxWidth int, wrapSymbol string, maxWraps int) []string {
	if maxWidth <= 0 || maxWraps == 0 {
		return []string{line}
	}

	width := ansi.StringWidth(line)
	if width <= maxWidth {
		return []string{line}
	}

	// For now, simple truncation if wrapping is enabled but limited
	// Full wrapping implementation would be more complex with ANSI handling
	if maxWraps > 0 {
		// Truncate with wrap symbol
		truncated := ansi.Truncate(line, maxWidth-ansi.StringWidth(wrapSymbol), "")
		return []string{truncated + wrapSymbol}
	}

	return []string{line}
}

// ProcessLineContent applies tab expansion, wrapping, and truncation
func ProcessLineContent(content string, cfg *config.Config) string {
	// Expand tabs
	if cfg.Display.Tabs > 0 {
		content = ExpandTabs(content, cfg.Display.Tabs)
	}

	// Truncate if needed
	if cfg.Display.MaxLineLength > 0 {
		content = TruncateWithEllipsis(content, cfg.Display.MaxLineLength, "...")
	}

	return content
}

// StripTrailingWhitespace removes trailing whitespace from a string
func StripTrailingWhitespace(s string) string {
	return strings.TrimRight(s, " \t\r\n")
}

// CountLeadingWhitespace counts leading whitespace characters
func CountLeadingWhitespace(s string) int {
	count := 0
	for _, r := range s {
		if r == ' ' || r == '\t' {
			count++
		} else {
			break
		}
	}
	return count
}

// VisualizeWhitespace replaces whitespace with visible characters
func VisualizeWhitespace(s string, showTabs, showSpaces, showNewlines bool) string {
	if !showTabs && !showSpaces && !showNewlines {
		return s
	}

	var result strings.Builder
	result.Grow(len(s))

	for _, r := range s {
		switch r {
		case '\t':
			if showTabs {
				result.WriteString("→   ")
			} else {
				result.WriteRune(r)
			}
		case ' ':
			if showSpaces {
				result.WriteString("·")
			} else {
				result.WriteRune(r)
			}
		case '\n':
			if showNewlines {
				result.WriteString("↵\n")
			} else {
				result.WriteRune(r)
			}
		default:
			result.WriteRune(r)
		}
	}

	return result.String()
}

// CalculateVisibleLength calculates the visible length of a string, excluding ANSI codes
func CalculateVisibleLength(s string) int {
	return ansi.StringWidth(s)
}

// RemoveANSI strips all ANSI escape sequences from a string
func RemoveANSI(s string) string {
	return ansi.Strip(s)
}

// PadToWidth pads a string to a specific width with spaces
func PadToWidth(s string, width int) string {
	currentWidth := ansi.StringWidth(s)
	if currentWidth >= width {
		return s
	}
	return s + strings.Repeat(" ", width-currentWidth)
}

// IsBlankLine checks if a line contains only whitespace
func IsBlankLine(s string) bool {
	return strings.TrimSpace(s) == ""
}

// NormalizeWhitespace normalizes whitespace in a string for comparison
func NormalizeWhitespace(s string, mode string) string {
	switch mode {
	case "ignore-all":
		// Remove all whitespace
		return strings.Map(func(r rune) rune {
			if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
				return -1
			}
			return r
		}, s)
	case "ignore-space-change":
		// Collapse multiple spaces to single space
		var result strings.Builder
		lastWasSpace := false
		for _, r := range s {
			isSpace := r == ' ' || r == '\t'
			if isSpace {
				if !lastWasSpace {
					result.WriteRune(' ')
					lastWasSpace = true
				}
			} else {
				result.WriteRune(r)
				lastWasSpace = false
			}
		}
		return result.String()
	default:
		return s
	}
}

// SplitByWidth splits a string by display width
// TODO: This function does not currently preserve ANSI escape sequences during splitting.
// ANSI sequences may be split mid-code, which could break formatting. For now, this is
// a simple width-based split. Future enhancement should use proper ANSI-aware splitting.
func SplitByWidth(s string, maxWidth int) []string {
	if maxWidth <= 0 {
		return []string{s}
	}

	width := ansi.StringWidth(s)
	if width <= maxWidth {
		return []string{s}
	}

	// Simple split by width - does not preserve ANSI sequences
	var lines []string
	currentLine := ""
	currentWidth := 0

	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		
		// Calculate display width for this rune
		runeWidth := 1
		if r == '\t' {
			runeWidth = 4 // Tab width assumption
		} else if r >= 0x1100 {
			// Approximate wide character detection (CJK, emojis, etc.)
			// For more accuracy, consider using github.com/mattn/go-runewidth
			runeWidth = 2
		}

		if currentWidth+runeWidth > maxWidth {
			lines = append(lines, currentLine)
			currentLine = string(r)
			currentWidth = runeWidth
		} else {
			currentLine += string(r)
			currentWidth += runeWidth
		}
		i += size
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}
