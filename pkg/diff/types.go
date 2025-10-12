package diff

import (
	"regexp"
)

// -------------------------------------------------------------------------
// Core Types
// -------------------------------------------------------------------------

// LineType represents the kind of line in a diff.
type LineType int

const (
	LineContext LineType = iota // Line exists in both files
	LineAdded                   // Line added in the new file
	LineRemoved                 // Line removed from the old file
)

var (
	ansiRegex = regexp.MustCompile(`\x1b(?:[@-Z\\-_]|\[[0-9?]*(?:;[0-9?]*)*[@-~])`)
)

// Segment represents a portion of a line for intra-line highlighting
type Segment struct {
	Start int
	End   int
	Type  LineType
	Text  string
}

// DiffLine represents a single line in a diff
type DiffLine struct {
	OldLineNo int       // Line number in old file (0 for added lines)
	NewLineNo int       // Line number in new file (0 for removed lines)
	Kind      LineType  // Type of line (added, removed, context)
	Content   string    // Content of the line
	Segments  []Segment // Segments for intraline highlighting
}

// Hunk represents a section of changes in a diff
type Hunk struct {
	Header string
	Lines  []DiffLine
}

// FileStatus represents the modification status of a file
type FileStatus int

const (
	FileModified FileStatus = iota // File was modified
	FileAdded                       // File was added
	FileDeleted                     // File was deleted
	FileRenamed                     // File was renamed
	FileCopied                      // File was copied
)

// String returns the short string representation of file status
func (fs FileStatus) String() string {
	switch fs {
	case FileModified:
		return "M"
	case FileAdded:
		return "A"
	case FileDeleted:
		return "D"
	case FileRenamed:
		return "R"
	case FileCopied:
		return "C"
	default:
		return "?"
	}
}

// FileMetadata contains metadata about a file in a diff
type FileMetadata struct {
	OldPath    string     // Path in old version
	NewPath    string     // Path in new version
	Status     FileStatus // File modification status
	LinesAdded int        // Number of lines added
	LinesRemoved int      // Number of lines removed
}

// DiffResult contains the parsed result of a diff
type DiffResult struct {
	OldFile  string
	NewFile  string
	Metadata FileMetadata // File metadata
	Hunks    []Hunk
}

// linePair represents a pair of lines for side-by-side display
type linePair struct {
	left  *DiffLine
	right *DiffLine
}

// UnifiedConfig configures the rendering of unified diffs
type UnifiedConfig struct {
	Width int
}

// UnifiedOption modifies a UnifiedConfig
type UnifiedOption func(*UnifiedConfig)

// NewUnifiedConfig creates a UnifiedConfig with default values
func NewUnifiedConfig(opts ...UnifiedOption) UnifiedConfig {
	config := UnifiedConfig{
		Width: 80,
	}
	for _, opt := range opts {
		opt(&config)
	}
	return config
}

// NewSideBySideConfig creates a SideBySideConfig with default values
func NewSideBySideConfig(opts ...UnifiedOption) UnifiedConfig {
	config := UnifiedConfig{
		Width: 160,
	}
	for _, opt := range opts {
		opt(&config)
	}
	return config
}

// WithWidth sets the width for unified view
func WithWidth(width int) UnifiedOption {
	return func(u *UnifiedConfig) {
		if width > 0 {
			u.Width = width
		}
	}
}
