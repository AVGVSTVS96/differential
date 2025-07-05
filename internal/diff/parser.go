package diff

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	// Regular expressions for parsing diff format
	fileHeaderRegex = regexp.MustCompile(`^diff --git a/(.+) b/(.+)$`)
	oldFileRegex    = regexp.MustCompile(`^--- (?:a/)?(.+?)(?:\s+\d{4}-\d{2}-\d{2}.*)?$`)
	newFileRegex    = regexp.MustCompile(`^\+\+\+ (?:b/)?(.+?)(?:\s+\d{4}-\d{2}-\d{2}.*)?$`)
	hunkHeaderRegex = regexp.MustCompile(`^@@ -(\d+)(?:,(\d+))? \+(\d+)(?:,(\d+))? @@`)
	binaryFileRegex = regexp.MustCompile(`^Binary files? .* differ$`)
)

// ParseUnifiedDiff parses a unified diff format string into a DiffResult
func ParseUnifiedDiff(diffText string) (*DiffResult, error) {
	if diffText == "" {
		return &DiffResult{}, nil
	}

	result := &DiffResult{
		Hunks: make([]Hunk, 0, 10),
	}

	scanner := bufio.NewScanner(strings.NewReader(diffText))
	var currentHunk *Hunk
	var oldLine, newLine int
	inFileHeader := true

	for scanner.Scan() {
		line := scanner.Text()

		// Check for binary file
		if binaryFileRegex.MatchString(line) {
			result.IsBinary = true
			return result, nil
		}

		// File headers
		if inFileHeader {
			if matches := oldFileRegex.FindStringSubmatch(line); matches != nil {
				result.OldFile = matches[1]
				continue
			}
			if matches := newFileRegex.FindStringSubmatch(line); matches != nil {
				result.NewFile = matches[1]
				inFileHeader = false
				continue
			}
			// Skip other header lines (index, mode, etc.)
			continue
		}

		// Hunk header: @@ -10,7 +10,7 @@ func main() {
		if matches := hunkHeaderRegex.FindStringSubmatch(line); matches != nil {
			// Save previous hunk
			if currentHunk != nil {
				result.Hunks = append(result.Hunks, *currentHunk)
			}

			// Parse line numbers
			oldLine, _ = strconv.Atoi(matches[1])
			newLine, _ = strconv.Atoi(matches[3])

			currentHunk = &Hunk{
				Header: line,
				Lines:  make([]DiffLine, 0, 10),
			}
			continue
		}

		// Skip "\ No newline at end of file"
		if strings.HasPrefix(line, "\\") {
			continue
		}

		// Parse diff lines
		if currentHunk != nil && len(line) > 0 {
			dl := parseDiffLine(line, &oldLine, &newLine)
			currentHunk.Lines = append(currentHunk.Lines, dl)
		}
	}

	// Don't forget the last hunk
	if currentHunk != nil {
		result.Hunks = append(result.Hunks, *currentHunk)
	}

	return result, scanner.Err()
}

// parseDiffLine parses a single line from a diff
func parseDiffLine(line string, oldLine, newLine *int) DiffLine {
	if len(line) == 0 {
		return DiffLine{
			Kind:      LineContext,
			OldLineNo: *oldLine,
			NewLineNo: *newLine,
			Content:   "",
		}
	}

	var dl DiffLine

	switch line[0] {
	case '+':
		dl.Kind = LineAdded
		dl.NewLineNo = *newLine
		dl.Content = line[1:]
		(*newLine)++

	case '-':
		dl.Kind = LineRemoved
		dl.OldLineNo = *oldLine
		dl.Content = line[1:]
		(*oldLine)++

	case ' ':
		dl.Kind = LineContext
		dl.OldLineNo = *oldLine
		dl.NewLineNo = *newLine
		dl.Content = line[1:]
		(*oldLine)++
		(*newLine)++

	default:
		// Malformed line, treat as context
		dl.Kind = LineContext
		dl.OldLineNo = *oldLine
		dl.NewLineNo = *newLine
		dl.Content = line
		(*oldLine)++
		(*newLine)++
	}

	return dl
}

// PairLines groups diff lines for side-by-side rendering
func PairLines(lines []DiffLine) []LinePair {
	var pairs []LinePair
	i := 0

	for i < len(lines) {
		switch lines[i].Kind {
		case LineRemoved:
			// Check if next line is addition (change)
			if i+1 < len(lines) && lines[i+1].Kind == LineAdded {
				pairs = append(pairs, LinePair{
					Left:  &lines[i],
					Right: &lines[i+1],
				})
				i += 2
			} else {
				// Deletion only
				pairs = append(pairs, LinePair{
					Left:  &lines[i],
					Right: nil,
				})
				i++
			}

		case LineAdded:
			// Addition only (not paired with removal)
			pairs = append(pairs, LinePair{
				Left:  nil,
				Right: &lines[i],
			})
			i++

		case LineContext:
			// Same line on both sides
			pairs = append(pairs, LinePair{
				Left:  &lines[i],
				Right: &lines[i],
			})
			i++
		}
	}

	return pairs
}

// GetFileExtension returns the file extension for syntax highlighting
func GetFileExtension(filename string) string {
	if filename == "" {
		return ""
	}

	parts := strings.Split(filename, ".")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}

	// Try to detect from common filenames without extensions
	base := parts[0]
	switch base {
	case "Dockerfile":
		return "dockerfile"
	case "Makefile":
		return "makefile"
	case "README":
		return "md"
	}

	return ""
}

// CountChanges returns the number of additions and deletions in a diff
func (d *DiffResult) CountChanges() (additions, deletions int) {
	for _, hunk := range d.Hunks {
		for _, line := range hunk.Lines {
			switch line.Kind {
			case LineAdded:
				additions++
			case LineRemoved:
				deletions++
			}
		}
	}
	return
}

// String returns a string representation of the diff result (for debugging)
func (d *DiffResult) String() string {
	return fmt.Sprintf("DiffResult{OldFile: %s, NewFile: %s, Hunks: %d}", 
		d.OldFile, d.NewFile, len(d.Hunks))
}