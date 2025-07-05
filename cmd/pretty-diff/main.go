package main

import (
	"fmt"
	"io"
	"os"

	"github.com/pretty-diff/pretty-diff/internal/app"
	"github.com/pretty-diff/pretty-diff/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version = "0.1.0"
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "pretty-diff [file1] [file2]",
	Short: "A beautiful TUI for file diffing with syntax highlighting",
	Long: `Pretty-Diff is a terminal UI for viewing diffs with syntax highlighting,
character-level changes, and interactive navigation.

It can be used as a drop-in replacement for git diff:
  git diff | pretty-diff
  pretty-diff file1.go file2.go
  pretty-diff HEAD~3 HEAD`,
	RunE: runDiff,
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/pretty-diff/config.toml)")
	rootCmd.Flags().StringP("theme", "t", "dracula", "Color theme to use")
	rootCmd.Flags().BoolP("side-by-side", "s", false, "Show diff in side-by-side view")
	rootCmd.Flags().BoolP("line-numbers", "n", true, "Show line numbers")
	rootCmd.Flags().IntP("context", "c", 3, "Number of context lines to show")
	rootCmd.Flags().BoolP("list-themes", "", false, "List available themes")
	rootCmd.Flags().BoolP("no-pager", "", false, "Disable pager for output")
	rootCmd.Flags().BoolP("pipe-mode", "p", false, "Force pipe mode (non-interactive)")

	viper.BindPFlags(rootCmd.Flags())
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err == nil {
			viper.AddConfigPath(home + "/.config/pretty-diff")
		}
		viper.SetConfigType("toml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("PRETTY_DIFF")

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func runDiff(cmd *cobra.Command, args []string) error {
	cfg := config.NewConfig()

	// Apply CLI flags
	if theme, _ := cmd.Flags().GetString("theme"); theme != "" {
		cfg.UI.Theme = theme
	}
	if sideBySide, _ := cmd.Flags().GetBool("side-by-side"); sideBySide {
		cfg.UI.DefaultView = "side-by-side"
	}
	if lineNumbers, _ := cmd.Flags().GetBool("line-numbers"); !lineNumbers {
		cfg.UI.LineNumbers = false
	}

	// List themes mode
	if listThemes, _ := cmd.Flags().GetBool("list-themes"); listThemes {
		themes := []string{
			"ayu", "catppuccin", "cobalt2", "dracula", "everforest",
			"github", "gruvbox", "kanagawa", "material", "matrix",
			"monokai", "nord", "one-dark", "opencode", "palenight",
			"rosepine", "solarized", "synthwave84", "tokyonight", "zenburn",
		}
		for _, theme := range themes {
			fmt.Println(theme)
		}
		return nil
	}

	// Determine mode
	isPipeMode := false
	var input io.Reader

	// Check if stdin has data
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		isPipeMode = true
		input = os.Stdin
	}

	// Force pipe mode flag
	if forceMode, _ := cmd.Flags().GetBool("pipe-mode"); forceMode {
		isPipeMode = true
		// If no stdin input but files provided, we'll generate diff in RunPipeMode
		if input == nil && len(args) > 0 {
			// Don't set input, let RunPipeMode handle file args
		}
	}

	if isPipeMode {
		// Pipe mode - render diff and exit
		return app.RunPipeMode(input, cfg, args)
	}

	// TUI mode
	return app.RunTUIMode(args, cfg)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}