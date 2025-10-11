# Configuration Guide

Differential uses git config just like delta, making it a true drop-in replacement. You can configure differential through:

1. **Git Config** (recommended, delta-compatible)
2. **TOML Config File** (optional, for non-git usage)
3. **Command-line Flags** (highest priority)

## Git Config (Like Delta)

Differential reads configuration from your `.gitconfig` file, just like delta. This means:
- **Existing delta configs work automatically**
- No need to duplicate configuration
- Consistent experience across tools

### Using Delta Config Section

If you already have delta configured, differential will use those settings:

```gitconfig
[delta]
    side-by-side = true
    line-numbers = true
    navigate = true
    hyperlinks = true
    syntax-theme = Dracula
```

Differential will automatically read these settings!

### Using Differential Config Section

You can also create a `[differential]` section (takes precedence over `[delta]`):

```gitconfig
[differential]
    theme = dracula
    side-by-side = true
    line-numbers = true
    tabs = 4
    width = 160
```

### Setting Differential as Git Pager

Configure git to use differential just like delta:

```bash
# Set for all diff commands
git config --global pager.diff differential
git config --global pager.show differential
git config --global pager.log differential

# Set with options
git config --global pager.diff "differential --side-by-side"

# Or use core.pager for everything
git config --global core.pager "differential"
```

## Configuration Options

### Core Settings

| Option | Description | Default | Example |
|--------|-------------|---------|---------|
| `theme` | Color theme to use | `opencode` | `dracula`, `nord`, `github` |
| `side-by-side` | Show side-by-side diff | `false` | `true` |
| `line-numbers` | Display line numbers | `false` | `true` |
| `pager` | Custom pager command | auto-detect | `less -R` |
| `no-pager` | Disable pager | `false` | `true` |

### Display Settings

| Option | Description | Default | Example |
|--------|-------------|---------|---------|
| `width` | Terminal width | auto-detect | `120`, `160` |
| `tabs` | Tab width in spaces | `4` | `2`, `8` |
| `unified` | Context lines | `3` | `5`, `10` |
| `max-line-length` | Truncate long lines | `3000` | `500`, `unlimited` |
| `relative-paths` | Show relative paths | `false` | `true` |
| `keep-plus-minus-markers` | Show +/- prefix | `false` | `true` |

### Color Settings

| Option | Description | Default | Example |
|--------|-------------|---------|---------|
| `syntax-theme` | Syntax highlighting theme | `auto` | `monokai`, `github` |
| `color` | When to use colors | `auto` | `always`, `never` |
| `dark` | Use dark color scheme | auto-detect | `true` |
| `light` | Use light color scheme | auto-detect | `true` |

## TOML Config File (Optional)

For non-git usage, you can create `~/.config/differential/config.toml`:

```toml
[core]
theme = "dracula"
side-by-side = false
line-numbers = true

[display]
width = 0  # auto-detect
tabs = 4
unified = 3
max-line-length = 3000

[colors]
syntax-theme = "auto"
color = "auto"

[style]
minus-style = "normal auto"
plus-style = "syntax auto"
zero-style = "normal"
commit-style = "raw"
file-style = "blue"
hunk-header-style = "line-number syntax"
```

## Configuration Precedence

Settings are applied in this order (later overrides earlier):

1. Built-in defaults
2. Git config `[delta]` section
3. Git config `[differential]` section
4. TOML config file (`~/.config/differential/config.toml`)
5. Custom config file (`--config path/to/file.toml`)
6. Command-line flags (highest priority)

## Examples

### Example 1: Delta Replacement

If you have this delta config:

```gitconfig
[delta]
    side-by-side = true
    line-numbers = true
    navigate = true
    dark = true
    syntax-theme = Dracula
```

Just run differential - it works automatically:

```bash
git diff | differential
differential file1.go file2.go
```

### Example 2: Customize Per-Project

In your project's `.git/config`:

```gitconfig
[differential]
    theme = monokai
    width = 200
    tabs = 2
```

### Example 3: Override with Flags

```bash
# Use git config, but override theme
differential --theme nord file1.go file2.go

# Use git config, but force no pager
git diff | differential --no-pager

# Override context lines
differential -U10 file1.go file2.go
```

### Example 4: Disable Git Config

If you want to ignore git config entirely:

```bash
differential --no-gitconfig file1.go file2.go
```

## Migrating from Delta

1. **Keep your delta config** - differential reads it automatically
2. **Test side-by-side**: `git diff | differential`
3. **Optionally switch pager**: `git config --global pager.diff differential`
4. **Use both tools** - they can coexist perfectly

## Compatibility Notes

### Supported Delta Options

Most common delta options work in differential:

✅ `side-by-side`, `line-numbers`, `width`, `tabs`  
✅ `syntax-theme`, `dark`, `light`  
✅ `pager`, `navigate`, `hyperlinks` (support coming soon)  
✅ Style options: `minus-style`, `plus-style`, `zero-style`, etc.

### Delta Options Not Yet Supported

Some advanced delta features are planned but not yet implemented:

⏳ `navigate` - Interactive navigation (Phase 4)  
⏳ `hyperlinks` - Clickable links (Phase 4)  
⏳ `grep-*` styles - Grep output styling (Phase 5)  
⏳ `blame-*` styles - Blame output styling (Phase 5)  
⏳ `merge-conflict-*` - Merge conflict handling (Phase 5)

See our [implementation roadmap](CLAUDE.md) for details.

## Troubleshooting

### Config Not Loading

```bash
# Check if git config is readable
git config --get-regexp "^differential\\."
git config --get-regexp "^delta\\."

# Test with explicit config
differential --config ~/.config/differential/config.toml file1 file2

# Disable git config to test
differential --no-gitconfig file1 file2
```

### Theme Not Working

```bash
# List available themes
differential --list-themes

# Test specific theme
differential --theme dracula file1 file2
```

### Width Issues

```bash
# Check detected width
differential --width 0 file1 file2  # auto-detect

# Force specific width
differential --width 160 file1 file2
```

## See Also

- [Delta Configuration Docs](https://dandavison.github.io/delta/configuration.html)
- [Quick Start Guide](QUICKSTART.md)
- [Full README](README.md)
