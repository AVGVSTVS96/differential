package diff

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/x/ansi"

	"github.com/avgvstvs96/differential/pkg/config"
	"github.com/avgvstvs96/differential/pkg/styles"
	"github.com/avgvstvs96/differential/pkg/theme"
)

// HunkInfo contains parsed information from a hunk header
type HunkInfo struct {
	OldStart int    // Starting line in old file
	OldCount int    // Number of lines in old file
	NewStart int    // Starting line in new file
	NewCount int    // Number of lines in new file
	Context  string // Optional function/context info
}

var hunkHeaderRegex = regexp.MustCompile(`^@@ -(\d+)(?:,(\d+))? \+(\d+)(?:,(\d+))? @@(.*)$`)

// ParseHunkHeader extracts line range information from a hunk header
func ParseHunkHeader(header string) (*HunkInfo, error) {
	matches := hunkHeaderRegex.FindStringSubmatch(header)
	if matches == nil {
		return nil, fmt.Errorf("invalid hunk header format: %s", header)
	}

	info := &HunkInfo{
		OldCount: 1, // Default to 1 if not specified
		NewCount: 1,
	}

	var err error
	info.OldStart, err = strconv.Atoi(matches[1])
	if err != nil {
		return nil, fmt.Errorf("invalid old start line: %w", err)
	}

	if matches[2] != "" {
		info.OldCount, err = strconv.Atoi(matches[2])
		if err != nil {
			return nil, fmt.Errorf("invalid old count: %w", err)
		}
	}

	info.NewStart, err = strconv.Atoi(matches[3])
	if err != nil {
		return nil, fmt.Errorf("invalid new start line: %w", err)
	}

	if matches[4] != "" {
		info.NewCount, err = strconv.Atoi(matches[4])
		if err != nil {
			return nil, fmt.Errorf("invalid new count: %w", err)
		}
	}

	// Context is everything after the "@@" markers (typically function name)
	if len(matches) > 5 {
		info.Context = strings.TrimSpace(matches[5])
	}

	return info, nil
}

// RenderHunkHeader renders a styled hunk header line
func RenderHunkHeader(h Hunk, cfg *config.Config, width int) string {
	if h.Header == "" {
		return ""
	}

	t := theme.CurrentTheme()
	
	// Parse the header to extract information
	info, err := ParseHunkHeader(h.Header)
	if err != nil {
		// Fallback to plain rendering if parsing fails
		return renderPlainHunkHeader(h.Header, t, width)
	}

	var sb strings.Builder

	// Hunk header style: muted background with distinct color
	headerStyle := styles.NewStyle().
		Foreground(t.TextMuted()).
		Background(t.DiffContextBg())

	// Build the header text
	oldRange := fmt.Sprintf("-%d,%d", info.OldStart, info.OldCount)
	newRange := fmt.Sprintf("+%d,%d", info.NewStart, info.NewCount)

	// Style the ranges differently
	oldStyle := styles.NewStyle().
		Foreground(t.DiffRemoved()).
		Background(t.DiffContextBg()).
		Bold(true)
	
	newStyle := styles.NewStyle().
		Foreground(t.DiffAdded()).
		Background(t.DiffContextBg()).
		Bold(true)

	sb.WriteString(headerStyle.Render("@@ "))
	sb.WriteString(oldStyle.Render(oldRange))
	sb.WriteString(headerStyle.Render(" "))
	sb.WriteString(newStyle.Render(newRange))
	sb.WriteString(headerStyle.Render(" @@"))

	// Add context/function name if present
	if info.Context != "" {
		contextStyle := styles.NewStyle().
			Foreground(t.Text()).
			Background(t.DiffContextBg()).
			Italic(true)
		sb.WriteString(headerStyle.Render(" "))
		sb.WriteString(contextStyle.Render(info.Context))
	}

	// Pad to full width with background color
	rendered := sb.String()
	currentWidth := ansi.StringWidth(rendered)
	if currentWidth < width {
		padding := width - currentWidth
		sb.WriteString(headerStyle.Render(strings.Repeat(" ", padding)))
	}

	return sb.String()
}

// renderPlainHunkHeader is a fallback renderer for unparseable headers
func renderPlainHunkHeader(header string, t theme.Theme, width int) string {
	headerStyle := styles.NewStyle().
		Foreground(t.TextMuted()).
		Background(t.DiffContextBg())

	// Pad to full width
	currentWidth := len(header)
	if currentWidth < width {
		padding := width - currentWidth
		header = header + strings.Repeat(" ", padding)
	}

	return headerStyle.Render(header)
}
