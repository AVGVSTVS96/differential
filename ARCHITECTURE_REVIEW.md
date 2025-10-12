# Architecture Review & Improvement Plan

## Current State Analysis

### File Structure
```
.
├── cmd/differential/
│   └── main.go (389 lines) ⚠️ Growing large
├── pkg/
│   ├── ansi/ (EMPTY) ❌ Remove
│   ├── config/
│   │   ├── defaults.go (97 lines) ✅
│   │   ├── loader.go (336 lines) ⚠️ Could split
│   │   └── types.go (52 lines) ✅
│   ├── diff/
│   │   ├── diff.go (976 lines) ❌ TOO LARGE - needs splitting
│   │   ├── parse.go (58 lines) ✅
│   │   └── textutil.go (266 lines) ⚠️ Misplaced
│   ├── styles/
│   │   └── styles.go (46 lines) ✅
│   ├── theme/
│   │   ├── loader.go (408 lines) ⚠️
│   │   ├── manager.go (229 lines) ⚠️
│   │   ├── system.go (330 lines) ⚠️
│   │   ├── theme.go (215 lines) ⚠️
│   │   └── themes/ (JSON) ✅
│   └── util/
│       └── concurrency.go (40 lines) ✅
```

---

## 🔴 Critical Issues

### 1. **diff.go is a God Object (976 lines)**

**Problem:** Single file handles:
- Type definitions (LineType, Segment, DiffLine, Hunk, DiffResult, linePair, UnifiedConfig)
- Parsing (ParseUnifiedDiff)
- Intra-line diff (HighlightIntralineChanges)
- Syntax highlighting (SyntaxHighlight, highlightLine, getColor, getChromaColor)
- Style creation (createStyles)
- ANSI handling (applyHighlighting)
- Rendering primitives (renderLinePrefix, renderLineContent)
- Unified rendering (renderUnifiedLine, RenderUnifiedHunk)
- Side-by-side rendering (renderDiffColumnLine, renderLeftColumn, renderRightColumn, pairLines, RenderSideBySideHunk)
- Public API (FormatDiff, FormatUnifiedDiff)

**Impact:**
- Hard to understand and navigate
- Difficult to test individual components
- High coupling between concerns
- Violates Single Responsibility Principle

**Recommended Split:**
```
pkg/diff/
├── types.go          # Core types (LineType, Segment, DiffLine, Hunk, etc.)
├── parser.go         # ParseUnifiedDiff (move from current code)
├── highlighter.go    # HighlightIntralineChanges, applyHighlighting
├── syntax.go         # SyntaxHighlight, highlightLine, getChromaColor
├── styles.go         # createStyles (or move to pkg/styles)
├── render.go         # renderLinePrefix, renderLineContent (shared primitives)
├── unified.go        # Unified view rendering (renderUnifiedLine, RenderUnifiedHunk)
├── sidebyside.go     # Side-by-side rendering (pairLines, renderDiffColumnLine, etc.)
└── format.go         # Public API (FormatDiff, FormatUnifiedDiff)
```

### 2. **pkg/ansi is Empty**

**Problem:** Empty directory serves no purpose

**Solution:** Remove it
```bash
rm -rf pkg/ansi
```

### 3. **textutil.go Misplaced in diff Package**

**Problem:** Text utilities are generic, not diff-specific
- ExpandTabs, TruncateWithEllipsis, StripTrailingWhitespace
- IsBlankLine, NormalizeWhitespace, PadToWidth
- These are general string manipulation utilities

**Solution:** Move to dedicated utility package
```
pkg/text/
├── formatting.go    # ExpandTabs, TruncateWithEllipsis, PadToWidth
├── whitespace.go    # StripTrailingWhitespace, NormalizeWhitespace, IsBlankLine
└── width.go         # CalculateVisibleLength, SplitByWidth
```

Or keep simpler:
```
pkg/textutil/
└── textutil.go      # All text utilities (current content)
```

### 4. **main.go Growing Too Large (389 lines)**

**Problem:** Mix of CLI setup, business logic, and utilities

**Current functions:**
- init() - Flag setup
- runDiff() - Core business logic
- runGitDiff() - Git integration
- getTerminalWidth() - Terminal detection
- shouldUsePager() / showWithPager() - Pager handling
- applyCLIFlags() - Config merging
- stripANSI() - Text utility
- main() - Entry point

**Recommended Split:**
```
cmd/differential/
├── main.go           # Entry point, flag setup (minimal)
├── run.go            # runDiff (core logic)
├── git.go            # runGitDiff (git integration)
├── pager.go          # shouldUsePager, showWithPager
├── terminal.go       # getTerminalWidth, terminal utilities
└── flags.go          # Flag definitions, applyCLIFlags
```

Or create a sub-package:
```
cmd/differential/
├── main.go
└── app/
    ├── run.go
    ├── git.go
    ├── pager.go
    ├── terminal.go
    └── config.go
```

### 5. **config/loader.go is Large (336 lines)**

**Problem:** Handles multiple concerns:
- Git config loading
- TOML file loading  
- Merging logic
- Value parsing (parseInt, parseBool)

**Recommended Split:**
```
pkg/config/
├── types.go          # Config struct definitions
├── defaults.go       # Default values, mappings
├── loader.go         # Main Load() function, coordination
├── gitconfig.go      # LoadFromGitConfig, applyGitConfigValue
├── merge.go          # mergeConfig logic
└── parse.go          # parseInt, parseBool, value parsing
```

### 6. **Theme Package Organization Unclear**

**Current:**
```
pkg/theme/
├── theme.go (215 lines)      # Theme interface and BaseTheme
├── system.go (330 lines)     # System colors and generation
├── loader.go (408 lines)     # JSON loading and chroma generation
├── manager.go (229 lines)    # Theme management and detection
└── themes/ (JSON files)
```

**Issue:** Unclear separation between:
- Theme definition vs theme loading
- Theme management vs theme application
- System/generated themes vs file-based themes

**Recommended Organization:**
```
pkg/theme/
├── theme.go          # Theme interface, core types
├── builtin.go        # Built-in/system theme generation (from system.go)
├── loader.go         # JSON file loading
├── registry.go       # Theme registration and lookup (from manager.go)
├── chroma.go         # Chroma style generation
├── terminal.go       # Terminal detection (from manager.go)
└── themes/           # JSON theme files
```

---

## 🟡 Medium Priority Issues

### 7. **Inconsistent Naming Conventions**

**Issues:**
- `pkg/diff/parse.go` has `ParseStats` but main parsing is in `diff.go`
- `pkg/styles/styles.go` - redundant naming
- `pkg/config/loader.go` vs `pkg/theme/loader.go` - same name, different packages

**Recommendations:**
- Rename `pkg/styles/styles.go` → `pkg/styles/style.go` or `pkg/styles/terminal.go`
- Consolidate parsing in `pkg/diff/parser.go`
- Use more specific names: `config/gitloader.go`, `theme/jsonloader.go`

### 8. **Missing Clear Package Purpose Documentation**

**Problem:** Package-level godocs are minimal or missing

**Solution:** Add comprehensive package comments:
```go
// Package diff provides diff parsing, rendering, and formatting.
//
// Architecture:
// - types.go: Core data structures
// - parser.go: Unified diff parsing
// - highlighter.go: Intra-line diff highlighting
// - syntax.go: Syntax highlighting integration
// - unified.go: Unified view rendering
// - sidebyside.go: Side-by-side view rendering
// - format.go: High-level API
package diff
```

### 9. **Parse vs Render vs Format Terminology**

**Issue:** Unclear distinction in naming
- `ParseUnifiedDiff` - clear
- `RenderUnifiedHunk` - clear
- `FormatUnifiedDiff` - what's the difference from Render?

**Recommendation:** Establish clear terminology:
- **Parse**: String → Data structures
- **Render**: Data structures → Styled strings (single hunk)
- **Format**: Data structures → Complete output (multiple hunks)

---

## 🟢 What's Working Well

1. **Config Package Structure** - Clear separation of types, defaults, loader
2. **Small, Focused Files** - util/concurrency.go, styles/styles.go
3. **Clear Public APIs** - FormatDiff, FormatUnifiedDiff
4. **Good Type Definitions** - DiffLine, Hunk, DiffResult are well-designed
5. **Theme JSON Organization** - Themes as data, not code

---

## 📋 Recommended Refactoring Plan

### Phase 1: Critical (Do First)
1. **Remove pkg/ansi** (empty directory)
2. **Split diff.go** into 8-9 files by responsibility
3. **Move textutil.go** to pkg/textutil or pkg/text

### Phase 2: Important (Do Soon)
4. **Split main.go** into focused files/subpackage
5. **Reorganize config package** - separate git config loading
6. **Clarify theme package** structure

### Phase 3: Polish (Nice to Have)
7. **Add package-level documentation** to all packages
8. **Standardize naming conventions**
9. **Create architecture diagram** for documentation

---

## 🎯 Proposed Final Structure

```
differential/
├── cmd/differential/
│   ├── main.go              # Entry point (~50 lines)
│   ├── flags.go             # Flag definitions (~100 lines)
│   ├── run.go               # Core runDiff logic (~100 lines)
│   ├── git.go               # Git integration (~50 lines)
│   ├── pager.go             # Pager handling (~80 lines)
│   └── terminal.go          # Terminal utilities (~60 lines)
│
├── pkg/
│   ├── config/
│   │   ├── types.go         # Config structs
│   │   ├── defaults.go      # Default values
│   │   ├── loader.go        # Main Load() orchestration
│   │   ├── gitconfig.go     # Git config loading (~150 lines)
│   │   ├── toml.go          # TOML file loading (~50 lines)
│   │   ├── merge.go         # Config merging (~100 lines)
│   │   └── parse.go         # parseInt, parseBool (~30 lines)
│   │
│   ├── diff/
│   │   ├── types.go         # Core types (~100 lines)
│   │   ├── parser.go        # ParseUnifiedDiff (~150 lines)
│   │   ├── stats.go         # ParseStats (from parse.go)
│   │   ├── highlighter.go   # Intra-line highlighting (~150 lines)
│   │   ├── syntax.go        # Syntax highlighting (~200 lines)
│   │   ├── render.go        # Shared rendering primitives (~100 lines)
│   │   ├── unified.go       # Unified view (~150 lines)
│   │   ├── sidebyside.go    # Side-by-side view (~200 lines)
│   │   └── format.go        # Public API (~50 lines)
│   │
│   ├── textutil/
│   │   └── textutil.go      # Text utilities (current content)
│   │
│   ├── styles/
│   │   └── style.go         # Terminal styles (renamed)
│   │
│   ├── theme/
│   │   ├── theme.go         # Theme interface (~100 lines)
│   │   ├── builtin.go       # System theme generation (~200 lines)
│   │   ├── loader.go        # JSON loading (~250 lines)
│   │   ├── registry.go      # Theme registry (~150 lines)
│   │   ├── chroma.go        # Chroma integration (~100 lines)
│   │   ├── terminal.go      # Terminal detection (~50 lines)
│   │   └── themes/          # JSON files
│   │
│   └── util/
│       └── concurrency.go   # Parallel utilities
│
└── (docs, configs, etc.)
```

---

## 🎨 Benefits of Proposed Structure

### Improved Maintainability
- **Smaller files** (<250 lines each) easier to understand
- **Clear responsibilities** - each file has one job
- **Easier testing** - can test components in isolation
- **Better navigation** - find code by responsibility

### Enhanced Readability
- **Self-documenting** file names tell you what's inside
- **Logical grouping** - related code stays together
- **Clearer imports** - see exactly what you're using
- **Reduced cognitive load** - focus on one concern at a time

### Better Development Workflow
- **Parallel development** - multiple devs, different files
- **Reduced merge conflicts** - changes in different files
- **Easier code review** - review one concern at a time
- **Faster compilation** (potentially) - smaller compilation units

### Future-Proof
- **Easy to extend** - add new renderers/parsers/themes
- **Plugin-ready** - could externalize themes/syntax
- **Testable** - each component can be tested independently
- **Documentable** - each file can have focused docs

---

## 🚀 Migration Strategy

### Step 1: Create New Structure (No Breaking Changes)
1. Create new files with refactored code
2. Keep old code in place temporarily
3. Update imports gradually
4. Run tests at each step

### Step 2: Internal Refactoring
1. Move small, independent pieces first (textutil)
2. Split diff.go incrementally
3. Update tests as you go

### Step 3: Cleanup
1. Remove old files once all imports updated
2. Update documentation
3. Remove empty pkg/ansi

### Step 4: Polish
1. Add package-level docs
2. Standardize naming
3. Update CLAUDE.md with new structure

---

## 💡 Quick Wins (Do These Now)

1. **Remove pkg/ansi** - 5 minutes
   ```bash
   rm -rf pkg/ansi
   ```

2. **Move textutil** - 10 minutes
   ```bash
   mkdir pkg/textutil
   git mv pkg/diff/textutil.go pkg/textutil/
   # Update imports
   ```

3. **Split diff.go types** - 15 minutes
   ```bash
   # Extract type definitions to types.go
   # Keep rest in diff.go for now
   ```

4. **Add package docs** - 20 minutes
   ```go
   // Package diff provides...
   // Package config manages...
   // etc.
   ```

---

## 📊 Metrics Comparison

| Metric | Current | Proposed | Improvement |
|--------|---------|----------|-------------|
| Largest file | 976 lines | ~250 lines | 73% reduction |
| Files >500 lines | 2 files | 0 files | 100% reduction |
| Avg file size | 192 lines | ~100 lines | 48% reduction |
| Package clarity | Medium | High | Subjective ++ |
| Empty packages | 1 (ansi) | 0 | Clean |

---

## 🎓 Key Principles Applied

1. **Single Responsibility Principle** - Each file has one job
2. **Separation of Concerns** - Parse ≠ Render ≠ Format
3. **High Cohesion** - Related code together
4. **Low Coupling** - Clean interfaces between components
5. **Self-Documenting Code** - Structure tells the story
6. **Ease of Navigation** - Find code by responsibility
7. **Testability** - Small, focused units

---

## 🤔 Considerations

### Why Not Go Further?
- Could split into even more packages (pkg/parse, pkg/render, pkg/format)
- Could create internal/ for non-public APIs
- Could use Go workspaces for true modularity

**Answer:** Current proposal balances:
- ✅ Improved structure without over-engineering
- ✅ Maintains Go idioms (flat package structure)
- ✅ Easy to understand for newcomers
- ✅ Room to grow in future phases

### Why Not Stay As-Is?
- Current structure works but won't scale
- 976-line file is past the tipping point
- Adding Phase 3-5 features will make it worse
- Better to refactor now than later

---

## 📝 Conclusion

The current architecture is **functional but not optimal**. Key issues:

1. ❌ **diff.go is too large** (976 lines) - urgent
2. ❌ **Empty ansi package** - trivial fix
3. ⚠️ **Misplaced utilities** - easy fix
4. ⚠️ **Growing main.go** - will get worse

**Recommendation:** 
- Do **Phase 1 refactoring** (split diff.go) before adding Phase 3-5 features
- At minimum, do the **Quick Wins** immediately
- The proposed structure will support the project through Phase 5 and beyond

The refactoring will **pay dividends** as:
- Phase 3: Adding git metadata parsing
- Phase 4: Adding navigation/hyperlinks
- Phase 5: Adding advanced styling

**Better to refactor now** while the codebase is manageable than wait until it's unwieldy.
