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
	noPager    bool
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
	rootCmd.Flags().BoolVar(&noPager, "no-pager", false, "Disable pager (output directly to stdout)")
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

	// Output the result (with pager if appropriate)
	if shouldUsePager() {
		if err := showWithPager(output); err != nil {
			// Fallback to direct output if pager fails
			fmt.Print(output)
		}
	} else {
		fmt.Print(output)
	}
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
	// Try to get terminal size from the controlling terminal
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	cmd.Stdout = nil
	cmd.Stderr = nil
	
	output, err := cmd.Output()
	if err == nil {
		var rows, cols int
		if n, _ := fmt.Sscanf(string(output), "%d %d", &rows, &cols); n == 2 && cols > 0 {
			return cols
		}
	}
	
	// Fallback to tput
	cmd = exec.Command("tput", "cols")
	output, err = cmd.Output()
	if err == nil {
		var w int
		if _, scanErr := fmt.Sscanf(string(output), "%d", &w); scanErr == nil && w > 0 {
			return w
		}
	}
	
	// Last resort: check COLUMNS environment variable
	if cols := os.Getenv("COLUMNS"); cols != "" {
		var w int
		if _, scanErr := fmt.Sscanf(cols, "%d", &w); scanErr == nil && w > 0 {
			return w
		}
	}
	
	return 120 // Final fallback
}

// shouldUsePager determines if output should be piped through a pager
func shouldUsePager() bool {
	// Don't use pager if explicitly disabled
	if noPager {
		return false
	}

	// Don't use pager if output is not a terminal (e.g., piped to file)
	stat, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		return false
	}

	return true
}

// showWithPager pipes output through a pager like less
func showWithPager(content string) error {
	// Check for pager preference in environment variables
	// Order: DIFFERENTIAL_PAGER, BAT_PAGER, PAGER, then default to "less"
	pagerCmd := os.Getenv("DIFFERENTIAL_PAGER")
	if pagerCmd == "" {
		pagerCmd = os.Getenv("BAT_PAGER")
	}
	if pagerCmd == "" {
		pagerCmd = os.Getenv("PAGER")
	}
	if pagerCmd == "" {
		pagerCmd = "less"
	}

	// Try to find the pager executable
	pagerPath, err := exec.LookPath(pagerCmd)
	if err != nil {
		// If pager not found, try alternatives
		for _, alternative := range []string{"less", "more", "cat"} {
			if path, err := exec.LookPath(alternative); err == nil {
				pagerPath = path
				pagerCmd = alternative
				break
			}
		}
		if pagerPath == "" {
			return fmt.Errorf("no pager found")
		}
	}

	// Set up pager command
	// Use -R flag for less to handle ANSI colors
	var cmd *exec.Cmd
	if pagerCmd == "less" {
		cmd = exec.Command(pagerPath, "-R")
	} else {
		cmd = exec.Command(pagerPath)
	}

	// Pipe content to pager
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the pager
	if err := cmd.Start(); err != nil {
		return err
	}

	// Write content to pager
	if _, err := io.WriteString(stdin, content); err != nil {
		stdin.Close()
		return err
	}
	stdin.Close()

	// Wait for pager to finish
	return cmd.Wait()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
