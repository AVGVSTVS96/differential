# Pretty-Diff

> [!WARNING]
> **Early Alpha - Vibe-Coded Proof of Concept**
> 
> This project was rapidly prototyped as a proof of concept and is in a very early alpha stage. It has only been tested for basic functionality and has not undergone thorough human code review. Use at your own risk and expect bugs, missing features, and potential breaking changes.

A beautiful terminal UI (TUI) for viewing diffs with syntax highlighting, character-level changes, and interactive navigation. Inspired by [delta](https://github.com/dandavison/delta) and based on OpenCode's diff implementation.

## Features

- 🎨 **Syntax Highlighting** - Full language-aware syntax highlighting using Chroma
- 🔍 **Character-level Diffs** - See exactly what changed within lines
- 🎯 **Intra-line Highlighting** - Smart highlighting of changed segments
- 🌈 **Multiple Themes** - 8 built-in themes with terminal background detection
- 📊 **Side-by-Side View** - Compare files in split view or unified view
- ⚡ **Fast Performance** - Parallel rendering with ANSI sequence preservation
- 🖥️ **Interactive TUI** - Navigate diffs with vim-like keybindings
- 🔧 **Git Integration** - Drop-in replacement for `git diff`

## Installation

```bash
# Clone the repository
git clone https://github.com/pretty-diff/pretty-diff
cd pretty-diff

# Build
go build ./cmd/pretty-diff

# Optional: Install to PATH
go install ./cmd/pretty-diff
```

## Usage

### Basic Usage

```bash
# Compare two files
pretty-diff file1.go file2.go

# Pipe git diff output
git diff | pretty-diff

# Compare git commits/branches
pretty-diff HEAD~3 HEAD
pretty-diff main feature-branch

# Compare directories
pretty-diff dir1/ dir2/
```

### Pipe Mode (Non-Interactive)

For scripting or when you want static output:

```bash
# Force pipe mode with --pipe-mode or -p
pretty-diff file1.go file2.go --pipe-mode

# Disable pager for continuous output
pretty-diff file1.go file2.go --pipe-mode --no-pager
```

### Themes

```bash
# List available themes
pretty-diff --list-themes

# Set theme with -t or --theme
pretty-diff file1.go file2.go -t monokai
pretty-diff file1.go file2.go --theme github

# Available themes:
# - dracula (default)
# - monokai
# - nord
# - github
# - gruvbox
# - catppuccin
# - tokyonight
# - solarized
```

### View Modes

```bash
# Side-by-side view with -s or --side-by-side
pretty-diff file1.go file2.go -s
pretty-diff file1.go file2.go --side-by-side

# In TUI mode, press Tab to toggle between unified and side-by-side
```

### Line Numbers and Context

```bash
# Hide line numbers
pretty-diff file1.go file2.go --no-line-numbers

# Change context lines (default: 3)
pretty-diff file1.go file2.go -c 5
pretty-diff file1.go file2.go --context 10
```

### Combining Options

```bash
# Side-by-side with Nord theme and 10 context lines
pretty-diff file1.go file2.go -s -t nord -c 10

# Pipe mode with GitHub theme and no line numbers
git diff | pretty-diff --pipe-mode -t github --no-line-numbers
```

## Interactive TUI Mode

When running without `--pipe-mode`, pretty-diff launches an interactive terminal UI:

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `j` / `↓` | Scroll down |
| `k` / `↑` | Scroll up |
| `g` / `Home` | Go to top |
| `G` / `End` | Go to bottom |
| `Ctrl+f` / `PgDn` | Page down |
| `Ctrl+b` / `PgUp` | Page up |
| `Tab` | Toggle unified/side-by-side view |
| `n` | Toggle line numbers |
| `?` | Show help |
| `q` / `Ctrl+c` | Quit |

### Navigation Features

- Smooth scrolling through large diffs
- Jump between hunks with `{` and `}` (coming soon)
- Search within diffs with `/` (coming soon)

## Configuration

Pretty-diff can be configured via a TOML file at `~/.config/pretty-diff/config.toml`:

```toml
[ui]
theme = "dracula"
default_view = "unified"  # or "side-by-side"
tab_width = 4
line_numbers = true
syntax_highlight = true
wrap_lines = false

[git]
default_context = 3
ignore_whitespace = false
show_stats = true

[keybindings]
quit = "q"
help = "?"
toggle_view = "tab"
next_hunk = "}"
prev_hunk = "{"
scroll_up = "k"
scroll_down = "j"
```

## Git Integration

Pretty-diff can be used as a drop-in replacement for git diff:

```bash
# Set as git pager
git config --global pager.diff "pretty-diff --pipe-mode"
git config --global pager.show "pretty-diff --pipe-mode"

# Or use directly
git diff | pretty-diff
git show HEAD | pretty-diff
```

## Examples

### Viewing Code Changes

```bash
# Review your staged changes with syntax highlighting
git diff --cached | pretty-diff -t github

# Compare feature branch with main in side-by-side view
pretty-diff main feature-branch -s

# Review a specific commit
git show abc123 | pretty-diff
```

### Different File Types

Pretty-diff automatically detects file types and applies appropriate syntax highlighting:

```bash
# Python files
pretty-diff old.py new.py

# JavaScript/TypeScript
pretty-diff src/old.ts src/new.ts -t monokai

# Markdown with GitHub theme
pretty-diff README.old.md README.new.md -t github

# JSON with pretty side-by-side view
pretty-diff config.old.json config.new.json -s
```

## Tips

1. **Terminal Colors**: Pretty-diff automatically detects if your terminal has a dark or light background and adjusts themes accordingly.

2. **Large Files**: For better performance with large files, use `--context` to limit the amount of context shown:
   ```bash
   pretty-diff large1.txt large2.txt -c 1
   ```

3. **Piping to Less**: If you want to use your own pager:
   ```bash
   pretty-diff file1 file2 --pipe-mode --no-pager | less -R
   ```

4. **Shell Aliases**: Add convenient aliases to your shell config:
   ```bash
   alias pd="pretty-diff"
   alias gd="git diff | pretty-diff --pipe-mode"
   ```

## Architecture

Pretty-diff is built with:
- **Bubble Tea** - Terminal UI framework
- **Lipgloss** - Styling and layout
- **Chroma** - Syntax highlighting
- **go-diff** - Diff algorithms

The diff rendering engine is based on OpenCode's implementation, featuring:
- Sophisticated ANSI escape sequence handling
- Character-level diff computation with proper UTF-8 support
- Parallel rendering for performance
- Theme system with dynamic Chroma style generation

## License

MIT License - see LICENSE file for details

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Acknowledgments

- OpenCode team for the excellent diff rendering implementation
- Delta for the inspiration and CLI design
- Charm.sh for the amazing TUI libraries