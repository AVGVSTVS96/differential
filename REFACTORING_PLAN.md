## Refactoring Plan

This branch will implement the complete Phase 1 refactoring from ARCHITECTURE_REVIEW.md:

**Goal:** Split pkg/diff/diff.go (976 lines) into 8-9 focused files

### Proposed Structure:
```
pkg/diff/
├── types.go          # Core types (DONE in Quick Wins)
├── parser.go         # ParseUnifiedDiff
├── highlighter.go    # HighlightIntralineChanges, applyHighlighting  
├── syntax.go         # SyntaxHighlight, highlightLine, getChromaColor
├── render.go         # renderLinePrefix, renderLineContent (shared primitives)
├── unified.go        # Unified view rendering
├── sidebyside.go     # Side-by-side view rendering  
└── format.go         # Public API (FormatDiff, FormatUnifiedDiff)
```

### Benefits:
- Files <250 lines each (easier to understand)
- Clear responsibilities (easier to test)
- Better navigation (find code by purpose)
- Ready for Phase 3-5 features

### Strategy:
1. Create new files with refactored code
2. Keep old diff.go temporarily  
3. Update imports gradually
4. Remove old file once complete
5. Run tests at each step

See ARCHITECTURE_REVIEW.md for full details.

