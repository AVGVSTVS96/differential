package tui

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/avgvstvs96/differential/pkg/diff"
	"github.com/avgvstvs96/differential/pkg/styles"
	"github.com/avgvstvs96/differential/pkg/theme"
)

// FileTreeItem represents a file in the tree
type FileTreeItem struct {
	Path        string
	Status      string // "M", "A", "D", etc.
	AddedLines  int
	DeletedLines int
	IsSelected  bool
}

// BuildFileTree creates a file tree from diff results
func BuildFileTree(results []diff.DiffResult, currentFile int) []FileTreeItem {
	items := make([]FileTreeItem, 0, len(results))
	
	for i, result := range results {
		// Count lines
		added, deleted := 0, 0
		for _, hunk := range result.Hunks {
			for _, line := range hunk.Lines {
				if line.Kind == diff.LineAdded {
					added++
				} else if line.Kind == diff.LineRemoved {
					deleted++
				}
			}
		}
		
		// Determine status
		status := "M" // Modified (default)
		if result.OldFile == "/dev/null" {
			status = "A" // Added
		} else if result.NewFile == "/dev/null" {
			status = "D" // Deleted
		}
		
		items = append(items, FileTreeItem{
			Path:         result.NewFile,
			Status:       status,
			AddedLines:   added,
			DeletedLines: deleted,
			IsSelected:   i == currentFile,
		})
	}
	
	return items
}

// RenderFileTree renders the file tree panel
func RenderFileTree(items []FileTreeItem, width, height int) string {
	t := theme.CurrentTheme()
	var sb strings.Builder
	
	// Title
	titleStyle := styles.NewStyle().
		Foreground(t.Accent()).
		Bold(true).
		Background(t.BackgroundPanel())
	
	title := " Files "
	if len(title) < width {
		title += strings.Repeat(" ", width-len(title))
	}
	sb.WriteString(titleStyle.Render(title))
	sb.WriteString("\n")
	
	// Render each file item
	visibleLines := height - 2 // Reserve space for title and padding
	for i, item := range items {
		if i >= visibleLines {
			break
		}
		
		sb.WriteString(renderFileTreeItem(item, width, t))
		sb.WriteString("\n")
	}
	
	// Fill remaining space
	for i := len(items); i < visibleLines; i++ {
		sb.WriteString(strings.Repeat(" ", width))
		sb.WriteString("\n")
	}
	
	return sb.String()
}

// renderFileTreeItem renders a single file tree item
func renderFileTreeItem(item FileTreeItem, width int, t theme.Theme) string {
	// Status icon with color
	statusIcon := ""
	statusColor := t.Text()
	
	switch item.Status {
	case "M":
		statusIcon = "[M]"
		statusColor = t.Info()
	case "A":
		statusIcon = "[A]"
		statusColor = t.Success()
	case "D":
		statusIcon = "[D]"
		statusColor = t.Error()
	}
	
	// File name (basename for brevity)
	fileName := filepath.Base(item.Path)
	
	// Stats
	stats := fmt.Sprintf("+%d -%d", item.AddedLines, item.DeletedLines)
	
	// Build the line
	var sb strings.Builder
	
	// Selection indicator
	if item.IsSelected {
		selStyle := styles.NewStyle().
			Background(t.BackgroundElement()).
			Foreground(t.Text()).
			Bold(true)
		
		// Status
		statusStyle := styles.NewStyle().
			Foreground(statusColor).
			Background(t.BackgroundElement()).
			Bold(true)
		
		sb.WriteString(selStyle.Render(" ▶ "))
		sb.WriteString(statusStyle.Render(statusIcon))
		sb.WriteString(selStyle.Render(" " + fileName))
		
		// Pad and add stats
		currentLen := 3 + len(statusIcon) + 1 + len(fileName)
		statsLen := len(stats)
		padding := width - currentLen - statsLen - 1
		if padding > 0 {
			sb.WriteString(selStyle.Render(strings.Repeat(" ", padding)))
		}
		
		statsStyle := styles.NewStyle().
			Foreground(t.TextMuted()).
			Background(t.BackgroundElement())
		sb.WriteString(statsStyle.Render(stats))
		
		// Fill to width
		if currentLen+padding+statsLen < width {
			sb.WriteString(selStyle.Render(strings.Repeat(" ", width-currentLen-padding-statsLen)))
		}
	} else {
		// Not selected - normal style
		statusStyle := styles.NewStyle().
			Foreground(statusColor).
			Bold(true)
		
		nameStyle := styles.NewStyle().
			Foreground(t.Text())
		
		sb.WriteString("   ") // No arrow for non-selected
		sb.WriteString(statusStyle.Render(statusIcon))
		sb.WriteString(nameStyle.Render(" " + fileName))
		
		// Pad and add stats
		currentLen := 3 + len(statusIcon) + 1 + len(fileName)
		statsLen := len(stats)
		padding := width - currentLen - statsLen - 1
		if padding > 0 {
			sb.WriteString(strings.Repeat(" ", padding))
		}
		
		statsStyle := styles.NewStyle().
			Foreground(t.TextMuted())
		sb.WriteString(statsStyle.Render(stats))
		
		// Fill to width
		if currentLen+padding+statsLen < width {
			sb.WriteString(strings.Repeat(" ", width-currentLen-padding-statsLen))
		}
	}
	
	return sb.String()
}
