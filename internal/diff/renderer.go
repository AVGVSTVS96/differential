package diff

import (
	"fmt"
	"strings"
	"sync"

	"github.com/charmbracelet/lipgloss"
	"github.com/avgvstvs96/differential/internal/themes"
)

// RenderUnifiedDiff renders a diff in unified format with syntax highlighting
func RenderUnifiedDiff(result *DiffResult, opts RenderOptions) string {
	if result.IsBinary {
		return fmt.Sprintf("Binary files %s and %s differ\n", result.OldFile, result.NewFile)
	}

	// Initialize themes if not already done
	themes.Initialize()
	theme := themes.GetCurrentTheme()

	var sb strings.Builder

	// Apply intra-line highlighting to all hunks
	for i := range result.Hunks {
		HighlightIntralineChanges(&result.Hunks[i])
	}

	// Render each hunk
	for _, hunk := range result.Hunks {
		sb.WriteString(renderUnifiedHunk(result.NewFile, hunk, theme, opts))
		sb.WriteString("\n")
	}

	return sb.String()
}

// renderUnifiedHunk renders a single hunk in unified format
func renderUnifiedHunk(filename string, hunk Hunk, theme *themes.ThemeColors, opts RenderOptions) string {
	var sb strings.Builder

	// Render hunk header
	headerStyle := lipgloss.NewStyle().
		Foreground(theme.TextMuted).
		Bold(true)
	sb.WriteString(headerStyle.Render(hunk.Header))
	sb.WriteString("\n")

	// Render lines in parallel for performance
	lines := make([]string, len(hunk.Lines))
	var wg sync.WaitGroup
	wg.Add(len(hunk.Lines))

	for i, line := range hunk.Lines {
		go func(idx int, dl DiffLine) {
			defer wg.Done()
			lines[idx] = renderUnifiedLine(filename, dl, theme, opts)
		}(i, line)
	}

	wg.Wait()

	// Join lines
	for _, line := range lines {
		sb.WriteString(line)
		sb.WriteString("\n")
	}

	return sb.String()
}

// renderUnifiedLine renders a single line in unified format
func renderUnifiedLine(filename string, dl DiffLine, theme *themes.ThemeColors, opts RenderOptions) string {
	var marker string
	var bgStyle lipgloss.Style
	var lineNumberStyle lipgloss.Style
	var highlightColor lipgloss.Color
	var lineNum string

	switch dl.Kind {
	case LineRemoved:
		marker = "-"
		bgStyle = lipgloss.NewStyle().Background(theme.DiffRemovedBg)
		lineNumberStyle = lipgloss.NewStyle().
			Background(theme.DiffRemovedLineNumberBg).
			Foreground(theme.DiffRemoved)
		highlightColor = theme.DiffHighlightRemoved
		if opts.ShowLineNumbers {
			lineNum = fmt.Sprintf("%6d       ", dl.OldLineNo)
		}

	case LineAdded:
		marker = "+"
		bgStyle = lipgloss.NewStyle().Background(theme.DiffAddedBg)
		lineNumberStyle = lipgloss.NewStyle().
			Background(theme.DiffAddedLineNumberBg).
			Foreground(theme.DiffAdded)
		highlightColor = theme.DiffHighlightAdded
		if opts.ShowLineNumbers {
			lineNum = fmt.Sprintf("       %6d", dl.NewLineNo)
		}

	case LineContext:
		marker = " "
		bgStyle = lipgloss.NewStyle().Background(theme.DiffContextBg)
		lineNumberStyle = lipgloss.NewStyle().
			Background(theme.DiffLineNumber).
			Foreground(theme.TextMuted)
		if opts.ShowLineNumbers {
			lineNum = fmt.Sprintf("%6d %6d", dl.OldLineNo, dl.NewLineNo)
		}
	}

	// Build the line
	var result strings.Builder

	// Line numbers
	if opts.ShowLineNumbers {
		result.WriteString(lineNumberStyle.Render(lineNum))
		result.WriteString(" ")
	}

	// Marker
	markerStyle := lipgloss.NewStyle().
		Background(bgStyle.GetBackground()).
		Foreground(bgStyle.GetForeground()).
		Bold(true)
	result.WriteString(markerStyle.Render(marker))

	// Content with syntax highlighting
	content := dl.Content

	// Apply syntax highlighting
	if filename != "" && dl.Kind == LineContext {
		// Only apply syntax highlighting to context lines
		// (added/removed lines will have diff colors)
		highlighted := themes.SyntaxHighlightLine(content, filename)
		content = highlighted
	}

	// Apply intra-line highlighting for added/removed lines
	if len(dl.Segments) > 0 && highlightColor != "" {
		// Create highlight style
		r, g, b := hexToRGB(string(highlightColor))
		highlightStyle := fmt.Sprintf("\x1b[48;2;%d;%d;%dm", r, g, b)
		content = ApplyHighlighting(content, dl.Segments, dl.Kind, highlightStyle)
	}

	// Apply background color to the entire line
	styledContent := bgStyle.Render(content)
	result.WriteString(styledContent)

	// Pad to width if needed
	if opts.Width > 0 {
		currentWidth := VisibleLength(result.String())
		if currentWidth < opts.Width {
			padding := strings.Repeat(" ", opts.Width-currentWidth)
			result.WriteString(bgStyle.Render(padding))
		}
	}

	return result.String()
}

// RenderSideBySideDiff renders a diff in side-by-side format
func RenderSideBySideDiff(result *DiffResult, opts RenderOptions) string {
	if result.IsBinary {
		return fmt.Sprintf("Binary files %s and %s differ\n", result.OldFile, result.NewFile)
	}

	// Initialize themes
	themes.Initialize()
	theme := themes.GetCurrentTheme()

	var sb strings.Builder

	// Apply intra-line highlighting
	for i := range result.Hunks {
		HighlightIntralineChanges(&result.Hunks[i])
	}

	// Calculate column widths
	halfWidth := opts.Width / 2
	if halfWidth < 40 {
		halfWidth = 40
	}

	// Render each hunk
	for _, hunk := range result.Hunks {
		sb.WriteString(renderSideBySideHunk(result.OldFile, result.NewFile, hunk, theme, opts, halfWidth))
		sb.WriteString("\n")
	}

	return sb.String()
}

// renderSideBySideHunk renders a single hunk in side-by-side format
func renderSideBySideHunk(oldFile, newFile string, hunk Hunk, theme *themes.ThemeColors, opts RenderOptions, halfWidth int) string {
	var sb strings.Builder

	// Render hunk header
	headerStyle := lipgloss.NewStyle().
		Foreground(theme.TextMuted).
		Bold(true)
	sb.WriteString(headerStyle.Render(hunk.Header))
	sb.WriteString("\n")

	// Pair lines for side-by-side rendering
	pairs := PairLines(hunk.Lines)

	// Render each pair
	for _, pair := range pairs {
		leftLine := renderSideBySideLine(oldFile, pair.Left, theme, opts, halfWidth, true)
		rightLine := renderSideBySideLine(newFile, pair.Right, theme, opts, halfWidth, false)

		sb.WriteString(leftLine)
		sb.WriteString(" ┃ ")
		sb.WriteString(rightLine)
		sb.WriteString("\n")
	}

	return sb.String()
}

// renderSideBySideLine renders a single line for side-by-side view
func renderSideBySideLine(filename string, dl *DiffLine, theme *themes.ThemeColors, opts RenderOptions, width int, isLeft bool) string {
	if dl == nil {
		// Empty side
		emptyStyle := lipgloss.NewStyle().Background(theme.Background)
		return emptyStyle.Render(strings.Repeat(" ", width))
	}

	// Similar to renderUnifiedLine but adapted for side-by-side
	var bgStyle lipgloss.Style
	var lineNumberStyle lipgloss.Style
	var highlightColor lipgloss.Color
	var lineNum string

	switch dl.Kind {
	case LineRemoved:
		bgStyle = lipgloss.NewStyle().Background(theme.DiffRemovedBg)
		lineNumberStyle = lipgloss.NewStyle().
			Background(theme.DiffRemovedLineNumberBg).
			Foreground(theme.DiffRemoved)
		highlightColor = theme.DiffHighlightRemoved
		if opts.ShowLineNumbers {
			lineNum = fmt.Sprintf("%6d", dl.OldLineNo)
		}

	case LineAdded:
		bgStyle = lipgloss.NewStyle().Background(theme.DiffAddedBg)
		lineNumberStyle = lipgloss.NewStyle().
			Background(theme.DiffAddedLineNumberBg).
			Foreground(theme.DiffAdded)
		highlightColor = theme.DiffHighlightAdded
		if opts.ShowLineNumbers {
			lineNum = fmt.Sprintf("%6d", dl.NewLineNo)
		}

	case LineContext:
		bgStyle = lipgloss.NewStyle().Background(theme.DiffContextBg)
		lineNumberStyle = lipgloss.NewStyle().
			Background(theme.DiffLineNumber).
			Foreground(theme.TextMuted)
		if opts.ShowLineNumbers {
			if isLeft {
				lineNum = fmt.Sprintf("%6d", dl.OldLineNo)
			} else {
				lineNum = fmt.Sprintf("%6d", dl.NewLineNo)
			}
		}
	}

	var result strings.Builder

	// Line numbers
	if opts.ShowLineNumbers {
		result.WriteString(lineNumberStyle.Render(lineNum))
		result.WriteString(" ")
	}

	// Content
	content := dl.Content

	// Apply syntax highlighting for context lines
	if filename != "" && dl.Kind == LineContext {
		content = themes.SyntaxHighlightLine(content, filename)
	}

	// Apply intra-line highlighting
	if len(dl.Segments) > 0 && highlightColor != "" {
		r, g, b := hexToRGB(string(highlightColor))
		highlightStyle := fmt.Sprintf("\x1b[48;2;%d;%d;%dm", r, g, b)
		content = ApplyHighlighting(content, dl.Segments, dl.Kind, highlightStyle)
	}

	// Truncate if needed
	contentWidth := width
	if opts.ShowLineNumbers {
		contentWidth -= 7 // Line number width
	}
	content = TruncateString(content, contentWidth)

	// Apply background and add to result
	styledContent := bgStyle.Render(content)
	result.WriteString(styledContent)

	// Pad to width
	currentWidth := VisibleLength(result.String())
	if currentWidth < width {
		padding := strings.Repeat(" ", width-currentWidth)
		result.WriteString(bgStyle.Render(padding))
	}

	return result.String()
}

// hexToRGB converts a hex color to RGB values
func hexToRGB(hex string) (r, g, b int) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return 255, 255, 255
	}

	fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	return
}

// FormatUnifiedDiff formats an entire diff in unified view
func FormatUnifiedDiff(filename, diffText string, opts RenderOptions) (string, error) {
	result, err := ParseUnifiedDiff(diffText)
	if err != nil {
		return "", err
	}

	// Use filename from diff if not provided
	if filename == "" && result.NewFile != "" {
		filename = result.NewFile
	}

	return RenderUnifiedDiff(result, opts), nil
}

// FormatSideBySideDiff formats an entire diff in side-by-side view
func FormatSideBySideDiff(filename, diffText string, opts RenderOptions) (string, error) {
	result, err := ParseUnifiedDiff(diffText)
	if err != nil {
		return "", err
	}

	// Use filename from diff if not provided
	if filename == "" && result.NewFile != "" {
		filename = result.NewFile
	}

	return RenderSideBySideDiff(result, opts), nil
}