// Package diff provides diff parsing, rendering, and formatting with syntax highlighting.
//
// This package handles the core diff functionality including:
//   - Parsing unified diff format (ParseUnifiedDiff in parser.go)
//   - Character-level diff highlighting (HighlightIntralineChanges in highlighter.go)
//   - Syntax highlighting integration (SyntaxHighlight in syntax.go)
//   - Unified view rendering (RenderUnifiedHunk, FormatUnifiedDiff in unified.go)
//   - Side-by-side view rendering (RenderSideBySideHunk, FormatDiff in sidebyside.go)
//
// High-level API:
//   - FormatUnifiedDiff(filename, diffText, config, opts...) - Full unified diff
//   - FormatDiff(filename, diffText, config, opts...) - Full side-by-side diff
//
// Architecture:
//   - types.go: Core data structures (LineType, DiffLine, Hunk, etc.)
//   - parser.go: Unified diff parsing
//   - highlighter.go: Intra-line diff highlighting
//   - syntax.go: Syntax highlighting integration with Chroma
//   - render.go: Shared rendering primitives
//   - unified.go: Unified view rendering
//   - sidebyside.go: Side-by-side view rendering
//
// This package was refactored from a single 976-line file into focused modules
// for better maintainability and testability.
package diff
