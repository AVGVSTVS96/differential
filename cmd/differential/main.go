package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/avgvstvs96/differential/pkg/config"
	"github.com/avgvstvs96/differential/pkg/diff"
	"github.com/avgvstvs96/differential/pkg/styles"
	"github.com/avgvstvs96/differential/pkg/theme"
	"github.com/avgvstvs96/differential/pkg/tui"
	"github.com/spf13/cobra"
)

var (
	version    = "0.3.0"
	
	// Config flags
	configFile     string
	noGitConfig    bool
	
	// Display flags
	themeName           string
	sideBySide          bool
	width               int
	listThemes          bool
	noPager             bool
	unified             int
	tabs                int
	colorMode           string
	lineNumbers         bool
	relativePaths       bool
	keepPlusMinusMakers bool
	ignoreAllSpace      bool
	ignoreSpaceChange   bool
	ignoreBlankLines    bool
	
	// Feature flags
	hyperlinks bool
	navigate   bool
	
	// Advanced styling flags
	graphicalStats   bool
	lineNumbersStyle string
	fileSeparator    string
	hunkSeparator    string
	
	// Global config (loaded from file + CLI flags)
	globalConfig *config.Config
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
	// Config flags
	rootCmd.Flags().StringVar(&configFile, "config", "", "Load config file at PATH instead of default")
	rootCmd.Flags().BoolVar(&noGitConfig, "no-gitconfig", false, "Do not read any settings from git config")
	
	// Core flags
	rootCmd.Flags().StringVarP(&themeName, "theme", "t", "", "Color theme to use")
	rootCmd.Flags().BoolVarP(&sideBySide, "side-by-side", "s", false, "Show diff in side-by-side view")
	rootCmd.Flags().BoolVar(&listThemes, "list-themes", false, "List available themes")
	rootCmd.Flags().BoolVarP(&lineNumbers, "line-numbers", "n", false, "Display line numbers next to the diff")
	
	// Display flags
	rootCmd.Flags().IntVarP(&width, "width", "w", 0, "Terminal width (auto-detect if 0)")
	rootCmd.Flags().IntVarP(&unified, "unified", "U", 0, "Number of context lines (default: 3)")
	rootCmd.Flags().IntVar(&tabs, "tabs", 0, "Number of spaces to replace tab characters with (default: 4)")
	rootCmd.Flags().StringVar(&colorMode, "color", "", "When to use colors: auto, always, never (default: auto)")
	rootCmd.Flags().BoolVar(&relativePaths, "relative-paths", false, "Output file paths relative to current directory")
	rootCmd.Flags().BoolVar(&keepPlusMinusMakers, "keep-plus-minus-markers", false, "Prefix added/removed lines with +/- characters")
	
	// Whitespace flags (no short flags to avoid conflicts)
	rootCmd.Flags().BoolVar(&ignoreAllSpace, "ignore-all-space", false, "Ignore all white space")
	rootCmd.Flags().BoolVar(&ignoreSpaceChange, "ignore-space-change", false, "Ignore changes in amount of white space")
	rootCmd.Flags().BoolVar(&ignoreBlankLines, "ignore-blank-lines", false, "Ignore changes whose lines are all blank")
	
	// Feature flags
	rootCmd.Flags().BoolVar(&hyperlinks, "hyperlinks", false, "Enable clickable file path hyperlinks (auto-detected by default)")
	rootCmd.Flags().BoolVar(&navigate, "navigate", false, "Enable navigation mode (not yet implemented)")
	
	// Advanced styling flags
	rootCmd.Flags().BoolVar(&graphicalStats, "graphical-stats", false, "Show graphical statistics bars for file changes")
	rootCmd.Flags().StringVar(&lineNumbersStyle, "line-numbers-style", "", "Line number style: standard, compact, minimal, left, right")
	rootCmd.Flags().StringVar(&fileSeparator, "file-separator", "", "Separator between files: none, blank, line, double-line, dashed, dotted")
	rootCmd.Flags().StringVar(&hunkSeparator, "hunk-separator", "", "Separator between hunks: none, blank, line, double-line, dashed, dotted")
	
	// Pager flags
	rootCmd.Flags().BoolVar(&noPager, "no-pager", false, "Disable pager (output directly to stdout)")
}

func runDiff(cmd *cobra.Command, args []string) error {
	// Load configuration
	var err error
	globalConfig, err = config.Load(configFile, noGitConfig)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	
	// Apply CLI flag overrides
	applyCLIFlags(cmd)
	
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
	if err := theme.SetTheme(globalConfig.Core.Theme); err != nil {
		return fmt.Errorf("invalid theme %q: %w", globalConfig.Core.Theme, err)
	}

	// Detect terminal info
	termInfo := theme.DetectTerminal()
	theme.Terminal = termInfo
	styles.Terminal.BackgroundIsDark = termInfo.BackgroundIsDark

	// Auto-detect width if not specified
	if globalConfig.Display.Width == 0 {
		globalConfig.Display.Width = getTerminalWidth()
	}

	// Get diff input
	var diffText string

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
		diffText, err = runGitDiff(args[0], args[1], globalConfig)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("usage: differential [file1] [file2] or pipe diff output to differential")
	}

	// Parse diff (may contain multiple files)
	results, err := diff.ParseMultiFileDiff(diffText)
	if err != nil {
		return fmt.Errorf("failed to parse diff: %w", err)
	}

	// Launch interactive TUI if navigate mode is enabled
	if globalConfig.Features.Navigate {
		return runInteractiveTUI(results, globalConfig)
	}

	// Render each file's diff
	var sb strings.Builder
	for i, result := range results {
		var output string
		
		// Convert result back to unified diff format for rendering
		// (This is a bit of a workaround; ideally we'd refactor FormatDiff to accept DiffResult)
		resultText := formatDiffResultAsText(result)
		
		if globalConfig.Core.SideBySide {
			output, err = diff.FormatDiff("", resultText, globalConfig, diff.WithWidth(globalConfig.Display.Width))
		} else {
			output, err = diff.FormatUnifiedDiff("", resultText, globalConfig, diff.WithWidth(globalConfig.Display.Width))
		}
		if err != nil {
			return fmt.Errorf("failed to render diff for file %s: %w", result.NewFile, err)
		}
		
		sb.WriteString(output)
		
		// Add file separator if configured (except after last file)
		if i < len(results)-1 {
			sepStyle := diff.ParseSeparatorStyle(globalConfig.Display.FileSeparator)
			if sepStyle != diff.SeparatorNone {
				sb.WriteString(diff.RenderFileSeparator(sepStyle, globalConfig.Display.Width, theme.CurrentTheme()))
			} else {
				// Default to a blank line if no separator configured
				sb.WriteString("\n")
			}
		}
	}

	output := sb.String()

	// Strip ANSI codes if color mode is never
	if globalConfig.Colors.ColorMode == "never" {
		output = stripANSI(output)
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

// runInteractiveTUI launches the interactive TUI mode
func runInteractiveTUI(results []diff.DiffResult, cfg *config.Config) error {
	// Create the TUI model
	model := tui.NewModel(results, cfg)
	
	// Create the Bubble Tea program
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)
	
	// Run the program
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("TUI error: %w", err)
	}
	
	// Check if there was an error in the model
	if m, ok := finalModel.(tui.Model); ok {
		if m.Err() != nil {
			return m.Err()
		}
	}
	
	return nil
}

func runGitDiff(file1, file2 string, cfg *config.Config) (string, error) {
	args := []string{"-u"}
	
	// Context lines
	if cfg.Display.Unified > 0 {
		args = append(args, fmt.Sprintf("-U%d", cfg.Display.Unified))
	}
	
	// Whitespace handling
	if cfg.Display.IgnoreAllSpace {
		args = append(args, "-w")
	}
	if cfg.Display.IgnoreSpaceChange {
		args = append(args, "-b")
	}
	if cfg.Display.IgnoreBlankLines {
		args = append(args, "-B")
	}
	
	args = append(args, file1, file2)
	cmd := exec.Command("diff", args...)
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
	if globalConfig.Core.NoPager {
		return false
	}
	
	// Force no pager if color mode is never (usually means piping to file)
	if globalConfig.Colors.ColorMode == "never" {
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
	// Check for pager preference (config, then env vars)
	pagerCmd := globalConfig.Core.Pager
	if pagerCmd == "" {
		pagerCmd = os.Getenv("DIFFERENTIAL_PAGER")
	}
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

// formatDiffResultAsText converts a DiffResult back to unified diff text format
func formatDiffResultAsText(result diff.DiffResult) string {
	var sb strings.Builder
	
	// Write file headers
	sb.WriteString(fmt.Sprintf("--- a/%s\n", result.OldFile))
	sb.WriteString(fmt.Sprintf("+++ b/%s\n", result.NewFile))
	
	// Write each hunk
	for _, hunk := range result.Hunks {
		sb.WriteString(hunk.Header)
		sb.WriteString("\n")
		
		for _, line := range hunk.Lines {
			switch line.Kind {
			case diff.LineAdded:
				sb.WriteString("+")
			case diff.LineRemoved:
				sb.WriteString("-")
			case diff.LineContext:
				sb.WriteString(" ")
			}
			sb.WriteString(line.Content)
			sb.WriteString("\n")
		}
	}
	
	return sb.String()
}

// applyCLIFlags applies command-line flag overrides to the loaded config
func applyCLIFlags(cmd *cobra.Command) {
	// CLI flags override config file settings
	if themeName != "" {
		globalConfig.Core.Theme = themeName
	}
	if sideBySide {
		globalConfig.Core.SideBySide = sideBySide
	}
	if width != 0 {
		globalConfig.Display.Width = width
	}
	if noPager {
		globalConfig.Core.NoPager = noPager
	}
	if unified != 0 {
		globalConfig.Display.Unified = unified
	}
	if tabs != 0 {
		globalConfig.Display.Tabs = tabs
	}
	if colorMode != "" {
		globalConfig.Colors.ColorMode = colorMode
	}
	if lineNumbers {
		globalConfig.Core.LineNumbers = lineNumbers
	}
	if relativePaths {
		globalConfig.Display.RelativePaths = relativePaths
	}
	if keepPlusMinusMakers {
		globalConfig.Display.KeepPlusMinusMarkers = keepPlusMinusMakers
	}
	if ignoreAllSpace {
		globalConfig.Display.IgnoreAllSpace = ignoreAllSpace
	}
	if ignoreSpaceChange {
		globalConfig.Display.IgnoreSpaceChange = ignoreSpaceChange
	}
	if ignoreBlankLines {
		globalConfig.Display.IgnoreBlankLines = ignoreBlankLines
	}
	// Feature flags - only override if explicitly set
	if cmd.Flags().Changed("hyperlinks") {
		globalConfig.Features.Hyperlinks = hyperlinks
	}
	if navigate {
		globalConfig.Features.Navigate = navigate
	}
	
	// Advanced styling flags
	if graphicalStats {
		globalConfig.Display.GraphicalStats = graphicalStats
	}
	if lineNumbersStyle != "" {
		globalConfig.Display.LineNumbersStyle = lineNumbersStyle
	}
	if fileSeparator != "" {
		globalConfig.Display.FileSeparator = fileSeparator
	}
	if hunkSeparator != "" {
		globalConfig.Display.HunkSeparator = hunkSeparator
	}
}

// stripANSI removes ANSI escape sequences from a string
func stripANSI(s string) string {
	// Match ANSI escape sequences
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	return ansiRegex.ReplaceAllString(s, "")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
