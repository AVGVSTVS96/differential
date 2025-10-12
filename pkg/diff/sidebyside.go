package diff

import (
	"fmt"
	"strings"
	"sync"

	"github.com/charmbracelet/lipgloss/v2/compat"
	"github.com/charmbracelet/x/ansi"

	"github.com/avgvstvs96/differential/pkg/config"
	"github.com/avgvstvs96/differential/pkg/styles"
	"github.com/avgvstvs96/differential/pkg/theme"
	"github.com/avgvstvs96/differential/pkg/util"
)

// pairLines converts a flat list of diff lines to pairs for side-by-side display
func pairLines(lines []DiffLine) []linePair {
	var pairs []linePair
	i := 0

	for i < len(lines) {
		switch lines[i].Kind {
		case LineRemoved:
			// Check if the next line is an addition, if so pair them
			if i+1 < len(lines) && lines[i+1].Kind == LineAdded {
				pairs = append(pairs, linePair{left: &lines[i], right: &lines[i+1]})
				i += 2
			} else {
				pairs = append(pairs, linePair{left: &lines[i], right: nil})
				i++
			}
		case LineAdded:
			pairs = append(pairs, linePair{left: nil, right: &lines[i]})
			i++
		case LineContext:
			pairs = append(pairs, linePair{left: &lines[i], right: &lines[i]})
			i++
		}
	}

	return pairs
}

// renderDiffColumnLine is a helper function that handles the common logic for rendering diff columns
func renderDiffColumnLine(
	fileName string,
	dl *DiffLine,
	colWidth int,
	isLeftColumn bool,
	t theme.Theme,
	cfg *config.Config,
) string {
	if dl == nil {
		// Use transparent background for empty cells
		contextLineStyle := styles.NewStyle()
		return contextLineStyle.Width(colWidth).Render("")
	}

	removedLineStyle, addedLineStyle, contextLineStyle, lineNumberStyle := createStyles(t)

	// Determine line style based on line type and column
	var marker string
	var bgStyle styles.Style
	var lineNum string
	var highlightColor compat.AdaptiveColor

	if isLeftColumn {
		// Left column logic
		switch dl.Kind {
		case LineRemoved:
			marker = "-"
			bgStyle = removedLineStyle
			lineNumberStyle = lineNumberStyle.Background(t.DiffRemovedLineNumberBg()).Foreground(t.DiffRemoved())
			highlightColor = t.DiffHighlightRemoved() // TODO: handle "none"
		case LineAdded:
			marker = "?"
			bgStyle = contextLineStyle
		case LineContext:
			marker = " "
			bgStyle = contextLineStyle
		}

		// Format line number for left column
		if dl.OldLineNo > 0 {
			lineNum = fmt.Sprintf("%6d", dl.OldLineNo)
		}
	} else {
		// Right column logic
		switch dl.Kind {
		case LineAdded:
			marker = "+"
			bgStyle = addedLineStyle
			lineNumberStyle = lineNumberStyle.Background(t.DiffAddedLineNumberBg()).Foreground(t.DiffAdded())
			highlightColor = t.DiffHighlightAdded()
		case LineRemoved:
			marker = "?"
			bgStyle = contextLineStyle
		case LineContext:
			marker = " "
			bgStyle = contextLineStyle
		}

		// Format line number for right column
		if dl.NewLineNo > 0 {
			lineNum = fmt.Sprintf("%6d", dl.NewLineNo)
		}
	}

	// Create the line prefix
	prefix := renderLinePrefix(*dl, lineNum, marker, lineNumberStyle, t)

	// Determine if we should render content
	shouldRenderContent := (dl.Kind == LineRemoved && isLeftColumn) ||
		(dl.Kind == LineAdded && !isLeftColumn) ||
		dl.Kind == LineContext

	if !shouldRenderContent {
		return bgStyle.Width(colWidth).Render("")
	}

	// Render the content
	prefixWidth := ansi.StringWidth(prefix)
	contentWidth := colWidth - prefixWidth
	content := renderLineContent(fileName, *dl, bgStyle, highlightColor, contentWidth, cfg)

	return prefix + content
}

// renderLeftColumn formats the left side of a side-by-side diff
func renderLeftColumn(fileName string, dl *DiffLine, colWidth int, cfg *config.Config) string {
	return renderDiffColumnLine(fileName, dl, colWidth, true, theme.CurrentTheme(), cfg)
}

// renderRightColumn formats the right side of a side-by-side diff
func renderRightColumn(fileName string, dl *DiffLine, colWidth int, cfg *config.Config) string {
	return renderDiffColumnLine(fileName, dl, colWidth, false, theme.CurrentTheme(), cfg)
}

// RenderSideBySideHunk formats a hunk for side-by-side display
func RenderSideBySideHunk(fileName string, h Hunk, cfg *config.Config, opts ...UnifiedOption) string {
	// Apply options to create the configuration
	renderConfig := NewSideBySideConfig(opts...)

	// Make a copy of the hunk so we don't modify the original
	hunkCopy := Hunk{Lines: make([]DiffLine, len(h.Lines))}
	copy(hunkCopy.Lines, h.Lines)

	// Highlight changes within lines
	HighlightIntralineChanges(&hunkCopy)

	// Pair lines for side-by-side display
	pairs := pairLines(hunkCopy.Lines)

	// Calculate column width
	colWidth := renderConfig.Width / 2

	leftWidth := colWidth
	rightWidth := renderConfig.Width - colWidth
	var sb strings.Builder

	util.WriteStringsPar(&sb, pairs, func(p linePair) string {
		wg := &sync.WaitGroup{}
		var leftStr, rightStr string
		wg.Add(2)
		go func() {
			defer wg.Done()
			leftStr = renderLeftColumn(fileName, p.left, leftWidth, cfg)
		}()
		go func() {
			defer wg.Done()
			rightStr = renderRightColumn(fileName, p.right, rightWidth, cfg)
		}()
		wg.Wait()
		return leftStr + rightStr + "\n"
	})

	return sb.String()
}

// FormatDiff creates a side-by-side formatted view of a diff
func FormatDiff(filename string, diffText string, cfg *config.Config, opts ...UnifiedOption) (string, error) {
	diffResult, err := ParseUnifiedDiff(diffText)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	util.WriteStringsPar(&sb, diffResult.Hunks, func(h Hunk) string {
		return RenderSideBySideHunk(filename, h, cfg, opts...)
	})

	return sb.String(), nil
}
