package themes_test

import (
	"testing"

	"github.com/avgvstvs96/differential/internal/themes"
)

func TestInitialize(t *testing.T) {
	err := themes.Initialize()
	if err != nil {
		t.Fatalf("failed to initialize themes: %v", err)
	}

	// Check that we have themes loaded
	themeList := themes.ListThemes()
	if len(themeList) == 0 {
		t.Error("expected themes to be loaded, got none")
	}
}

func TestSetTheme(t *testing.T) {
	// Initialize first
	if err := themes.Initialize(); err != nil {
		t.Fatalf("failed to initialize themes: %v", err)
	}

	tests := []struct {
		name    string
		theme   string
		wantErr bool
	}{
		{
			name:    "valid theme",
			theme:   "dracula",
			wantErr: false,
		},
		{
			name:    "invalid theme",
			theme:   "nonexistent",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := themes.SetTheme(tt.theme)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetTheme(%q) error = %v, wantErr %v", tt.theme, err, tt.wantErr)
			}
		})
	}
}

func TestListThemes(t *testing.T) {
	// Initialize first
	if err := themes.Initialize(); err != nil {
		t.Fatalf("failed to initialize themes: %v", err)
	}

	themes := themes.ListThemes()
	
	// We should have at least the embedded themes
	expectedThemes := []string{"dracula", "monokai", "nord", "github", "gruvbox", "catppuccin", "tokyonight", "solarized"}
	
	if len(themes) < len(expectedThemes) {
		t.Errorf("expected at least %d themes, got %d", len(expectedThemes), len(themes))
	}

	// Check that specific themes exist
	themeMap := make(map[string]bool)
	for _, theme := range themes {
		themeMap[theme] = true
	}

	for _, expected := range expectedThemes {
		if !themeMap[expected] {
			t.Errorf("expected theme %q not found", expected)
		}
	}
}

func TestGetCurrentTheme(t *testing.T) {
	// Initialize and set a theme
	if err := themes.Initialize(); err != nil {
		t.Fatalf("failed to initialize themes: %v", err)
	}
	if err := themes.SetTheme("dracula"); err != nil {
		t.Fatalf("failed to set theme: %v", err)
	}

	theme := themes.GetCurrentTheme()
	if theme == nil {
		t.Fatal("expected current theme, got nil")
	}

	// Check that some colors are set
	if theme.Text == "" {
		t.Error("expected Text color to be set")
	}
	if theme.DiffAdded == "" {
		t.Error("expected DiffAdded color to be set")
	}
	if theme.DiffRemoved == "" {
		t.Error("expected DiffRemoved color to be set")
	}
}