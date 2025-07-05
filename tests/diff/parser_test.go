package diff_test

import (
	"testing"

	"github.com/avgvstvs96/differential/internal/diff"
)

func TestParseUnifiedDiff_Empty(t *testing.T) {
	result, err := diff.ParseUnifiedDiff("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Hunks) != 0 {
		t.Errorf("expected 0 hunks, got %d", len(result.Hunks))
	}
}

func TestParseUnifiedDiff_Basic(t *testing.T) {
	input := `--- a/test.go	2025-01-01 00:00:00
+++ b/test.go	2025-01-01 00:00:01
@@ -1,3 +1,3 @@
 func main() {
-	fmt.Println("Hello")
+	fmt.Println("World")
 }`

	result, err := diff.ParseUnifiedDiff(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.OldFile != "test.go" {
		t.Errorf("expected OldFile 'test.go', got '%s'", result.OldFile)
	}
	if result.NewFile != "test.go" {
		t.Errorf("expected NewFile 'test.go', got '%s'", result.NewFile)
	}
	if len(result.Hunks) != 1 {
		t.Fatalf("expected 1 hunk, got %d", len(result.Hunks))
	}

	hunk := result.Hunks[0]
	if len(hunk.Lines) != 4 {
		t.Errorf("expected 4 lines, got %d", len(hunk.Lines))
	}
	
	// Verify line types
	expectedKinds := []diff.LineType{
		diff.LineContext,  // func main() {
		diff.LineRemoved,  // -	fmt.Println("Hello")
		diff.LineAdded,    // +	fmt.Println("World")
		diff.LineContext,  // }
	}
	
	for i, expected := range expectedKinds {
		if hunk.Lines[i].Kind != expected {
			t.Errorf("line %d: expected kind %v, got %v", i, expected, hunk.Lines[i].Kind)
		}
	}
}

func TestParseUnifiedDiff_GitStyle(t *testing.T) {
	input := `--- a/file.txt
+++ b/file.txt
@@ -1,2 +1,2 @@
 line1
-line2
+modified`

	result, err := diff.ParseUnifiedDiff(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.OldFile != "file.txt" {
		t.Errorf("expected OldFile 'file.txt', got '%s'", result.OldFile)
	}
	if result.NewFile != "file.txt" {
		t.Errorf("expected NewFile 'file.txt', got '%s'", result.NewFile)
	}
}

func TestParseUnifiedDiff_BinaryFile(t *testing.T) {
	input := `--- a/image.png
+++ b/image.png
Binary files a/image.png and b/image.png differ`

	result, err := diff.ParseUnifiedDiff(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.IsBinary {
		t.Errorf("expected IsBinary to be true")
	}
	if len(result.Hunks) != 0 {
		t.Errorf("expected 0 hunks for binary file, got %d", len(result.Hunks))
	}
}