# Differential

A terminal diff viewer with syntax highlighting powered by [OpenTUI](https://github.com/sst/opentui)'s DiffRenderable component.

## Features

- **Syntax Highlighting** - Tree-sitter based syntax highlighting for 20+ languages
- **Split & Unified Views** - Toggle between side-by-side and unified diff views
- **Multiple Themes** - GitHub Dark, Dracula, Monokai, One Dark, Nord, Tokyo Night
- **Line Numbers** - Optional line number display
- **Word-level Diffs** - Highlights specific changed words within lines
- **Multi-file Navigation** - Navigate between files in multi-file diffs
- **Pipe Support** - Works with `git diff`, `git show`, and any unified diff input
- **Interactive TUI** - Full-screen terminal interface with keyboard navigation

## Installation

Requires [Bun](https://bun.sh) runtime.

```bash
# Clone the repository
git clone https://github.com/avgvstvs96/differential
cd differential

# Install dependencies
bun install

# Run directly
bun run src/index.tsx --help

# Or link globally
bun link
```

## Usage

```bash
# Compare two files
differential old.ts new.ts

# View git diff for a file
differential src/index.ts

# Pipe git diff output
git diff | differential
git show HEAD | differential

# View staged changes
differential --staged

# View specific commit
differential --commit abc123

# Side-by-side view
differential --view split old.ts new.ts

# Different theme
differential --theme dracula old.ts new.ts

# Pipe mode (non-interactive, colored output)
differential --pipe old.ts new.ts
git diff | differential -p
```

## Options

| Flag | Short | Description |
|------|-------|-------------|
| `--view <mode>` | `-v` | View mode: `unified` or `split` (default: unified) |
| `--theme <name>` | `-t` | Theme name (default: github-dark) |
| `--no-line-numbers` | `-n` | Hide line numbers |
| `--pipe` | `-p` | Output colored diff to stdout (non-interactive) |
| `--staged` | | Show staged changes (git diff --cached) |
| `--commit <ref>` | `-c` | Show changes for a specific commit |
| `--help` | `-h` | Show help message |
| `--version` | | Show version |

## Keybindings (Interactive Mode)

| Key | Action |
|-----|--------|
| `q`, `Esc` | Quit |
| `v` | Toggle unified/split view |
| `t` | Cycle through themes |
| `L` | Toggle line numbers |
| `n`, `→`, `l` | Next file (in multi-file diffs) |
| `p`, `←`, `h` | Previous file (in multi-file diffs) |
| `j`, `↓` | Scroll down |
| `k`, `↑` | Scroll up |

## Themes

- **github-dark** - GitHub's dark theme (default)
- **dracula** - Popular dark theme
- **monokai** - Classic Sublime Text theme
- **one-dark** - Atom One Dark theme
- **nord** - Arctic, north-bluish color palette
- **tokyonight** - Clean, dark theme inspired by Tokyo at night

## Examples

```bash
# Review staged changes
git diff --cached | differential

# Compare branches
git diff main..feature | differential --view split

# View a commit with Dracula theme
git show abc123 | differential --theme dracula

# Quick colored output for scripts
differential -p old.py new.py | less -R
```

## Architecture

Differential is a thin CLI wrapper around [OpenTUI](https://github.com/sst/opentui)'s built-in `<diff>` component. OpenTUI provides:

- High-performance terminal rendering via native Zig code
- Tree-sitter based syntax highlighting
- Word-level diff highlighting
- Split and unified view modes
- Flexible layout system (Yoga flexbox)

The implementation uses:
- **@opentui/core** - Core rendering engine
- **@opentui/react** - React reconciler for building the UI
- **diff** - Unified diff parsing
- **cac** - CLI argument parsing

## License

MIT License - see LICENSE file for details.

## Acknowledgments

- **[OpenTUI](https://github.com/sst/opentui)** - The TUI framework powering the diff rendering
- **[Critique](https://github.com/remorses/critique)** - Inspiration for multi-file navigation
- **[Delta](https://github.com/dandavison/delta)** - Inspiration for the CLI interface
