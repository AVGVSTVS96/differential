# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Differential is a terminal UI (TUI) for viewing diffs with syntax highlighting and character-level changes. It's inspired by delta and based on OpenCode's diff implementation. This is an early alpha proof-of-concept that hasn't undergone thorough human review.

## Build and Development Commands

```bash
# Build the binary
go build ./cmd/differential

# Run tests (no tests implemented yet)
go test ./...

# Install dependencies
go mod tidy

# Run in pipe mode (non-interactive)
./differential file1.go file2.go --pipe-mode --no-pager

# Run in TUI mode (interactive - requires terminal)
./differential file1.go file2.go
```

## Architecture Overview

### Core Components

1. **Diff Engine** (`internal/diff/`)
   - `parser.go`: Parses unified diff format, handles both git-style (`--- a/file`) and standard (`--- file`) formats
   - `highlighter.go`: Implements character-level diff highlighting using diffmatchpatch, preserves ANSI sequences
   - `renderer.go`: Renders diffs in unified/side-by-side views with syntax highlighting and parallel processing
   - `types.go`: Core data structures (DiffLine, Hunk, DiffResult, Segment)

2. **Theme System** (`internal/themes/`)
   - `theme.go`: Theme loader and color resolver, detects terminal background (dark/light)
   - `chroma.go`: Generates dynamic Chroma styles from theme colors for syntax highlighting
   - `embedded.go`: Loads 8 built-in themes from JSON files
   - `themes/*.json`: Theme definitions with dark/light variants

3. **TUI Application** (`internal/app/`)
   - `app.go`: Main application logic with two modes:
     - Pipe mode: Static output for scripting/piping
     - TUI mode: Interactive with Bubble Tea framework
   - Handles git integration via exec commands

4. **CLI Interface** (`cmd/differential/`)
   - `main.go`: Cobra-based CLI with flags for themes, view modes, context lines

### Key Technical Details

- **ANSI Sequence Preservation**: The highlighter carefully preserves existing ANSI sequences while applying diff colors by mapping visible positions to byte positions
- **Parallel Rendering**: Diff lines are rendered in parallel using goroutines for performance
- **Character-level Diffs**: Uses diffmatchpatch to compute intra-line changes and highlights specific segments
- **Theme System**: Supports both dark and light terminal backgrounds with automatic detection
- **Git Integration**: Can parse both piped input and run git/diff commands directly

### Important Implementation Notes

1. The regex patterns in `parser.go` support both git-style diffs (`--- a/file`) and standard diffs (`--- file timestamp`)
2. Terminal width detection uses `tput cols` command
3. The TUI requires a proper terminal and will fail with "could not open a new TTY" in non-terminal environments
4. Syntax highlighting is only applied to context lines; added/removed lines use diff colors
5. The renderer handles UTF-8 and zero-width characters correctly

## Common Issues and Solutions

- **No output in pipe mode**: Check if the diff is being parsed correctly - the parser expects proper unified diff format
- **TTY errors**: Use `--pipe-mode` flag when running in non-terminal environments
- **Theme not loading**: The theme system will fall back to a default theme if initialization fails