package diff

// LineType represents the type of change for a line in a diff
type LineType int

const (
	LineContext LineType = iota // Line exists in both files
	LineAdded                   // Line added in the new file
	LineRemoved                 // Line removed from the old file
)

// Segment represents a highlighted segment within a line for character-level diff
type Segment struct {
	Start int      // Start position in the line (byte offset)
	End   int      // End position in the line (byte offset)
	Type  LineType // Type of change (Added/Removed)
	Text  string   // The actual text of the segment
}

// DiffLine represents a single line in a diff
type DiffLine struct {
	OldLineNo int       // Line number in old file (0 for added lines)
	NewLineNo int       // Line number in new file (0 for removed lines)
	Kind      LineType  // Type of line (added, removed, context)
	Content   string    // Content of the line (without diff markers)
	Segments  []Segment // Segments for intraline highlighting
}

// Hunk represents a contiguous block of changes in a diff
type Hunk struct {
	Header string     // The @@ header line
	Lines  []DiffLine // All lines in this hunk
}

// DiffResult contains the complete parsed diff
type DiffResult struct {
	OldFile string // Old file path
	NewFile string // New file path
	Hunks   []Hunk // All hunks in the diff
	IsBinary bool  // Whether this is a binary file diff
}

// LinePair is used for side-by-side rendering
type LinePair struct {
	Left  *DiffLine // nil for added lines
	Right *DiffLine // nil for removed lines
}

// ViewMode represents how the diff should be displayed
type ViewMode int

const (
	ViewUnified ViewMode = iota
	ViewSideBySide
)

// RenderOptions contains options for rendering diffs
type RenderOptions struct {
	Width           int      // Terminal width
	ViewMode        ViewMode // Unified or side-by-side
	ShowLineNumbers bool     // Whether to show line numbers
	ContextLines    int      // Number of context lines
	TabWidth        int      // Tab character width
}