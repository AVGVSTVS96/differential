package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/avgvstvs96/differential/pkg/diff"
	"github.com/avgvstvs96/differential/pkg/styles"
	"github.com/avgvstvs96/differential/pkg/theme"
	"github.com/spf13/cobra"
)

var (
	version    = "0.2.0"
	themeName  string
	sideBySide bool
	width      int
	listThemes bool
)

var rootCmd = &cobra.Command{
	Use:   "differential [file1] [file2]",
	Short: "A beautiful diff tool with syntax highlighting",
	Long: `Differential is a terminal tool for viewing diffs with syntax highlighting,
character-level changes, and beautiful themes.

It can be used as a drop-in replacement for git diff:
  git diff | differential
  differential file1.go file2.go`,
	RunE: runDiff,
}

func init() {
	rootCmd.Flags().StringVarP(&themeName, "theme", "t", "opencode", "Color theme to use")
	rootCmd.Flags().BoolVarP(&sideBySide, "side-by-side", "s", false, "Show diff in side-by-side view")
	rootCmd.Flags().IntVarP(&width, "width", "w", 0, "Terminal width (auto-detect if 0)")
	rootCmd.Flags().BoolVar(&listThemes, "list-themes", false, "List available themes")
}

func runDiff(cmd *cobra.Command, args []string) error {
	// Initialize themes
	if err := theme.LoadThemesFromJSON(); err != nil {
		return fmt.Errorf("failed to load themes: %w", err)
	}

	// List themes mode
	if listThemes {
		themes := theme.AvailableThemes()
		fmt.Println("Available themes:")
		for _, name := range themes {
			fmt.Printf("  - %s\n", name)
		}
		return nil
	}

	// Set theme
	if err := theme.SetTheme(themeName); err != nil {
		return fmt.Errorf("invalid theme %q: %w", themeName, err)
	}

	// Detect terminal info
	termInfo := theme.DetectTerminal()
	theme.Terminal = termInfo
	styles.Terminal.BackgroundIsDark = termInfo.BackgroundIsDark

	// Auto-detect width if not specified
	if width == 0 {
		width = getTerminalWidth()
	}

	// Get diff input
	var diffText string
	var err error

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// Reading from pipe
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read stdin: %w", err)
		}
		diffText = string(data)
	} else if len(args) == 2 {
		// Compare two files
		diffText, err = runGitDiff(args[0], args[1])
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("usage: differential [file1] [file2] or pipe diff output to differential")
	}

	// Render diff
	var output string
	if sideBySide {
		output, err = diff.FormatDiff("", diffText, diff.WithWidth(width))
	} else {
		output, err = diff.FormatUnifiedDiff("", diffText, diff.WithWidth(width))
	}
	if err != nil {
		return fmt.Errorf("failed to render diff: %w", err)
	}

	fmt.Print(output)
	return nil
}

func runGitDiff(file1, file2 string) (string, error) {
	cmd := exec.Command("diff", "-u", file1, file2)
	output, err := cmd.Output()
	if err != nil {
		// diff returns exit code 1 when files differ, which is normal
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return string(output), nil
		}
		return "", fmt.Errorf("diff command failed: %w", err)
	}
	return string(output), nil
}

func getTerminalWidth() int {
	cmd := exec.Command("tput", "cols")
	output, err := cmd.Output()
	if err != nil {
		return 120 // Default width
	}

	var w int
	fmt.Sscanf(string(output), "%d", &w)
	if w <= 0 {
		return 120
	}
	return w
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
