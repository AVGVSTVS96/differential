package themes

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Theme represents a color theme for the diff viewer
type Theme struct {
	Name  string                       `json:"name"`
	Defs  map[string]string           `json:"defs"`
	Theme map[string]map[string]string `json:"theme"`
}

// ThemeColors contains resolved color values for rendering
type ThemeColors struct {
	// Text colors
	Text       lipgloss.Color
	TextMuted  lipgloss.Color
	Error      lipgloss.Color

	// Diff colors
	DiffAdded           lipgloss.Color
	DiffRemoved         lipgloss.Color
	DiffContext         lipgloss.Color
	DiffAddedBg         lipgloss.Color
	DiffRemovedBg       lipgloss.Color
	DiffContextBg       lipgloss.Color
	DiffHighlightAdded  lipgloss.Color
	DiffHighlightRemoved lipgloss.Color
	DiffLineNumber      lipgloss.Color
	DiffAddedLineNumberBg   lipgloss.Color
	DiffRemovedLineNumberBg lipgloss.Color

	// Syntax colors
	SyntaxKeyword     lipgloss.Color
	SyntaxFunction    lipgloss.Color
	SyntaxType        lipgloss.Color
	SyntaxVariable    lipgloss.Color
	SyntaxString      lipgloss.Color
	SyntaxNumber      lipgloss.Color
	SyntaxComment     lipgloss.Color
	SyntaxOperator    lipgloss.Color
	SyntaxPunctuation lipgloss.Color

	// UI colors
	Background      lipgloss.Color
	BackgroundPanel lipgloss.Color
	Border          lipgloss.Color
	Selection       lipgloss.Color
}

var (
	currentTheme     *ThemeColors
	availableThemes  map[string]*Theme
	terminalIsDark   = true
)

// Initialize sets up the theme system
func Initialize() error {
	availableThemes = make(map[string]*Theme)
	
	// Detect terminal background
	detectTerminalBackground()
	
	// Load embedded themes
	if err := loadEmbeddedThemes(); err != nil {
		return fmt.Errorf("failed to load themes: %w", err)
	}
	
	// Set default theme
	if err := SetTheme("dracula"); err != nil {
		return err
	}
	
	return nil
}

// SetTheme activates a theme by name
func SetTheme(name string) error {
	theme, ok := availableThemes[name]
	if !ok {
		return fmt.Errorf("theme %s not found", name)
	}
	
	currentTheme = resolveTheme(theme)
	return nil
}

// GetCurrentTheme returns the current active theme
func GetCurrentTheme() *ThemeColors {
	if currentTheme == nil {
		// Return a default theme if not initialized
		return getDefaultTheme()
	}
	return currentTheme
}

// ListThemes returns all available theme names
func ListThemes() []string {
	themes := make([]string, 0, len(availableThemes))
	for name := range availableThemes {
		themes = append(themes, name)
	}
	return themes
}

// resolveTheme converts a Theme definition to resolved ThemeColors
func resolveTheme(theme *Theme) *ThemeColors {
	tc := &ThemeColors{}
	
	// Helper to resolve color references
	resolveColor := func(key string) lipgloss.Color {
		variant := "dark"
		if !terminalIsDark {
			variant = "light"
		}
		
		if colorMap, ok := theme.Theme[key]; ok {
			if color, ok := colorMap[variant]; ok {
				// Check if it's a reference to a defined color
				if definedColor, ok := theme.Defs[color]; ok {
					return lipgloss.Color(definedColor)
				}
				return lipgloss.Color(color)
			}
		}
		
		// Default color
		return lipgloss.Color("#ffffff")
	}
	
	// Resolve all colors
	tc.Text = resolveColor("text")
	tc.TextMuted = resolveColor("textMuted")
	tc.Error = resolveColor("error")
	
	tc.DiffAdded = resolveColor("diffAdded")
	tc.DiffRemoved = resolveColor("diffRemoved")
	tc.DiffContext = resolveColor("diffContext")
	tc.DiffAddedBg = resolveColor("diffAddedBg")
	tc.DiffRemovedBg = resolveColor("diffRemovedBg")
	tc.DiffContextBg = resolveColor("diffContextBg")
	tc.DiffHighlightAdded = resolveColor("diffHighlightAdded")
	tc.DiffHighlightRemoved = resolveColor("diffHighlightRemoved")
	tc.DiffLineNumber = resolveColor("diffLineNumber")
	tc.DiffAddedLineNumberBg = resolveColor("diffAddedLineNumberBg")
	tc.DiffRemovedLineNumberBg = resolveColor("diffRemovedLineNumberBg")
	
	tc.SyntaxKeyword = resolveColor("syntaxKeyword")
	tc.SyntaxFunction = resolveColor("syntaxFunction")
	tc.SyntaxType = resolveColor("syntaxType")
	tc.SyntaxVariable = resolveColor("syntaxVariable")
	tc.SyntaxString = resolveColor("syntaxString")
	tc.SyntaxNumber = resolveColor("syntaxNumber")
	tc.SyntaxComment = resolveColor("syntaxComment")
	tc.SyntaxOperator = resolveColor("syntaxOperator")
	tc.SyntaxPunctuation = resolveColor("syntaxPunctuation")
	
	tc.Background = resolveColor("background")
	tc.BackgroundPanel = resolveColor("backgroundPanel")
	tc.Border = resolveColor("border")
	tc.Selection = resolveColor("selection")
	
	return tc
}

// getDefaultTheme returns a basic default theme
func getDefaultTheme() *ThemeColors {
	return &ThemeColors{
		Text:                    lipgloss.Color("#ffffff"),
		TextMuted:               lipgloss.Color("#999999"),
		Error:                   lipgloss.Color("#ff0000"),
		DiffAdded:               lipgloss.Color("#00ff00"),
		DiffRemoved:             lipgloss.Color("#ff0000"),
		DiffContext:             lipgloss.Color("#ffffff"),
		DiffAddedBg:             lipgloss.Color("#001100"),
		DiffRemovedBg:           lipgloss.Color("#110000"),
		DiffContextBg:           lipgloss.Color("#000000"),
		DiffHighlightAdded:      lipgloss.Color("#00ff00"),
		DiffHighlightRemoved:    lipgloss.Color("#ff0000"),
		DiffLineNumber:          lipgloss.Color("#666666"),
		DiffAddedLineNumberBg:   lipgloss.Color("#002200"),
		DiffRemovedLineNumberBg: lipgloss.Color("#220000"),
		SyntaxKeyword:           lipgloss.Color("#ff79c6"),
		SyntaxFunction:          lipgloss.Color("#50fa7b"),
		SyntaxType:              lipgloss.Color("#8be9fd"),
		SyntaxVariable:          lipgloss.Color("#f8f8f2"),
		SyntaxString:            lipgloss.Color("#f1fa8c"),
		SyntaxNumber:            lipgloss.Color("#bd93f9"),
		SyntaxComment:           lipgloss.Color("#6272a4"),
		SyntaxOperator:          lipgloss.Color("#ff79c6"),
		SyntaxPunctuation:       lipgloss.Color("#f8f8f2"),
		Background:              lipgloss.Color("#282a36"),
		BackgroundPanel:         lipgloss.Color("#44475a"),
		Border:                  lipgloss.Color("#6272a4"),
		Selection:               lipgloss.Color("#44475a"),
	}
}

// detectTerminalBackground attempts to detect if the terminal has a dark background
func detectTerminalBackground() {
	// Check environment variables
	colorScheme := os.Getenv("COLORFGBG")
	if colorScheme != "" {
		parts := strings.Split(colorScheme, ";")
		if len(parts) >= 2 {
			// If background is greater than 7, it's likely light
			if parts[1] > "7" {
				terminalIsDark = false
				return
			}
		}
	}
	
	// Check terminal name
	term := os.Getenv("TERM")
	if strings.Contains(term, "light") {
		terminalIsDark = false
		return
	}
	
	// Default to dark
	terminalIsDark = true
}

// LoadThemeFromJSON loads a theme from a JSON file
func LoadThemeFromJSON(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read theme file: %w", err)
	}
	
	var theme Theme
	if err := json.Unmarshal(data, &theme); err != nil {
		return fmt.Errorf("failed to parse theme JSON: %w", err)
	}
	
	if theme.Name == "" {
		// Extract name from filename
		parts := strings.Split(path, "/")
		filename := parts[len(parts)-1]
		theme.Name = strings.TrimSuffix(filename, ".json")
	}
	
	availableThemes[theme.Name] = &theme
	return nil
}