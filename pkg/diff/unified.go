package diff

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss/v2/compat"
	"github.com/charmbracelet/x/ansi"

	"github.com/avgvstvs96/differential/pkg/config"
	"github.com/avgvstvs96/differential/pkg/styles"
	"github.com/avgvstvs96/differential/pkg/theme"
	"github.com/avgvstvs96/differential/pkg/util"
)

// renderUnifiedLine renders a single line in unified diff format
func renderUnifiedLine(fileName string, dl DiffLine, width int, t theme.Theme, cfg *config.Config) string {
	removedLineStyle, addedLineStyle, contextLineStyle, lineNumberStyle := createStyles(t)

	// Determine line style and marker based on line type
	var marker string
	var bgStyle styles.Style
	var lineNum string
	var highlightColor compat.AdaptiveColor

	switch dl.Kind {
	case LineRemoved:
		marker = "-"
		bgStyle = removedLineStyle
		lineNumberStyle = lineNumberStyle.Background(t.DiffRemovedLineNumberBg()).Foreground(t.DiffRemoved())
		highlightColor = t.DiffHighlightRemoved() // TODO: handle "none"
		if dl.OldLineNo > 0 {
			lineNum = fmt.Sprintf("%6d       ", dl.OldLineNo)
		} else {
			lineNum = "            "
		}
	case LineAdded:
		marker = "+"
		bgStyle = addedLineStyle
		lineNumberStyle = lineNumberStyle.Background(t.DiffAddedLineNumberBg()).Foreground(t.DiffAdded())
		highlightColor = t.DiffHighlightAdded() // TODO: handle "none"
		if dl.NewLineNo > 0 {
			lineNum = fmt.Sprintf("      %7d", dl.NewLineNo)
		} else {
			lineNum = "            "
		}
	case LineContext:
		marker = " "
		bgStyle = contextLineStyle
		if dl.OldLineNo > 0 && dl.NewLineNo > 0 {
			lineNum = fmt.Sprintf("%6d %6d", dl.OldLineNo, dl.NewLineNo)
		} else {
			lineNum = "            "
		}
	}

	// Create the line prefix
	prefix := renderLinePrefix(dl, lineNum, marker, lineNumberStyle, t)

	// Render the content
	prefixWidth := ansi.StringWidth(prefix)
	contentWidth := width - prefixWidth
	content := renderLineContent(fileName, dl, bgStyle, highlightColor, contentWidth, cfg)

	return prefix + content
}

// RenderUnifiedHunk formats a hunk for unified display
func RenderUnifiedHunk(fileName string, h Hunk, cfg *config.Config, opts ...UnifiedOption) string {
	// Apply options to create the configuration
	renderConfig := NewUnifiedConfig(opts...)

	var sb strings.Builder

	// Render hunk header first
	if h.Header != "" {
		header := RenderHunkHeader(h, cfg, renderConfig.Width)
		sb.WriteString(header)
		sb.WriteString("\n")
	}

	// Make a copy of the hunk so we don't modify the original
	hunkCopy := Hunk{Lines: make([]DiffLine, len(h.Lines))}
	copy(hunkCopy.Lines, h.Lines)

	// Highlight changes within lines
	HighlightIntralineChanges(&hunkCopy)

	sb.Grow(len(hunkCopy.Lines) * renderConfig.Width)

	util.WriteStringsPar(&sb, hunkCopy.Lines, func(line DiffLine) string {
		return renderUnifiedLine(fileName, line, renderConfig.Width, theme.CurrentTheme(), cfg) + "\n"
	})

	return sb.String()
}

// FormatUnifiedDiff creates a unified formatted view of a diff
func FormatUnifiedDiff(filename string, diffText string, cfg *config.Config, opts ...UnifiedOption) (string, error) {
	diffResult, err := ParseUnifiedDiff(diffText)
	if err != nil {
		return "", err
	}

	var sb strings.Builder

	// Add file header if enabled
	if cfg != nil && cfg.Display.ShowFileHeaders {
		sb.WriteString(RenderFileHeader(diffResult, cfg.Display.ShowStats, cfg.Features.Hyperlinks, cfg.Display.GraphicalStats))
	}

	util.WriteStringsPar(&sb, diffResult.Hunks, func(h Hunk) string {
		return RenderUnifiedHunk(filename, h, cfg, opts...)
	})

	return sb.String(), nil
}
