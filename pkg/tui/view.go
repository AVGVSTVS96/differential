package tui

import (
	"fmt"
	"strings"

	"github.com/avgvstvs96/differential/pkg/diff"
	"github.com/avgvstvs96/differential/pkg/styles"
	"github.com/avgvstvs96/differential/pkg/theme"
)

// View renders the current view of the diff
func (m Model) View() string {
	if m.quitting {
		return ""
	}
	
	if m.showHelp {
		return m.renderHelp()
	}
	
	// Render all diffs if not yet cached
	if len(m.renderedLines) == 0 {
		m.renderAllContent()
	}
	
	// Show file tree in split panel if enabled
	if m.showFileTree {
		return m.renderSplitView()
	}
	
	// Full-width diff view (no file tree)
	return m.renderFullView()
}

// renderFullView renders the diff in full width (no file tree)
func (m Model) renderFullView() string {
	var sb strings.Builder
	
	// Calculate visible window
	startLine := m.scrollOffset
	endLine := m.scrollOffset + m.height - 2 // Reserve space for status line
	
	if endLine > len(m.renderedLines) {
		endLine = len(m.renderedLines)
	}
	if startLine < 0 {
		startLine = 0
	}
	
	// Render visible lines
	for i := startLine; i < endLine; i++ {
		if i < len(m.renderedLines) {
			sb.WriteString(m.renderedLines[i])
			sb.WriteString("\n")
		}
	}
	
	// Add status line at bottom
	sb.WriteString(m.renderStatusLine())
	
	return sb.String()
}

// renderSplitView renders file tree + diff in split panel
func (m Model) renderSplitView() string {
	var sb strings.Builder
	
	// Calculate panel widths
	treeWidth := m.width / 4 // 25% for tree
	if treeWidth < 20 {
		treeWidth = 20
	}
	if treeWidth > 40 {
		treeWidth = 40
	}
	
	diffWidth := m.width - treeWidth - 1 // -1 for separator
	
	// Build file tree
	fileTree := BuildFileTree(m.results, m.currentFile)
	treeLines := strings.Split(RenderFileTree(fileTree, treeWidth, m.height-1), "\n")
	
	// Calculate visible diff window
	startLine := m.scrollOffset
	endLine := m.scrollOffset + m.height - 2
	
	if endLine > len(m.renderedLines) {
		endLine = len(m.renderedLines)
	}
	if startLine < 0 {
		startLine = 0
	}
	
	// Render lines side by side
	for i := 0; i < m.height-1; i++ {
		// File tree column
		if i < len(treeLines) {
			sb.WriteString(treeLines[i])
		} else {
			sb.WriteString(strings.Repeat(" ", treeWidth))
		}
		
		// Separator
		t := theme.CurrentTheme()
		sepStyle := styles.NewStyle().Foreground(t.Border())
		sb.WriteString(sepStyle.Render("│"))
		
		// Diff column
		diffLineIdx := startLine + i
		if diffLineIdx < len(m.renderedLines) {
			diffLine := m.renderedLines[diffLineIdx]
			// Truncate if too wide
			if len(diffLine) > diffWidth {
				diffLine = diffLine[:diffWidth]
			}
			sb.WriteString(diffLine)
		}
		
		sb.WriteString("\n")
	}
	
	// Add status line at bottom
	sb.WriteString(m.renderStatusLine())
	
	return sb.String()
}

// renderAllContent renders all diff content to lines
func (m *Model) renderAllContent() {
	var sb strings.Builder
	
	for i, result := range m.results {
		// Render file header
		if m.cfg.Display.ShowFileHeaders {
			header := diff.RenderFileHeader(result, m.cfg.Display.ShowStats, m.cfg.Features.Hyperlinks, m.cfg.Display.GraphicalStats)
			sb.WriteString(header)
		}
		
		// Render each hunk
		for j, hunk := range result.Hunks {
			totalHunks := len(result.Hunks)
			
			var hunkOutput string
			if m.cfg.Core.SideBySide {
				hunkOutput = diff.RenderSideBySideHunk("", hunk, m.cfg, j, totalHunks, diff.WithWidth(m.width))
			} else {
				hunkOutput = diff.RenderUnifiedHunk("", hunk, m.cfg, j, totalHunks, diff.WithWidth(m.width))
			}
			sb.WriteString(hunkOutput)
		}
		
		// Add file separator if not last file
		if i < len(m.results)-1 {
			sepStyle := diff.ParseSeparatorStyle(m.cfg.Display.FileSeparator)
			if sepStyle != diff.SeparatorNone {
				sb.WriteString(diff.RenderFileSeparator(sepStyle, m.width, theme.CurrentTheme()))
			} else {
				sb.WriteString("\n")
			}
		}
	}
	
	// Split into lines for scrolling
	m.renderedLines = strings.Split(sb.String(), "\n")
}

// findGlobalHunkIndex finds the global index for a file/hunk combination
func (m Model) findGlobalHunkIndex(fileIdx, hunkIdx int) int {
	for i, meta := range m.hunkIndex {
		if meta.FileIndex == fileIdx && meta.HunkIndex == hunkIdx {
			return i
		}
	}
	return 0
}

// renderStatusLine renders the bottom status line
func (m Model) renderStatusLine() string {
	t := theme.CurrentTheme()
	statusStyle := styles.NewStyle().
		Background(t.BackgroundPanel()).
		Foreground(t.Text()).
		Bold(true)
	
	statusText := m.GetStatusLine()
	
	// Pad to full width
	if len(statusText) < m.width {
		statusText += strings.Repeat(" ", m.width-len(statusText))
	}
	
	return statusStyle.Render(statusText)
}

// renderHelp renders the help screen
func (m Model) renderHelp() string {
	t := theme.CurrentTheme()
	var sb strings.Builder
	
	titleStyle := styles.NewStyle().Foreground(t.Accent()).Bold(true)
	keyStyle := styles.NewStyle().Foreground(t.Info()).Bold(true)
	descStyle := styles.NewStyle().Foreground(t.Text())
	
	sb.WriteString(titleStyle.Render("Differential - Keyboard Shortcuts"))
	sb.WriteString("\n\n")
	
	shortcuts := []struct{ keys, desc string }{
		{"j / ↓", "Move down one line"},
		{"k / ↑", "Move up one line"},
		{"n / N", "Next file"},
		{"p / P", "Previous file"},
		{"] / Tab", "Next hunk"},
		{"[ / Shift+Tab", "Previous hunk"},
		{"Ctrl+f / PgDn", "Page down"},
		{"Ctrl+b / PgUp", "Page up"},
		{"g / Home", "Go to top"},
		{"G / End", "Go to bottom"},
		{"/", "Search"},
		{"t / T", "Toggle file tree"},
		{"? / h", "Toggle this help"},
		{"q / Esc", "Quit"},
	}
	
	for _, s := range shortcuts {
		sb.WriteString("  ")
		sb.WriteString(keyStyle.Render(fmt.Sprintf("%-20s", s.keys)))
		sb.WriteString(descStyle.Render(s.desc))
		sb.WriteString("\n")
	}
	
	sb.WriteString("\n")
	sb.WriteString(descStyle.Render("Press any key to close this help"))
	
	return sb.String()
}
