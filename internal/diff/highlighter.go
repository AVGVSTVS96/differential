package diff

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// HighlightIntralineChanges computes character-level differences within changed lines
func HighlightIntralineChanges(h *Hunk) {
	dmp := diffmatchpatch.New()

	for i := 0; i < len(h.Lines); i++ {
		// Look for adjacent removed/added line pairs
		if i+1 < len(h.Lines) &&
			h.Lines[i].Kind == LineRemoved &&
			h.Lines[i+1].Kind == LineAdded {

			oldLine := &h.Lines[i]
			newLine := &h.Lines[i+1]

			// Compute character-level differences
			diffs := dmp.DiffMain(oldLine.Content, newLine.Content, false)
			diffs = dmp.DiffCleanupSemantic(diffs)

			// Build segments for highlighting
			oldSegments := []Segment{}
			newSegments := []Segment{}
			oldPos, newPos := 0, 0

			for _, diff := range diffs {
				switch diff.Type {
				case diffmatchpatch.DiffDelete:
					oldSegments = append(oldSegments, Segment{
						Start: oldPos,
						End:   oldPos + len(diff.Text),
						Type:  LineRemoved,
						Text:  diff.Text,
					})
					oldPos += len(diff.Text)

				case diffmatchpatch.DiffInsert:
					newSegments = append(newSegments, Segment{
						Start: newPos,
						End:   newPos + len(diff.Text),
						Type:  LineAdded,
						Text:  diff.Text,
					})
					newPos += len(diff.Text)

				case diffmatchpatch.DiffEqual:
					oldPos += len(diff.Text)
					newPos += len(diff.Text)
				}
			}

			// Apply segments to lines
			oldLine.Segments = oldSegments
			newLine.Segments = newSegments

			i++ // Skip the next line since we processed it
		}
	}
}

// ApplyHighlighting applies ANSI color codes to highlight segments while preserving existing ANSI sequences
func ApplyHighlighting(content string, segments []Segment, segmentType LineType, highlightStyle string) string {
	if len(segments) == 0 {
		return content
	}

	// Find all ANSI sequences in the content
	ansiRegex := regexp.MustCompile(`\x1b(?:[@-Z\\-_]|\[[0-9?]*(?:;[0-9?]*)*[@-~])`)
	ansiMatches := ansiRegex.FindAllStringIndex(content, -1)

	// Build mapping of visible position to byte position and ANSI sequences
	visibleToBytePos := make(map[int]int)
	ansiSequences := make(map[int]string)
	lastAnsiSeq := "\x1b[0m" // Default reset

	visibleIdx := 0
	byteIdx := 0

	for byteIdx < len(content) {
		isAnsi := false

		// Check if current position is an ANSI sequence
		for _, match := range ansiMatches {
			if match[0] == byteIdx {
				ansiSeq := content[match[0]:match[1]]
				ansiSequences[visibleIdx] = ansiSeq
				lastAnsiSeq = ansiSeq
				byteIdx = match[1]
				isAnsi = true
				break
			}
		}

		if !isAnsi {
			// Map visible position to byte position
			visibleToBytePos[visibleIdx] = byteIdx

			// Store last ANSI sequence for this position
			if _, exists := ansiSequences[visibleIdx]; !exists {
				ansiSequences[visibleIdx] = lastAnsiSeq
			}

			// Advance by one rune
			_, size := utf8.DecodeRuneInString(content[byteIdx:])
			byteIdx += size
			visibleIdx++
		}
	}

	// Apply highlighting
	var sb strings.Builder
	inSelection := false
	currentPos := 0

	for i := 0; i < len(content); {
		// Check if this is an ANSI sequence
		isAnsiSeq := false
		for _, match := range ansiMatches {
			if match[0] == i {
				sb.WriteString(content[match[0]:match[1]])
				i = match[1]
				isAnsiSeq = true
				break
			}
		}

		if !isAnsiSeq {
			// Check if we're entering or leaving a highlighted segment
			for _, seg := range segments {
				if seg.Type == segmentType {
					if currentPos == seg.Start && !inSelection {
						sb.WriteString(highlightStyle)
						inSelection = true
					}
					if currentPos == seg.End && inSelection {
						sb.WriteString("\x1b[0m") // Reset
						// Restore previous ANSI state
						if prevAnsi, ok := ansiSequences[currentPos]; ok {
							sb.WriteString(prevAnsi)
						}
						inSelection = false
					}
				}
			}

			// Write the character
			r, size := utf8.DecodeRuneInString(content[i:])
			sb.WriteRune(r)
			currentPos++
			i += size
		}
	}

	// Make sure we reset if still in selection
	if inSelection {
		sb.WriteString("\x1b[0m")
	}

	return sb.String()
}

// CreateHighlightStyle creates ANSI escape sequences for highlighting
func CreateHighlightStyle(fg, bg string) string {
	return fmt.Sprintf("\x1b[38;2;%s;48;2;%sm", fg, bg)
}

// StripANSI removes all ANSI escape sequences from a string
func StripANSI(str string) string {
	ansiRegex := regexp.MustCompile(`\x1b(?:[@-Z\\-_]|\[[0-9?]*(?:;[0-9?]*)*[@-~])`)
	return ansiRegex.ReplaceAllString(str, "")
}

// VisibleLength returns the visible length of a string (excluding ANSI sequences)
func VisibleLength(str string) int {
	stripped := StripANSI(str)
	return utf8.RuneCountInString(stripped)
}

// TruncateString truncates a string to a visible width, preserving ANSI sequences
func TruncateString(str string, width int) string {
	if width <= 0 {
		return ""
	}

	ansiRegex := regexp.MustCompile(`\x1b(?:[@-Z\\-_]|\[[0-9?]*(?:;[0-9?]*)*[@-~])`)
	ansiMatches := ansiRegex.FindAllStringIndex(str, -1)

	var sb strings.Builder
	visibleCount := 0
	i := 0

	for i < len(str) && visibleCount < width {
		// Check if we're at an ANSI sequence
		isAnsi := false
		for _, match := range ansiMatches {
			if match[0] == i {
				sb.WriteString(str[match[0]:match[1]])
				i = match[1]
				isAnsi = true
				break
			}
		}

		if !isAnsi {
			r, size := utf8.DecodeRuneInString(str[i:])
			sb.WriteRune(r)
			i += size
			visibleCount++
		}
	}

	// Copy any remaining ANSI sequences
	for _, match := range ansiMatches {
		if match[0] >= i {
			sb.WriteString(str[match[0]:match[1]])
		}
	}

	return sb.String()
}