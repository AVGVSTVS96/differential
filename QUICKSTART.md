# Quick Start Guide

## The Binary is Ready!

The `differential` binary is already built and located at:
```
/Users/bassimshahidy/Documents/GitHub/side-projects/pretty-diff/differential
```

## Basic Usage

### 1. Compare Two Files
```bash
cd /Users/bassimshahidy/Documents/GitHub/side-projects/pretty-diff
./differential file1.go file2.go
```

### 2. Pipe Git Diff
```bash
cd your-git-repo
git diff | /Users/bassimshahidy/Documents/GitHub/side-projects/pretty-diff/differential
```

### 3. Use Different Themes
```bash
./differential file1.go file2.go --theme dracula
./differential file1.go file2.go --theme nord
./differential file1.go file2.go --theme monokai
```

### 4. Side-by-Side View
```bash
./differential file1.go file2.go --side-by-side
```

### 5. List All Themes
```bash
./differential --list-themes
```

## Install Globally (Optional)

To use `differential` from anywhere:

```bash
# Copy to your PATH
sudo cp differential /usr/local/bin/

# Or add to PATH in your shell config (~/.zshrc or ~/.bashrc)
export PATH="$PATH:/Users/bassimshahidy/Documents/GitHub/side-projects/pretty-diff"
```

Then you can use it from anywhere:
```bash
differential file1.go file2.go
git diff | differential
```

## Quick Test

Run this to see it in action:
```bash
cd /Users/bassimshahidy/Documents/GitHub/side-projects/pretty-diff

# Create test files
echo 'package main
func main() {
    println("old")
}' > /tmp/test1.go

echo 'package main
func main() {
    println("new")
}' > /tmp/test2.go

# Compare them
./differential /tmp/test1.go /tmp/test2.go
```

## Git Integration

Set as your default git pager:
```bash
git config --global pager.diff "differential"
git config --global pager.show "differential"
```

Now `git diff` and `git show` will automatically use differential!

## Available Themes

opencode (default), aura, ayu, catppuccin, cobalt2, dracula, everforest, github, gruvbox, kanagawa, material, matrix, monokai, nord, one-dark, palenight, rosepine, solarized, synthwave84, tokyonight, vesper, zenburn

## What's Working

✅ All 22 themes with proper colors
✅ Character-level diff highlighting
✅ Syntax highlighting for 100+ languages  
✅ Side-by-side and unified views
✅ Works with any diff format
✅ Fast parallel rendering
