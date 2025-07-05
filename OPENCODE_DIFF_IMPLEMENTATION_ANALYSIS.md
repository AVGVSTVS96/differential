# OpenCode Diff Implementation Analysis

## Overview

OpenCode implements a sophisticated and beautiful terminal diff display that outperforms many existing diff tools. This document analyzes the implementation details to understand how it achieves its superior visual presentation.

## Core Architecture

### 1. Main Diff Component
**Location**: `/packages/tui/internal/components/diff/diff.go`

The core diff rendering is implemented in a single, well-structured Go file that handles:
- Unified diff parsing
- Character-level (intra-line) diff highlighting
- Syntax highlighting integration
- Both unified and side-by-side view modes
- Parallel rendering for performance

### 2. Key Libraries Used
- **Bubble Tea**: Terminal UI framework (charmbracelet/bubbletea/v2)
- **Lipgloss**: Terminal styling (charmbracelet/lipgloss/v2)
- **Chroma**: Syntax highlighting (alecthomas/chroma/v2)
- **go-diff**: Character-level diff computation (sergi/go-diff/diffmatchpatch)
- **ANSI utilities**: String width and truncation (charmbracelet/x/ansi)

## Implementation Details

### 1. Diff Data Structures

```go
type LineType int
const (
    LineContext LineType = iota // Line exists in both files
    LineAdded                   // Line added in the new file
    LineRemoved                 // Line removed from the old file
)

type Segment struct {
    Start int      // Start position in the line (byte offset)
    End   int      // End position in the line (byte offset)
    Type  LineType // Type of change (Added/Removed)
    Text  string   // The actual text of the segment
}

type DiffLine struct {
    OldLineNo int       // Line number in old file (0 for added lines)
    NewLineNo int       // Line number in new file (0 for removed lines)
    Kind      LineType  // Type of line (added, removed, context)
    Content   string    // Content of the line (without diff markers)
    Segments  []Segment // Segments for intraline highlighting
}

type Hunk struct {
    Header string      // The @@ header line
    Lines  []DiffLine  // All lines in this hunk
}

type DiffResult struct {
    OldFile string // Old file path
    NewFile string // New file path
    Hunks   []Hunk // All hunks in the diff
}

// For side-by-side rendering
type linePair struct {
    left  *DiffLine // nil for added lines
    right *DiffLine // nil for removed lines
}
```

**Data Flow**:
1. Git diff output → `ParseUnifiedDiff()` → `DiffResult`
2. `DiffResult` → `HighlightIntralineChanges()` → Updates `Segments` in paired lines
3. `DiffResult` → `RenderUnifiedHunk()` or `RenderSideBySideHunk()` → Terminal output

### 2. Character-Level Diff Highlighting

The most impressive feature is the character-level highlighting within changed lines.

**Algorithm** (`HighlightIntralineChanges`):
```go
func HighlightIntralineChanges(h *Hunk) {
    dmp := diffmatchpatch.New()
    
    for i := 0; i < len(h.Lines); i++ {
        // Look for adjacent removed/added line pairs
        if i+1 < len(h.Lines) &&
           h.Lines[i].Kind == LineRemoved &&
           h.Lines[i+1].Kind == LineAdded {
            
            oldLine := h.Lines[i]
            newLine := h.Lines[i+1]
            
            // Compute character-level differences
            diffs := dmp.DiffMain(oldLine.Content, newLine.Content, false)
            diffs = dmp.DiffCleanupSemantic(diffs)
            
            // Build segments for highlighting
            segments := []Segment{}
            oldPos, newPos := 0, 0
            
            for _, diff := range diffs {
                switch diff.Type {
                case diffmatchpatch.DiffDelete:
                    segments = append(segments, Segment{
                        Start: oldPos,
                        End:   oldPos + len(diff.Text),
                        Type:  LineRemoved,
                        Text:  diff.Text,
                    })
                    oldPos += len(diff.Text)
                    
                case diffmatchpatch.DiffInsert:
                    segments = append(segments, Segment{
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
            
            // Apply segments to both lines
            oldLine.Segments = segments
            newLine.Segments = segments
            
            i++ // Skip the next line since we processed it
        }
    }
}
```

**Key Points**:
- Only processes adjacent removed/added line pairs
- Uses semantic diff cleanup for better results
- Segments track byte positions, not character positions
- Both lines share the same segments for consistent highlighting

### 3. Syntax Highlighting Integration

OpenCode integrates Chroma for syntax highlighting with a sophisticated theme mapping system.

**Complete Chroma Theme Mapping**:
```go
func SyntaxHighlight(w io.Writer, source, fileName, formatter string, bg color.Color) error {
    t := theme.CurrentTheme()
    
    // Determine language lexer
    l := lexers.Match(fileName)
    if l == nil {
        l = lexers.Analyse(source)
    }
    if l == nil {
        l = lexers.Fallback
    }
    
    // Generate Chroma theme XML from OpenCode theme
    syntaxThemeXml := fmt.Sprintf(`
    <style name="opencode-theme">
        <!-- Base -->
        <entry type="Background" style="bg:%s"/>
        <entry type="Text" style="%s"/>
        <entry type="Error" style="%s"/>
        
        <!-- Keywords (all map to syntaxKeyword) -->
        <entry type="Keyword" style="%s"/>
        <entry type="KeywordConstant" style="%s"/>
        <entry type="KeywordDeclaration" style="%s"/>
        <entry type="KeywordNamespace" style="%s"/>
        <entry type="KeywordType" style="%s"/>
        
        <!-- Names -->
        <entry type="NameFunction" style="%s"/>
        <entry type="NameClass" style="%s"/>
        <entry type="NameVariable" style="%s"/>
        
        <!-- Literals -->
        <entry type="LiteralString" style="%s"/>
        <entry type="LiteralNumber" style="%s"/>
        
        <!-- Comments -->
        <entry type="Comment" style="%s"/>
        
        <!-- Operators -->
        <entry type="Operator" style="%s"/>
        <entry type="Punctuation" style="%s"/>
    </style>`,
        getChromaColor(t.BackgroundPanel()),
        getChromaColor(t.Text()),
        getChromaColor(t.Error()),
        getChromaColor(t.SyntaxKeyword()),
        getChromaColor(t.SyntaxKeyword()),
        getChromaColor(t.SyntaxKeyword()),
        getChromaColor(t.SyntaxKeyword()),
        getChromaColor(t.SyntaxType()),
        getChromaColor(t.SyntaxFunction()),
        getChromaColor(t.SyntaxType()),
        getChromaColor(t.SyntaxVariable()),
        getChromaColor(t.SyntaxString()),
        getChromaColor(t.SyntaxNumber()),
        getChromaColor(t.SyntaxComment()),
        getChromaColor(t.SyntaxOperator()),
        getChromaColor(t.SyntaxPunctuation()),
    )
    
    // Create and apply the style
    style := chroma.MustNewXMLStyle(strings.NewReader(syntaxThemeXml))
    
    // Tokenize and format
    it, _ := l.Tokenise(nil, source)
    return formatters.Get(formatter).Format(w, style, it)
}
```

**Key Mappings**:
- All keyword types → `syntaxKeyword` color
- Function/method names → `syntaxFunction` color
- Classes/types → `syntaxType` color
- Variables → `syntaxVariable` color
- Strings → `syntaxString` color
- Numbers → `syntaxNumber` color
- Comments → `syntaxComment` color

### 4. Color Scheme

The theme system provides extensive diff-related colors:
- `DiffAdded/Removed/Context`: Text colors
- `DiffAddedBg/RemovedBg/ContextBg`: Background colors
- `DiffHighlightAdded/Removed`: Intra-line highlight colors
- `DiffLineNumber`: Line number styling
- `DiffAddedLineNumberBg/RemovedLineNumberBg`: Line number backgrounds

### 5. Rendering Pipeline

1. **Parse Unified Diff**: Convert diff text to structured data
   - Extract file names, line numbers, and content
   - Identify line types (added/removed/context)
2. **Highlight Intra-line Changes**: Find and mark character-level differences
   - Compare removed/added line pairs
   - Use diffmatchpatch to find exact changes
   - Create segments marking changed portions
3. **Apply Syntax Highlighting**: Use Chroma with dynamic theme
   - Generate Chroma theme from current OpenCode theme colors
   - Apply language-specific highlighting to each line
4. **Apply Diff Highlighting**: Add background colors and intra-line highlights
   - Layer diff backgrounds (green/red) over syntax highlighting
   - Apply stronger colors to intra-line changes
   - Preserve all ANSI sequences from syntax highlighting
5. **Render Line Numbers**: Format with appropriate styles
   - Different formatting for old vs new line numbers
   - Color-coded backgrounds for line number columns
6. **Parallel Processing**: Use goroutines for efficient multi-line rendering
   - Each line rendered independently
   - Results collected in order

### 6. Performance Optimizations

- **Parallel Rendering**: `WriteStringsPar` function processes lines concurrently
  - Uses goroutines to render each line independently
  - Maintains order while processing in parallel
  - Especially beneficial for large diffs
- **Pre-allocation**: Arrays pre-allocated with reasonable capacities
  - `make([]Hunk, 0, 10)` for hunks
  - `make([]DiffLine, 0, 10)` for lines
- **String Builder**: Efficient string concatenation
  - Pre-calculates buffer size: `sb.Grow(len(hunkCopy.Lines) * config.Width)`
- **ANSI Optimization**: Minimal escape sequences for better performance

## Key Technical Innovations

### 1. ANSI Escape Sequence Handling

The `applyHighlighting` function demonstrates sophisticated ANSI handling that preserves syntax highlighting.

**Complete Implementation**:
```go
func applyHighlighting(content string, segments []Segment, segmentType LineType, highlightBg compat.AdaptiveColor) string {
    // Find all ANSI sequences
    ansiRegex := regexp.MustCompile(`\x1b(?:[@-Z\\-_]|\[[0-9?]*(?:;[0-9?]*)*[@-~])`)
    ansiMatches := ansiRegex.FindAllStringIndex(content, -1)
    
    // Build visible position → byte position mapping
    visibleIdx := 0
    ansiSequences := make(map[int]string)
    lastAnsiSeq := "\x1b[0m" // Default reset
    
    for i := 0; i < len(content); {
        isAnsi := false
        
        // Check if current position is an ANSI sequence
        for _, match := range ansiMatches {
            if match[0] == i {
                ansiSequences[visibleIdx] = content[match[0]:match[1]]
                lastAnsiSeq = content[match[0]:match[1]]
                i = match[1]
                isAnsi = true
                break
            }
        }
        
        if !isAnsi {
            // Store last ANSI for this visible position
            if _, exists := ansiSequences[visibleIdx]; !exists {
                ansiSequences[visibleIdx] = lastAnsiSeq
            }
            visibleIdx++
            
            // Advance by UTF-8 rune
            _, size := utf8.DecodeRuneInString(content[i:])
            i += size
        }
    }
    
    // Apply highlighting
    var sb strings.Builder
    inSelection := false
    currentPos := 0
    
    // Get highlight colors
    bg := getColor(highlightBg)
    fg := getColor(theme.CurrentTheme().BackgroundPanel())
    
    for i := 0; i < len(content); {
        // Write ANSI sequences
        for _, match := range ansiMatches {
            if match[0] == i {
                sb.WriteString(content[match[0]:match[1]])
                i = match[1]
                goto nextChar
            }
        }
        
        // Check segment boundaries
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
        
        // Get character
        r, size := utf8.DecodeRuneInString(content[i:])
        
        if inSelection {
            // Apply highlight
            sb.WriteString(fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r>>8, g>>8, b>>8)) // fg
            sb.WriteString(fmt.Sprintf("\x1b[48;2;%d;%d;%dm", r>>8, g>>8, b>>8)) // bg
            sb.WriteRune(r)
            sb.WriteString("\x1b[0m") // Reset
            sb.WriteString(ansiSequences[currentPos]) // Restore
        } else {
            sb.WriteRune(r)
        }
        
        currentPos++
        i += size
        
    nextChar:
    }
    
    return sb.String()
}
```

**Example ANSI Sequences**:
- `\x1b[0m` - Reset all attributes
- `\x1b[31m` - Red foreground (ANSI-16)
- `\x1b[48;2;255;0;0m` - Red background (RGB)
- `\x1b[1;34m` - Bold blue
- `\x1b[38;5;214m` - Orange (256-color)

**Edge Cases Handled**:
1. Multiple ANSI sequences in a row
2. UTF-8 characters (emoji, CJK)
3. Incomplete ANSI sequences
4. Zero-width characters
5. Overlapping highlights

### 2. Color Handling

- Uses lipgloss color system for terminal colors
- Converts hex colors to appropriate terminal formats
- Handles RGB and ANSI color codes

### 3. View Modes

**Unified View**:
- Traditional diff format with +/- markers
- Full syntax highlighting preserved
- Character-level change highlighting

**Side-by-Side View**:
- Split screen showing old and new versions
- Intelligent line pairing for changes
- Synchronized scrolling

### Line Rendering Implementation

**Unified Line Rendering**:
```go
func renderUnifiedLine(fileName string, dl DiffLine, width int, t Theme) string {
    // Get styles based on line type
    var marker string
    var bgStyle Style
    var lineNumberStyle Style
    var highlightColor AdaptiveColor
    var lineNum string
    
    switch dl.Kind {
    case LineRemoved:
        marker = "-"
        bgStyle = NewStyle().Background(t.DiffRemovedBg())
        lineNumberStyle = NewStyle().
            Background(t.DiffRemovedLineNumberBg()).
            Foreground(t.DiffRemoved())
        highlightColor = t.DiffHighlightRemoved()
        lineNum = fmt.Sprintf("%6d       ", dl.OldLineNo)
        
    case LineAdded:
        marker = "+"
        bgStyle = NewStyle().Background(t.DiffAddedBg())
        lineNumberStyle = NewStyle().
            Background(t.DiffAddedLineNumberBg()).
            Foreground(t.DiffAdded())
        highlightColor = t.DiffHighlightAdded()
        lineNum = fmt.Sprintf("      %7d", dl.NewLineNo)
        
    case LineContext:
        marker = " "
        bgStyle = NewStyle().Background(t.DiffContextBg())
        lineNumberStyle = NewStyle().
            Background(t.DiffLineNumber()).
            Foreground(t.TextMuted())
        lineNum = fmt.Sprintf("%6d %6d", dl.OldLineNo, dl.NewLineNo)
    }
    
    // Render components
    prefix := renderLinePrefix(dl, lineNum, marker, lineNumberStyle, t)
    content := renderLineContent(fileName, dl, bgStyle, highlightColor, width-14)
    
    return prefix + content
}
```

**Side-by-Side Pairing Algorithm**:
```go
func pairLines(lines []DiffLine) []linePair {
    var pairs []linePair
    i := 0
    
    for i < len(lines) {
        switch lines[i].Kind {
        case LineRemoved:
            // Check if next line is addition (change)
            if i+1 < len(lines) && lines[i+1].Kind == LineAdded {
                pairs = append(pairs, linePair{
                    left:  &lines[i],
                    right: &lines[i+1],
                })
                i += 2
            } else {
                // Deletion only
                pairs = append(pairs, linePair{
                    left:  &lines[i],
                    right: nil,
                })
                i++
            }
            
        case LineAdded:
            // Addition only (not paired with removal)
            pairs = append(pairs, linePair{
                left:  nil,
                right: &lines[i],
            })
            i++
            
        case LineContext:
            // Same line on both sides
            pairs = append(pairs, linePair{
                left:  &lines[i],
                right: &lines[i],
            })
            i++
        }
    }
    
    return pairs
}
```

## Visual Elements

1. **Line Numbers**: 
   - Different styles for old/new lines
   - Color-coded backgrounds matching line types

2. **Markers**: 
   - `+` for added lines
   - `-` for removed lines
   - ` ` for context lines

3. **Background Colors**:
   - Green background for added lines
   - Red background for removed lines
   - Default background for context

4. **Intra-line Highlights**:
   - Stronger colors for exact changes within lines
   - Preserves syntax highlighting underneath

## Integration Points

The diff component is used in:
- Chat messages to show file edits
- File viewer for comparing versions
- Potentially in git operations

## Advantages Over Traditional Diff Tools

1. **Syntax Awareness**: Full syntax highlighting preserved in diffs
2. **Character Precision**: Exact character-level change highlighting
3. **Beautiful Themes**: 20 professionally designed color schemes
4. **Performance**: Parallel rendering for large diffs
5. **Flexibility**: Multiple view modes (unified and side-by-side)

## Theme System (Simplified)

OpenCode has 20 built-in themes that we'll copy directly:
- ayu, catppuccin, cobalt2, dracula, everforest, github, gruvbox, kanagawa, material, matrix, monokai, nord, one-dark, opencode, palenight, rosepine, solarized, synthwave84, tokyonight, zenburn

### Theme Structure
Each theme JSON contains:
- `defs`: Color definitions that can be referenced
- `theme`: Maps semantic colors to actual values

Example from dracula.json:
```json
{
  "defs": {
    "background": "#282a36",
    "green": "#50fa7b",
    "red": "#ff5555"
  },
  "theme": {
    "diffAdded": {"dark": "green", "light": "green"},
    "diffAddedBg": {"dark": "#1a3a1a", "light": "#e0ffe0"}
  }
}
```

### Key Colors to Extract
From each theme JSON, we'll extract:
- `diffAdded/Removed`: Text colors for diff lines
- `diffAddedBg/RemovedBg`: Background colors
- `diffHighlightAdded/Removed`: Character-level highlight colors
- `diffLineNumber`: Line number styling
- Basic syntax colors: `syntaxKeyword`, `syntaxString`, `syntaxComment`, etc.

## Unified Diff Format Parsing

Understanding the unified diff format is crucial for implementing the parser.

### Format Structure
```
diff --git a/file.go b/file.go
index abc123..def456 100644
--- a/file.go
+++ b/file.go
@@ -10,7 +10,7 @@ func main() {
 	context line
 	context line
-	removed line
+	added line
 	context line
 	more context
```

### Parser Implementation Details
```go
func ParseUnifiedDiff(diff string) (DiffResult, error) {
    var result DiffResult
    var currentHunk *Hunk
    
    scanner := bufio.NewScanner(strings.NewReader(diff))
    var oldLine, newLine int
    inFileHeader := true
    
    for scanner.Scan() {
        line := scanner.Text()
        
        // File headers
        if inFileHeader {
            if strings.HasPrefix(line, "--- a/") {
                result.OldFile = line[6:]
                continue
            }
            if strings.HasPrefix(line, "+++ b/") {
                result.NewFile = line[6:]
                inFileHeader = false
                continue
            }
        }
        
        // Hunk header: @@ -10,7 +10,7 @@ func main() {
        if strings.HasPrefix(line, "@@") {
            // Save previous hunk
            if currentHunk != nil {
                result.Hunks = append(result.Hunks, *currentHunk)
            }
            
            // Parse line numbers
            re := regexp.MustCompile(`@@ -(\d+)(?:,(\d+))? \+(\d+)(?:,(\d+))? @@`)
            matches := re.FindStringSubmatch(line)
            if matches != nil {
                oldLine, _ = strconv.Atoi(matches[1])
                newLine, _ = strconv.Atoi(matches[3])
            }
            
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
            var dl DiffLine
            
            switch line[0] {
            case '+':
                dl.Kind = LineAdded
                dl.NewLineNo = newLine
                dl.Content = line[1:]
                newLine++
                
            case '-':
                dl.Kind = LineRemoved
                dl.OldLineNo = oldLine
                dl.Content = line[1:]
                oldLine++
                
            case ' ':
                dl.Kind = LineContext
                dl.OldLineNo = oldLine
                dl.NewLineNo = newLine
                dl.Content = line[1:]
                oldLine++
                newLine++
                
            default:
                // Malformed line, treat as context
                dl.Kind = LineContext
                dl.OldLineNo = oldLine
                dl.NewLineNo = newLine
                dl.Content = line
                oldLine++
                newLine++
            }
            
            currentHunk.Lines = append(currentHunk.Lines, dl)
        }
    }
    
    // Don't forget the last hunk
    if currentHunk != nil {
        result.Hunks = append(result.Hunks, *currentHunk)
    }
    
    return result, scanner.Err()
}
```

## Building a Standalone Tool

To recreate this as a standalone git diff tool:

### 1. Core Diff Rendering
- **Copy `diff.go`** but remove Bubble Tea dependencies
- **Keep the essential functions**:
  - `ParseUnifiedDiff()` - Parse git diff output
  - `HighlightIntralineChanges()` - Character-level highlighting
  - `SyntaxHighlight()` - Keep the dynamic Chroma theme generation (it's essential!)
  - `RenderUnifiedHunk()` / `RenderSideBySideHunk()` - Display modes

### 2. Minimal Dependencies
- **Lipgloss**: Terminal styling (keep as-is)
- **Chroma**: Syntax highlighting (keep as-is)
- **go-diff/diffmatchpatch**: Character-level diffs (keep as-is)
- **Remove**: Bubble Tea, Glamour, and other UI framework dependencies

### 3. Theme System Implementation
- **Copy theme JSON files** from `/packages/tui/internal/theme/themes/`
- **Theme loader features**:
  - Parse JSON with color reference resolution
  - Support both dark and light variants
  - Handle color references (e.g., "green" → "#50fa7b")
- **Terminal background detection**:
  - Use existing terminal info detection
  - Select dark/light variant based on terminal background
- **Embed themes** in binary using `go:embed`
- **Chroma theme generation**:
  - Keep the XML generation that maps theme colors to Chroma syntax elements
  - This ensures syntax highlighting matches the selected theme

### 4. CLI Tool Structure
```
git-diff-pretty/
├── main.go           # CLI entry point
├── diff.go           # Core diff rendering (from OpenCode)
├── themes.go         # Simple theme loader
├── themes/           # Copy of OpenCode theme JSONs
│   ├── dracula.json
│   ├── monokai.json
│   └── ...
└── go.mod
```

### 5. Command Usage
```bash
# Basic usage
git diff | git-diff-pretty

# With theme selection
git diff | git-diff-pretty --theme dracula

# Side-by-side view
git diff | git-diff-pretty --split

# Direct git integration
git-diff-pretty HEAD~1 HEAD --theme monokai
```

### 6. Implementation Steps

1. **Extract Core Components**:
   - Copy `diff.go` and remove Bubble Tea imports/dependencies
   - Copy theme-related utilities (color handling, theme interface)
   - Create simplified versions of style utilities

2. **Implement Theme Loading**:
   - Create theme parser that handles JSON structure
   - Implement color reference resolution (defs → theme colors)
   - Add terminal background detection for dark/light selection

3. **Create CLI Interface**:
   - Parse command-line arguments (theme, view mode, etc.)
   - Handle both piped input and direct git command execution
   - Implement proper terminal width detection

4. **Integration Details**:
   - Ensure proper ANSI color support detection
   - Handle edge cases (empty diffs, binary files)
   - Add pager support for large diffs

5. **Testing Strategy**:
   - Test with various programming languages for syntax highlighting
   - Verify ANSI sequence preservation
   - Ensure character-level highlighting works correctly
   - Test all 20 themes in both dark and light terminals

### Key Implementation Challenges

1. **ANSI Sequence Handling**: The most complex part - must perfectly preserve syntax highlighting while adding diff colors
2. **Terminal Detection**: Need reliable terminal background detection for theme selection
3. **Width Calculation**: Proper handling of terminal width for side-by-side view
4. **Performance**: Maintaining parallel rendering performance without Bubble Tea framework

The result will be a standalone tool that brings OpenCode's beautiful diff rendering to any git repository, with full theme support and syntax highlighting.

## Complete Function Reference

### Core Functions

```go
// Parse unified diff format into structured data
func ParseUnifiedDiff(diff string) (DiffResult, error)

// Add character-level highlighting to a hunk
func HighlightIntralineChanges(h *Hunk)

// Apply syntax highlighting to a line
func SyntaxHighlight(w io.Writer, source, fileName, formatter string, bg color.Color) error

// Render a hunk in unified format
func RenderUnifiedHunk(fileName string, h Hunk, opts ...UnifiedOption) string

// Render a hunk in side-by-side format
func RenderSideBySideHunk(fileName string, h Hunk, opts ...UnifiedOption) string

// Format entire diff in unified view
func FormatUnifiedDiff(filename string, diffText string, opts ...UnifiedOption) (string, error)

// Format entire diff in side-by-side view
func FormatDiff(filename string, diffText string, opts ...UnifiedOption) (string, error)
```

### Helper Functions

```go
// Render multiple strings in parallel while maintaining order
func WriteStringsPar[T any](sb *strings.Builder, items []T, fn func(T) string)

// Convert adaptive color to hex string based on terminal background
func AdaptiveColorToString(color compat.AdaptiveColor) *string

// Get terminal info (background color, dark/light)
func DetectTerminalBackground() TerminalInfo

// Create line number prefix with proper styling
func renderLinePrefix(dl DiffLine, lineNum string, marker string, style Style, t Theme) string

// Render line content with all highlighting layers
func renderLineContent(fileName string, dl DiffLine, bgStyle Style, highlightColor AdaptiveColor, width int) string
```

### Usage Examples

**Basic Diff Rendering**:
```go
// Read git diff
diffText := "..." // from git diff command

// Parse and render
result, err := ParseUnifiedDiff(diffText)
if err != nil {
    return err
}

// Apply intra-line highlighting
for i := range result.Hunks {
    HighlightIntralineChanges(&result.Hunks[i])
}

// Render unified view
output := FormatUnifiedDiff("main.go", diffText, WithWidth(120))
fmt.Print(output)
```

**Side-by-Side View**:
```go
// Configure for wide terminal
config := NewSideBySideConfig(WithWidth(160))

// Render each hunk
for _, hunk := range result.Hunks {
    output := RenderSideBySideHunk("main.go", hunk, WithWidth(160))
    fmt.Print(output)
}
```

**Theme Integration**:
```go
// Load theme
theme := LoadTheme("dracula")
SetCurrentTheme(theme)

// Theme colors automatically used in rendering
output := FormatUnifiedDiff("main.go", diffText)
```

## Terminal Compatibility

### Background Detection

```go
type TerminalInfo struct {
    Background       color.Color
    BackgroundIsDark bool
}

func DetectTerminalBackground() TerminalInfo {
    // Try OSC 11 query
    fmt.Print("\x1b]11;?\x07")
    
    // Read response with timeout
    response := readWithTimeout(100 * time.Millisecond)
    
    // Parse RGB values from response
    if matches := regexp.MustCompile(`rgb:([0-9a-f]+)/([0-9a-f]+)/([0-9a-f]+)`).FindStringSubmatch(response); matches != nil {
        r, _ := strconv.ParseInt(matches[1][:2], 16, 64)
        g, _ := strconv.ParseInt(matches[2][:2], 16, 64)
        b, _ := strconv.ParseInt(matches[3][:2], 16, 64)
        
        bg := color.RGBA{uint8(r), uint8(g), uint8(b), 255}
        isDark := (r*299 + g*587 + b*114) / 1000 < 128
        
        return TerminalInfo{
            Background:       bg,
            BackgroundIsDark: isDark,
        }
    }
    
    // Fallback to dark background
    return TerminalInfo{
        Background:       color.Black,
        BackgroundIsDark: true,
    }
}
```

### Color Support Detection

```go
func DetectColorSupport() ColorLevel {
    // Check COLORTERM
    if os.Getenv("COLORTERM") == "truecolor" {
        return TrueColor // 16M colors
    }
    
    // Check TERM
    term := os.Getenv("TERM")
    if strings.Contains(term, "256color") {
        return Color256 // 256 colors
    }
    
    if strings.Contains(term, "color") {
        return Color16 // 16 colors
    }
    
    return NoColor
}
```

## Edge Cases and Error Handling

### 1. Binary Files
```go
if strings.Contains(diffText, "Binary files") {
    return "Binary file diff not shown", nil
}
```

### 2. Large Files
```go
const maxLineLength = 1000
if len(line) > maxLineLength {
    line = line[:maxLineLength] + "..."
}
```

### 3. Malformed Diffs
```go
// Handle missing headers
if !strings.HasPrefix(line, "@@") {
    continue
}

// Validate line numbers
if oldLine < 0 || newLine < 0 {
    return fmt.Errorf("invalid line numbers")
}
```

### 4. Unicode Handling
```go
// Always use rune-aware operations
for i, r := range content {
    // Process rune r at position i
}

// Correct width calculation
width := ansi.StringWidth(content)
```

### 5. Performance Considerations
```go
// Pre-allocate slices
hunks := make([]Hunk, 0, estimatedHunks)

// Use string builder
var sb strings.Builder
sb.Grow(estimatedSize)

// Parallel processing threshold
if len(lines) > 100 {
    WriteStringsPar(&sb, lines, renderFunc)
} else {
    // Sequential for small diffs
}
```

## Testing Strategies

### 1. ANSI Sequence Preservation
```go
func TestANSIPreservation(t *testing.T) {
    input := "\x1b[31mred\x1b[0m text"
    segments := []Segment{{Start: 0, End: 3, Type: LineAdded}}
    
    output := applyHighlighting(input, segments, LineAdded, theme.DiffHighlightAdded())
    
    // Verify original ANSI sequences are preserved
    assert.Contains(t, output, "\x1b[31m")
}
```

### 2. Character-Level Diff Accuracy
```go
func TestIntralineDiff(t *testing.T) {
    oldLine := "hello world"
    newLine := "hello beautiful world"
    
    // Should highlight "beautiful " as added
    segments := computeSegments(oldLine, newLine)
    assert.Equal(t, 1, len(segments))
    assert.Equal(t, "beautiful ", segments[0].Text)
}
```

### 3. Theme Color Application
```go
func TestThemeColors(t *testing.T) {
    SetTheme("dracula")
    output := RenderUnifiedLine("test.go", addedLine, 80, theme)
    
    // Verify Dracula's green is used
    assert.Contains(t, output, "#50fa7b")
}
```

## Git Integration and CLI Implementation

### Main CLI Structure
```go
package main

import (
    "flag"
    "fmt"
    "io"
    "os"
    "os/exec"
)

func main() {
    var (
        theme    = flag.String("theme", "opencode", "Color theme to use")
        split    = flag.Bool("split", false, "Use side-by-side view")
        width    = flag.Int("width", 0, "Terminal width (auto-detect if 0)")
        pager    = flag.Bool("pager", true, "Use pager for output")
        listThemes = flag.Bool("list-themes", false, "List available themes")
    )
    flag.Parse()
    
    if *listThemes {
        for _, name := range AvailableThemes() {
            fmt.Println(name)
        }
        return
    }
    
    // Load themes
    if err := LoadThemesFromJSON(); err != nil {
        fmt.Fprintf(os.Stderr, "Failed to load themes: %v\n", err)
        os.Exit(1)
    }
    
    // Set theme
    if err := SetTheme(*theme); err != nil {
        fmt.Fprintf(os.Stderr, "Invalid theme: %v\n", err)
        os.Exit(1)
    }
    
    // Get diff input
    var diffText string
    var err error
    
    if flag.NArg() > 0 {
        // Run git diff with provided args
        diffText, err = runGitDiff(flag.Args())
    } else {
        // Read from stdin
        diffText, err = readStdin()
    }
    
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
    
    // Detect terminal info
    termInfo := DetectTerminalBackground()
    Terminal.BackgroundIsDark = termInfo.BackgroundIsDark
    
    // Auto-detect width if needed
    if *width == 0 {
        *width = getTerminalWidth()
    }
    
    // Render diff
    var output string
    if *split {
        output, err = FormatDiff("", diffText, WithWidth(*width))
    } else {
        output, err = FormatUnifiedDiff("", diffText, WithWidth(*width))
    }
    
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to render diff: %v\n", err)
        os.Exit(1)
    }
    
    // Output with optional pager
    if *pager && isTerminal(os.Stdout) {
        showWithPager(output)
    } else {
        fmt.Print(output)
    }
}

func runGitDiff(args []string) (string, error) {
    cmd := exec.Command("git", append([]string{"diff", "--no-color"}, args...)...)
    output, err := cmd.Output()
    if err != nil {
        return "", fmt.Errorf("git diff failed: %w", err)
    }
    return string(output), nil
}

func readStdin() (string, error) {
    data, err := io.ReadAll(os.Stdin)
    if err != nil {
        return "", fmt.Errorf("failed to read stdin: %w", err)
    }
    return string(data), nil
}

func getTerminalWidth() int {
    cmd := exec.Command("tput", "cols")
    output, err := cmd.Output()
    if err != nil {
        return 80 // Default
    }
    
    width := 0
    fmt.Sscanf(string(output), "%d", &width)
    if width <= 0 {
        return 80
    }
    return width
}

func showWithPager(content string) {
    // Try common pagers
    pagers := []string{"less", "more"}
    
    for _, pager := range pagers {
        if _, err := exec.LookPath(pager); err == nil {
            cmd := exec.Command(pager, "-R") // -R for ANSI colors
            cmd.Stdin = strings.NewReader(content)
            cmd.Stdout = os.Stdout
            cmd.Stderr = os.Stderr
            
            if err := cmd.Run(); err == nil {
                return
            }
        }
    }
    
    // Fallback to direct output
    fmt.Print(content)
}

func isTerminal(f *os.File) bool {
    fi, _ := f.Stat()
    return fi.Mode()&os.ModeCharDevice != 0
}
```

### Build and Installation
```bash
# Build
go build -o git-diff-pretty .

# Install globally
sudo cp git-diff-pretty /usr/local/bin/

# Or add to PATH
export PATH=$PATH:$(pwd)
```

### Usage in Git Config
```bash
# Set as default diff tool
git config --global core.pager "git-diff-pretty --theme dracula"

# Or create an alias
git config --global alias.dp "!git diff | git-diff-pretty"
```

## Summary

This document provides everything needed to implement OpenCode's beautiful diff rendering as a standalone tool:

1. **Complete data structures** with relationships and flow
2. **Detailed algorithms** for parsing, highlighting, and rendering
3. **ANSI handling logic** with full implementation
4. **Theme system** with color mappings and Chroma integration
5. **Function signatures** and usage examples
6. **Terminal compatibility** code
7. **Edge case handling** and performance tips
8. **Complete CLI implementation** with Git integration

The key to OpenCode's beauty is the careful layering of syntax highlighting, diff backgrounds, and character-level highlights, all while preserving ANSI sequences perfectly. With this guide, a developer can recreate the full experience.