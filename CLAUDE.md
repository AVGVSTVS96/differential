# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Differential is a terminal diff viewer built on OpenTUI's built-in `<diff>` component. It provides syntax-highlighted diff viewing with both interactive TUI mode and pipe mode for scripting.

## Build and Development Commands

```bash
# Install dependencies
bun install

# Run the CLI
bun run src/index.tsx --help

# Run with arguments
bun run src/index.tsx file1.ts file2.ts
bun run src/index.tsx --view split --theme dracula file1.ts file2.ts

# Pipe mode (non-interactive)
bun run src/index.tsx --pipe file1.ts file2.ts

# View git changes
bun run src/index.tsx --staged
bun run src/index.tsx --commit HEAD~1
```

## Architecture Overview

### Project Structure

```
src/
  index.tsx   # CLI entrypoint - argument parsing and renderer setup
  app.tsx     # Main React component - file navigation, keyboard handling, layout
  theme.ts    # Theme definitions (colors for diff, syntax, UI)
  utils.ts    # Helper functions (language detection, ANSI colors)
test/
  *.ts        # Test files for manual testing
```

### Core Dependencies

- **@opentui/core** - TUI framework providing CliRenderer and core renderables
- **@opentui/react** - React reconciler for OpenTUI
- **diff** - Unified diff parsing (parsePatch)
- **cac** - CLI argument parsing
- **react** - UI framework

### Key Components

1. **CLI (index.tsx)**
   - Parses arguments with cac
   - Reads diff from files, stdin, or git commands
   - Creates CliRenderer and renders React app
   - Handles pipe mode for non-interactive output

2. **App (app.tsx)**
   - Uses OpenTUI's built-in `<diff>` component for rendering
   - Manages state: current file, theme, view mode, line numbers
   - Handles keyboard navigation (file switching, view toggling, theme cycling)
   - Parses multi-file diffs and allows navigation between them

3. **OpenTUI's `<diff>` Component**
   - Handles unified and split view rendering
   - Tree-sitter syntax highlighting
   - Word-level diff highlighting
   - Line numbers with diff-aware coloring
   - Scrolling and viewport management

### Themes (theme.ts)

Six built-in themes:
- `github-dark` (default)
- `dracula`
- `monokai`
- `one-dark`
- `nord`
- `tokyonight`

Each theme specifies:
- Background colors (main, panel)
- Text colors (normal, muted)
- Diff colors (added, removed, context backgrounds and foregrounds)
- Line number colors
- UI colors (border, accent, hunk headers)

### Input Modes

1. **Two files**: `differential file1 file2` - generates diff using `diff -u`
2. **Single file**: `differential file` - runs `git diff` on the file
3. **Stdin**: `git diff | differential` - reads unified diff from pipe
4. **Git flags**: `--staged` for cached, `--commit <ref>` for specific commits

### Interactive Keybindings

- `q` / `Esc` - Quit
- `v` - Toggle unified/split view
- `t` - Cycle themes
- `L` - Toggle line numbers
- `n` / `p` / `←` / `→` - Navigate files

## Important Implementation Notes

1. **Built-in diff component** - Uses OpenTUI's `<diff>` component which handles all the complex diff rendering, syntax highlighting, and word-level diffs automatically.

2. **OpenTUI requires a TTY** - Interactive mode needs a real terminal. Use `--pipe` for scripts/CI.

3. **Automatic pipe mode** - If stdout is not a TTY, pipe mode is enabled automatically.

4. **Multi-file parsing** - Uses the `diff` library's `parsePatch()` to split multi-file diffs and enable file navigation.

5. **SyntaxStyle** - Themes define colors that are converted to OpenTUI's SyntaxStyle format for syntax highlighting.
