package diff

import (
	"bufio"
	"strconv"
	"strings"
)

// ParseUnifiedDiff parses a unified diff format string into structured data
func ParseUnifiedDiff(diff string) (DiffResult, error) {
	var result DiffResult
	var currentHunk *Hunk
	result.Hunks = make([]Hunk, 0, 10) // Pre-allocate with a reasonable capacity

	scanner := bufio.NewScanner(strings.NewReader(diff))
	var oldLine, newLine int
	inFileHeader := true

	for scanner.Scan() {
		line := scanner.Text()

		if inFileHeader {
			if strings.HasPrefix(line, "--- a/") {
				result.OldFile = line[6:]
				result.Metadata.OldPath = result.OldFile
				continue
			}
			if strings.HasPrefix(line, "--- /dev/null") {
				result.OldFile = "/dev/null"
				result.Metadata.OldPath = "/dev/null"
				result.Metadata.Status = FileAdded
				continue
			}
			if strings.HasPrefix(line, "+++ b/") {
				result.NewFile = line[6:]
				result.Metadata.NewPath = result.NewFile
				inFileHeader = false
				continue
			}
			if strings.HasPrefix(line, "+++ /dev/null") {
				result.NewFile = "/dev/null"
				result.Metadata.NewPath = "/dev/null"
				result.Metadata.Status = FileDeleted
				inFileHeader = false
				continue
			}
		}

		if strings.HasPrefix(line, "@@") {
			if currentHunk != nil {
				result.Hunks = append(result.Hunks, *currentHunk)
			}
			currentHunk = &Hunk{
				Header: line,
				Lines:  make([]DiffLine, 0, 10), // Pre-allocate
			}

			// Manual parsing of hunk header is faster than regex
			parts := strings.Split(line, " ")
			if len(parts) > 2 {
				oldRange := strings.Split(parts[1][1:], ",")
				newRange := strings.Split(parts[2][1:], ",")
				oldLine, _ = strconv.Atoi(oldRange[0])
				newLine, _ = strconv.Atoi(newRange[0])
			}
			continue
		}

		if strings.HasPrefix(line, "\\ No newline at end of file") || currentHunk == nil {
			continue
		}

		var dl DiffLine
		dl.Content = line
		if len(line) > 0 {
			switch line[0] {
			case '+':
				dl.Kind = LineAdded
				dl.NewLineNo = newLine
				dl.Content = line[1:]
				newLine++
				result.Metadata.LinesAdded++
			case '-':
				dl.Kind = LineRemoved
				dl.OldLineNo = oldLine
				dl.Content = line[1:]
				oldLine++
				result.Metadata.LinesRemoved++
			default: // context line
				dl.Kind = LineContext
				dl.OldLineNo = oldLine
				dl.NewLineNo = newLine
				oldLine++
				newLine++
			}
		} else { // empty context line
			dl.Kind = LineContext
			dl.OldLineNo = oldLine
			dl.NewLineNo = newLine
			oldLine++
			newLine++
		}
		currentHunk.Lines = append(currentHunk.Lines, dl)
	}

	if currentHunk != nil {
		result.Hunks = append(result.Hunks, *currentHunk)
	}

	// Determine file status if not already set
	if result.Metadata.Status == 0 {
		if result.OldFile == "/dev/null" {
			result.Metadata.Status = FileAdded
		} else if result.NewFile == "/dev/null" {
			result.Metadata.Status = FileDeleted
		} else if result.OldFile != result.NewFile {
			result.Metadata.Status = FileRenamed
		} else {
			result.Metadata.Status = FileModified
		}
	}

	return result, scanner.Err()
}
