package diff_test

import (
	"testing"

	"github.com/avgvstvs96/differential/internal/diff"
)

func TestStripANSI(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "plain text",
			input:    "Hello, World!",
			expected: "Hello, World!",
		},
		{
			name:     "simple color",
			input:    "\x1b[31mRed Text\x1b[0m",
			expected: "Red Text",
		},
		{
			name:     "multiple escapes",
			input:    "\x1b[1m\x1b[32mBold Green\x1b[0m Normal",
			expected: "Bold Green Normal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := diff.StripANSI(tt.input)
			if result != tt.expected {
				t.Errorf("StripANSI(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestVisibleLength(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "plain text",
			input:    "Hello",
			expected: 5,
		},
		{
			name:     "text with ANSI",
			input:    "\x1b[31mHello\x1b[0m",
			expected: 5,
		},
		{
			name:     "unicode text",
			input:    "Hello 世界",
			expected: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := diff.VisibleLength(tt.input)
			if result != tt.expected {
				t.Errorf("VisibleLength(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestHighlightIntralineChanges(t *testing.T) {
	hunk := &diff.Hunk{
		Lines: []diff.DiffLine{
			{
				Kind:    diff.LineRemoved,
				Content: "Hello World",
			},
			{
				Kind:    diff.LineAdded,
				Content: "Hello Differential",
			},
		},
	}

	diff.HighlightIntralineChanges(hunk)

	// Check that segments were created
	if len(hunk.Lines[0].Segments) == 0 {
		t.Error("expected segments for removed line")
	}
	if len(hunk.Lines[1].Segments) == 0 {
		t.Error("expected segments for added line")
	}

	// Check that "World" -> "Differential" was highlighted
	foundRemoved := false
	for _, seg := range hunk.Lines[0].Segments {
		if seg.Text == "World" && seg.Type == diff.LineRemoved {
			foundRemoved = true
			break
		}
	}
	if !foundRemoved {
		t.Error("expected 'World' to be highlighted in removed line")
	}

	foundAdded := false
	for _, seg := range hunk.Lines[1].Segments {
		if seg.Text == "Differential" && seg.Type == diff.LineAdded {
			foundAdded = true
			break
		}
	}
	if !foundAdded {
		t.Error("expected 'Differential' to be highlighted in added line")
	}
}