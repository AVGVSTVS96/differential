package diff

import (
	"image/color"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/lipgloss/v2/compat"
	"github.com/charmbracelet/x/ansi"

	"github.com/avgvstvs96/differential/pkg/config"
	"github.com/avgvstvs96/differential/pkg/styles"
	"github.com/avgvstvs96/differential/pkg/textutil"
	"github.com/avgvstvs96/differential/pkg/theme"
)

// createStyles generates the lipgloss styles needed for rendering diffs
func createStyles(t theme.Theme) (removedLineStyle, addedLineStyle, contextLineStyle, lineNumberStyle styles.Style) {
	removedLineStyle = styles.NewStyle().Background(t.DiffRemovedBg())
	addedLineStyle = styles.NewStyle().Background(t.DiffAddedBg())
	// Use terminal's default background for context lines (no background set = transparent)
	contextLineStyle = styles.NewStyle()
	// Use transparent background for line numbers with just muted text color
	lineNumberStyle = styles.NewStyle().Foreground(t.TextMuted())
	return
}

// renderLinePrefix renders the line number and marker prefix for a diff line
func renderLinePrefix(dl DiffLine, lineNum string, marker string, lineNumberStyle styles.Style, t theme.Theme) string {
	// Style the marker based on line type
	var styledMarker string
	switch dl.Kind {
	case LineRemoved:
		styledMarker = styles.NewStyle().Foreground(t.DiffRemoved()).Background(t.DiffRemovedBg()).Render(marker)
	case LineAdded:
		styledMarker = styles.NewStyle().Foreground(t.DiffAdded()).Background(t.DiffAddedBg()).Render(marker)
	case LineContext:
		styledMarker = styles.NewStyle().Foreground(t.TextMuted()).Render(marker)
	default:
		styledMarker = marker
	}

	return lineNumberStyle.Render(lineNum + " " + styledMarker)
}

// renderLineContent renders the content of a diff line with syntax and intra-line highlighting
func renderLineContent(fileName string, dl DiffLine, bgStyle styles.Style, highlightColor compat.AdaptiveColor, width int, cfg *config.Config) string {
	// Process content (tab expansion, truncation)
	content := dl.Content
	if cfg != nil {
		content = textutil.ProcessLineContent(content, cfg)
	}

	// Apply syntax highlighting
	// For context lines, use NoColor to prevent syntax highlighter from adding background
	bg := bgStyle.GetBackground()
	if _, ok := bg.(lipgloss.NoColor); ok {
		bg = lipgloss.NoColor{}
	}
	content = highlightLine(fileName, content, bg)

	// Apply intra-line highlighting if needed
	if len(dl.Segments) > 0 && (dl.Kind == LineRemoved || dl.Kind == LineAdded) {
		// Get foreground color based on line type
		t := theme.CurrentTheme()
		var fgColor compat.AdaptiveColor
		if dl.Kind == LineRemoved {
			fgColor = t.BackgroundPanel()
		} else {
			fgColor = t.BackgroundPanel()
		}
		content = applyHighlighting(content, dl.Segments, dl.Kind, highlightColor, fgColor)
	}

	// Add a padding space for added/removed lines
	if dl.Kind == LineRemoved || dl.Kind == LineAdded {
		content = bgStyle.Render(" ") + content
	}

	// Create the final line and truncate if needed
	return bgStyle.MaxHeight(1).Width(width).Render(
		ansi.Truncate(
			content,
			width,
			"...",
		),
	)
}

// adaptiveColorToColor converts an AdaptiveColor to a color.Color
func adaptiveColorToColor(ac compat.AdaptiveColor) color.Color {
	colorStr := styles.AdaptiveColorToString(ac)
	if colorStr == nil {
		return nil
	}
	return lipgloss.Color(*colorStr)
}
