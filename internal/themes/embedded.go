package themes

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

// Embedded theme JSON files
var (
	//go:embed themes/dracula.json
	draculaTheme string

	//go:embed themes/monokai.json
	monokaiTheme string

	//go:embed themes/nord.json
	nordTheme string

	//go:embed themes/github.json
	githubTheme string

	//go:embed themes/gruvbox.json
	gruvboxTheme string

	//go:embed themes/catppuccin.json
	catppuccinTheme string

	//go:embed themes/tokyonight.json
	tokyonightTheme string

	//go:embed themes/solarized.json
	solarizedTheme string
)

// loadEmbeddedThemes loads all embedded theme files
func loadEmbeddedThemes() error {
	themeData := map[string]string{
		"dracula":    draculaTheme,
		"monokai":    monokaiTheme,
		"nord":       nordTheme,
		"github":     githubTheme,
		"gruvbox":    gruvboxTheme,
		"catppuccin": catppuccinTheme,
		"tokyonight": tokyonightTheme,
		"solarized":  solarizedTheme,
	}

	for name, data := range themeData {
		var theme Theme
		if err := json.Unmarshal([]byte(data), &theme); err != nil {
			return fmt.Errorf("failed to parse %s theme: %w", name, err)
		}
		theme.Name = name
		availableThemes[name] = &theme
	}

	return nil
}