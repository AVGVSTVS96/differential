package tui

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/avgvstvs96/differential/pkg/styles"
	"github.com/avgvstvs96/differential/pkg/theme"
)

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

// stripANSI removes ANSI escape codes from a string
func stripANSI(s string) string {
	return ansiRegex.ReplaceAllString(s, "")
}

// visualWidth returns the visual width of a string (without ANSI codes)
func visualWidth(s string) int {
	return len(stripANSI(s))
}

// View renders the current view of the diff
func (m Model) View() string {
	if m.quitting {
		return ""
	}
	
	if m.showHelp {
		return m.renderHelp()
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
	// For now, just show full view with a message that split view is disabled
	// TODO: Fix split view rendering with proper ANSI handling
	return m.renderFullView()
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
