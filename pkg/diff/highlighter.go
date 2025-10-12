package diff

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss/v2/compat"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// HighlightIntralineChanges updates lines in a hunk to show character-level differences
func HighlightIntralineChanges(h *Hunk) {
	var updated []DiffLine
	dmp := diffmatchpatch.New()

	for i := 0; i < len(h.Lines); i++ {
		// Look for removed line followed by added line
		if i+1 < len(h.Lines) &&
			h.Lines[i].Kind == LineRemoved &&
			h.Lines[i+1].Kind == LineAdded {

			oldLine := h.Lines[i]
			newLine := h.Lines[i+1]

			// Find character-level differences
			patches := dmp.DiffMain(oldLine.Content, newLine.Content, false)
			patches = dmp.DiffCleanupSemantic(patches)
			patches = dmp.DiffCleanupMerge(patches)
			patches = dmp.DiffCleanupEfficiency(patches)

			segments := make([]Segment, 0)

			removeStart := 0
			addStart := 0
			for _, patch := range patches {
				switch patch.Type {
				case diffmatchpatch.DiffDelete:
					segments = append(segments, Segment{
						Start: removeStart,
						End:   removeStart + len(patch.Text),
						Type:  LineRemoved,
						Text:  patch.Text,
					})
					removeStart += len(patch.Text)
				case diffmatchpatch.DiffInsert:
					segments = append(segments, Segment{
						Start: addStart,
						End:   addStart + len(patch.Text),
						Type:  LineAdded,
						Text:  patch.Text,
					})
					addStart += len(patch.Text)
				default:
					// Context text, no highlighting needed
					removeStart += len(patch.Text)
					addStart += len(patch.Text)
				}
			}
			oldLine.Segments = segments
			newLine.Segments = segments

			updated = append(updated, oldLine, newLine)
			i++ // Skip the next line as we've already processed it
		} else {
			updated = append(updated, h.Lines[i])
		}
	}

	h.Lines = updated
}

// applyHighlighting applies intra-line highlighting to a piece of text
func applyHighlighting(content string, segments []Segment, segmentType LineType, highlightBg compat.AdaptiveColor, highlightFg compat.AdaptiveColor) string {
	// Find all ANSI sequences in the content
	ansiMatches := ansiRegex.FindAllStringIndex(content, -1)

	// Build a mapping of visible character positions to their actual indices
	visibleIdx := 0
	ansiSequences := make(map[int]string)
	lastAnsiSeq := "\x1b[0m" // Default reset sequence

	for i := 0; i < len(content); {
		isAnsi := false
		for _, match := range ansiMatches {
			if match[0] == i {
				ansiSequences[visibleIdx] = content[match[0]:match[1]]
				lastAnsiSeq = content[match[0]:match[1]]
				i = match[1]
				isAnsi = true
				break
			}
		}
		if isAnsi {
			continue
		}

		// For non-ANSI positions, store the last ANSI sequence
		if _, exists := ansiSequences[visibleIdx]; !exists {
			ansiSequences[visibleIdx] = lastAnsiSeq
		}
		visibleIdx++

		// Properly advance by UTF-8 rune, not byte
		_, size := utf8.DecodeRuneInString(content[i:])
		i += size
	}

	// Apply highlighting
	var sb strings.Builder
	inSelection := false
	currentPos := 0

	// Get foreground and background colors
	bgColor := adaptiveColorToColor(highlightBg)
	fgColor := adaptiveColorToColor(highlightFg)

	for i := 0; i < len(content); {
		// Check if we're at an ANSI sequence
		isAnsi := false
		for _, match := range ansiMatches {
			if match[0] == i {
				sb.WriteString(content[match[0]:match[1]]) // Preserve ANSI sequence
				i = match[1]
				isAnsi = true
				break
			}
		}
		if isAnsi {
			continue
		}

		// Check for segment boundaries
		for _, seg := range segments {
			if seg.Type == segmentType {
				if currentPos == seg.Start {
					inSelection = true
				}
				if currentPos == seg.End {
					inSelection = false
				}
			}
		}

		// Get current character (properly handle UTF-8)
		r, size := utf8.DecodeRuneInString(content[i:])
		char := string(r)

		if inSelection {
			// Get the current styling
			currentStyle := ansiSequences[currentPos]

			// Apply foreground and background highlight
			if fgColor != nil {
				sb.WriteString("\x1b[38;2;")
				r, g, b, _ := fgColor.RGBA()
				sb.WriteString(fmt.Sprintf("%d;%d;%dm", r>>8, g>>8, b>>8))
			} else {
				sb.WriteString("\x1b[49m")
			}
			if bgColor != nil {
				sb.WriteString("\x1b[48;2;")
				r, g, b, _ := bgColor.RGBA()
				sb.WriteString(fmt.Sprintf("%d;%d;%dm", r>>8, g>>8, b>>8))
			} else {
				sb.WriteString("\x1b[39m")
			}
			sb.WriteString(char)

			// Full reset of all attributes to ensure clean state
			sb.WriteString("\x1b[0m")

			// Reapply the original ANSI sequence
			sb.WriteString(currentStyle)
		} else {
			// Not in selection, just copy the character
			sb.WriteString(char)
		}

		currentPos++
		i += size
	}

	return sb.String()
}
