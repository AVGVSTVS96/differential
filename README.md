# Differential

> [!WARNING]
> **Vibe-Coded Proof of Concept**
>
> This is an early alpha proof-of-concept that has not been human reviewed. Use at your own risk.

> [!NOTE]
> **Complete Rewrite - Now Based on OpenCode**
>
> This project has been completely rewritten to use [OpenCode](https://opencode.ai/)'s battle-tested diff rendering engine. The implementation now exactly matches OpenCode's visual output with all 22 themes, proper ANSI handling, and sophisticated character-level diff highlighting.

A beautiful terminal diff viewer with syntax highlighting, character-level changes, and 22 professionally designed themes. Built on OpenCode's proven diff rendering engine and designed as a delta-like CLI tool.

## Features

- 🎨 **Syntax Highlighting** - Full language-aware syntax highlighting using Chroma
- 🔍 **Character-level Diffs** - See exactly what changed within lines
- 🎯 **Intra-line Highlighting** - Sophisticated highlighting of changed segments
- 🌈 **22 Beautiful Themes** - All themes from OpenCode with dark/light terminal support
- 📊 **Side-by-Side View** - Compare files in split view or unified view
- ⚡ **Fast Performance** - Parallel rendering with ANSI sequence preservation
- 🔧 **Git Integration** - Drop-in replacement for `git diff`
- 🎭 **Perfect ANSI Handling** - Preserves syntax highlighting while applying diff colors

## Installation

```bash
# Clone the repository
git clone https://github.com/avgvstvs96/differential
cd differential

# Build
go build ./cmd/differential

# Optional: Install to PATH
go install ./cmd/differential
```

## Usage

### Basic Usage

```bash
# Compare two files
differential file1.go file2.go

# Pipe git diff output
git diff | differential

# Compare git commits/branches
differential HEAD~3 HEAD
differential main feature-branch

# Compare directories
differential dir1/ dir2/
```

### Standard Output

Differential outputs directly to stdout, making it perfect for piping:

```bash
# Pipe to less for paging
differential file1.go file2.go | less -R

# Save to file
differential file1.go file2.go > diff.txt
```

### Themes

Differential includes all 22 themes from OpenCode:

```bash
# List available themes
differential --list-themes

# Set theme with -t or --theme
differential file1.go file2.go -t dracula
differential file1.go file2.go --theme nord

# Available themes:
# opencode (default), aura, ayu, catppuccin, cobalt2, dracula, everforest,
# github, gruvbox, kanagawa, material, matrix, monokai, nord, one-dark,
# palenight, rosepine, solarized, synthwave84, tokyonight, vesper, zenburn
```

### View Modes

```bash
# Side-by-side view with -s or --side-by-side
differential file1.go file2.go -s
differential file1.go file2.go --side-by-side

# Default is unified view
differential file1.go file2.go
```

### Width Control

```bash
# Auto-detect terminal width (default)
differential file1.go file2.go

# Set specific width
differential file1.go file2.go --width 120
differential file1.go file2.go -w 160
```

### Combining Options

```bash
# Side-by-side with Nord theme and custom width
differential file1.go file2.go -s -t nord -w 160

# Piped input with GitHub theme
git diff | differential -t github

# Side-by-side with Dracula theme
differential file1.go file2.go --side-by-side --theme dracula
```

## Configuration

Differential uses a simple CLI interface. For persistent theme preferences, you can create a shell alias:

```bash
# Add to ~/.bashrc or ~/.zshrc
alias diff='differential --theme dracula'
alias gdiff='git diff | differential --theme nord'
```

## Git Integration

Differential can be used as a drop-in replacement for git diff:

```bash
# Set as git pager
git config --global pager.diff differential
git config --global pager.show differential

# Or pipe git output
git diff | differential
git show HEAD | differential --theme dracula
```

## Examples

### Viewing Code Changes

```bash
# Review your staged changes with syntax highlighting
git diff --cached | differential -t github

# Compare feature branch with main in side-by-side view
differential main feature-branch -s

# Review a specific commit
git show abc123 | differential
```

### Different File Types

Differential automatically detects file types and applies appropriate syntax highlighting:

```bash
# Python files
differential old.py new.py

# JavaScript/TypeScript
differential src/old.ts src/new.ts -t monokai

# Markdown with GitHub theme
differential README.old.md README.new.md -t github

# JSON with pretty side-by-side view
differential config.old.json config.new.json -s
```

## Tips

1. **Terminal Colors**: Differential automatically detects if your terminal has a dark or light background and adjusts themes accordingly.

2. **Piping to Less**: Differential outputs directly to stdout, making it easy to pipe:
   ```bash
   differential file1 file2 | less -R
   ```

3. **Shell Aliases**: Add convenient aliases to your shell config:
   ```bash
   alias d="differential"
   alias gd="git diff | differential"
   alias gds="git diff | differential --side-by-side"
   ```

4. **Custom Width**: Adjust for your terminal size:
   ```bash
   differential file1 file2 --width 160  # For wide terminals
   differential file1 file2 --side-by-side --width 200
   ```

## Architecture

Differential is built directly on OpenCode's diff rendering engine with:
- **Lipgloss v2** - Terminal styling and colors
- **Chroma v2** - Syntax highlighting with dynamic theme generation
- **go-diff** - Character-level diff computation
- **OpenCode's Theme System** - 22 professional themes with color resolution

Key features from OpenCode's implementation:
- Sophisticated ANSI escape sequence preservation
- Character-level diff with proper UTF-8 and zero-width character support
- Parallel rendering using goroutines for performance
- Dynamic Chroma theme generation from JSON theme files
- Terminal background detection (dark/light mode)

## Roadmap

### Potential Enhancements
- 📝 **Context control** - Flag to adjust number of context lines shown
- 🎨 **Custom user themes** - Support for user-defined theme JSON files
- 🔗 **Hyperlink support** - Clickable file paths (if terminal supports)
- 📊 **Diff statistics** - Show file change summaries
- 🚀 **Performance optimizations** - Streaming parser for very large diffs
- 🌐 **Color output control** - Flags for forcing/disabling color output

### What's Working Now
- ✅ All 22 OpenCode themes with proper color application
- ✅ Character-level diff highlighting
- ✅ Syntax highlighting for 100+ languages
- ✅ Side-by-side and unified views
- ✅ Parallel rendering for performance
- ✅ Perfect ANSI sequence preservation

## License

MIT License - see LICENSE file for details

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Acknowledgments

- **[OpenCode](https://opencode.ai/)** - This project is built entirely on OpenCode's diff rendering engine. All credit for the sophisticated diff implementation, theme system, and ANSI handling goes to the OpenCode team.
- **[Delta](https://github.com/dandavison/delta)** - Inspiration for the delta-like CLI interface
- **[Charm.sh](https://charm.sh/)** - For the excellent Lipgloss and related libraries
