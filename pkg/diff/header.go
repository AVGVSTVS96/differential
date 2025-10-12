package diff

import (
	"fmt"
	"strings"

	"github.com/avgvstvs96/differential/pkg/styles"
	"github.com/avgvstvs96/differential/pkg/theme"
)

// RenderFileHeader renders a file header with status and statistics
func RenderFileHeader(result DiffResult, showStats bool, hyperlinksEnabled bool, graphicalStats bool) string {
	t := theme.CurrentTheme()
	var sb strings.Builder

	// Status indicator with color
	status := result.Metadata.Status.String()
	var statusStyle styles.Style
	switch result.Metadata.Status {
	case FileAdded:
		statusStyle = styles.NewStyle().Foreground(t.DiffAdded()).Bold(true)
	case FileDeleted:
		statusStyle = styles.NewStyle().Foreground(t.DiffRemoved()).Bold(true)
	case FileRenamed:
		statusStyle = styles.NewStyle().Foreground(t.Warning()).Bold(true)
	case FileModified:
		statusStyle = styles.NewStyle().Foreground(t.Info()).Bold(true)
	default:
		statusStyle = styles.NewStyle().Foreground(t.TextMuted())
	}

	// File path
	filePath := result.NewFile
	if result.Metadata.Status == FileDeleted {
		filePath = result.OldFile
	}
	
	// Apply hyperlink if enabled
	displayPath := filePath
	if hyperlinksEnabled {
		displayPath = CreateFileHyperlink(filePath, filePath, true)
	}
	
	filePathStyle := styles.NewStyle().Foreground(t.Text()).Bold(true)

	// Build header line
	sb.WriteString(statusStyle.Render(fmt.Sprintf("[%s]", status)))
	sb.WriteString(" ")
	sb.WriteString(filePathStyle.Render(displayPath))

	// Add rename information if applicable
	if result.Metadata.Status == FileRenamed && result.OldFile != result.NewFile {
		sb.WriteString(" ")
		sb.WriteString(styles.NewStyle().Foreground(t.TextMuted()).Render(
			fmt.Sprintf("(renamed from %s)", result.OldFile),
		))
	}

	// Add statistics if requested
	if showStats {
		added := result.Metadata.LinesAdded
		removed := result.Metadata.LinesRemoved
		total := added + removed

		if total > 0 {
			sb.WriteString(" ")
			
			if graphicalStats {
				// Graphical representation with bars
				sb.WriteString(renderGraphicalStats(added, removed, t))
			} else {
				// Text-only representation
				addedStyle := styles.NewStyle().Foreground(t.DiffAdded())
				removedStyle := styles.NewStyle().Foreground(t.DiffRemoved())
				mutedStyle := styles.NewStyle().Foreground(t.TextMuted())

				sb.WriteString(mutedStyle.Render("("))
				if added > 0 {
					sb.WriteString(addedStyle.Render(fmt.Sprintf("+%d", added)))
				}
				if added > 0 && removed > 0 {
					sb.WriteString(mutedStyle.Render("/"))
				}
				if removed > 0 {
					sb.WriteString(removedStyle.Render(fmt.Sprintf("-%d", removed)))
				}
				sb.WriteString(mutedStyle.Render(")"))
			}
		}
	}

	sb.WriteString("\n")
	return sb.String()
}

// renderGraphicalStats renders statistics as a graphical bar
func renderGraphicalStats(added, removed int, t theme.Theme) string {
	total := added + removed
	if total == 0 {
		return ""
	}
	
	var sb strings.Builder
	mutedStyle := styles.NewStyle().Foreground(t.TextMuted())
	addedStyle := styles.NewStyle().Foreground(t.DiffAdded())
	removedStyle := styles.NewStyle().Foreground(t.DiffRemoved())
	
	// Calculate bar segments (max 10 chars total)
	maxBarWidth := 10
	addedBars := 0
	removedBars := 0
	
	if total > 0 {
		addedBars = (added * maxBarWidth) / total
		removedBars = maxBarWidth - addedBars
		
		// Ensure at least one bar if there are changes
		if added > 0 && addedBars == 0 {
			addedBars = 1
			removedBars--
		}
		if removed > 0 && removedBars == 0 {
			removedBars = 1
			addedBars--
		}
	}
	
	// Build the graphical representation
	sb.WriteString(mutedStyle.Render("["))
	if addedBars > 0 {
		sb.WriteString(addedStyle.Render(strings.Repeat("=", addedBars)))
	}
	if removedBars > 0 {
		sb.WriteString(removedStyle.Render(strings.Repeat("=", removedBars)))
	}
	// Pad to full width
	remaining := maxBarWidth - addedBars - removedBars
	if remaining > 0 {
		sb.WriteString(mutedStyle.Render(strings.Repeat(" ", remaining)))
	}
	sb.WriteString(mutedStyle.Render("]"))
	
	// Add numeric stats
	sb.WriteString(" ")
	sb.WriteString(addedStyle.Render(fmt.Sprintf("+%d", added)))
	sb.WriteString(" ")
	sb.WriteString(removedStyle.Render(fmt.Sprintf("-%d", removed)))
	
	return sb.String()
}

// RenderDiffStats renders a summary of all file changes
func RenderDiffStats(results []DiffResult) string {
	if len(results) == 0 {
		return ""
	}

	t := theme.CurrentTheme()
	var sb strings.Builder

	// Count files by status
	statusCounts := make(map[FileStatus]int)
	totalAdded := 0
	totalRemoved := 0

	for _, result := range results {
		statusCounts[result.Metadata.Status]++
		totalAdded += result.Metadata.LinesAdded
		totalRemoved += result.Metadata.LinesRemoved
	}

	// Summary line
	totalFiles := len(results)
	sb.WriteString(styles.NewStyle().Bold(true).Render(
		fmt.Sprintf("%d file(s) changed", totalFiles),
	))
	sb.WriteString(", ")
	sb.WriteString(styles.NewStyle().Foreground(t.DiffAdded()).Render(
		fmt.Sprintf("%d insertion(s)", totalAdded),
	))
	sb.WriteString(", ")
	sb.WriteString(styles.NewStyle().Foreground(t.DiffRemoved()).Render(
		fmt.Sprintf("%d deletion(s)", totalRemoved),
	))
	sb.WriteString("\n")

	// Status breakdown if there are different types
	if len(statusCounts) > 1 || (len(statusCounts) == 1 && statusCounts[FileModified] != totalFiles) {
		sb.WriteString(styles.NewStyle().Foreground(t.TextMuted()).Render("  "))
		
		first := true
		for status, count := range statusCounts {
			if !first {
				sb.WriteString(styles.NewStyle().Foreground(t.TextMuted()).Render(", "))
			}
			first = false

			var statusStyle styles.Style
			var statusName string
			switch status {
			case FileAdded:
				statusStyle = styles.NewStyle().Foreground(t.DiffAdded())
				statusName = "added"
			case FileDeleted:
				statusStyle = styles.NewStyle().Foreground(t.DiffRemoved())
				statusName = "deleted"
			case FileRenamed:
				statusStyle = styles.NewStyle().Foreground(t.Warning())
				statusName = "renamed"
			case FileModified:
				statusStyle = styles.NewStyle().Foreground(t.Info())
				statusName = "modified"
			default:
				statusStyle = styles.NewStyle().Foreground(t.TextMuted())
				statusName = "unknown"
			}

			sb.WriteString(statusStyle.Render(fmt.Sprintf("%d %s", count, statusName)))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	return sb.String()
}
