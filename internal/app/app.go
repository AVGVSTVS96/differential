package app

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/avgvstvs96/differential/internal/config"
	"github.com/avgvstvs96/differential/internal/diff"
	"github.com/avgvstvs96/differential/internal/themes"
)

// Mode represents the current mode of the application
type Mode int

const (
	ModeBrowse Mode = iota
	ModeDiff
	ModeSearch
	ModeHelp
)

// Model represents the main application state
type Model struct {
	// Application state
	mode         Mode
	config       *config.Config
	windowWidth  int
	windowHeight int
	ready        bool
	err          error

	// Current diff
	diffResult   *diff.DiffResult
	diffText     string
	filename     string
	viewMode     diff.ViewMode

	// Navigation
	scrollOffset int
	selectedHunk int
	selectedLine int

	// UI state
	showLineNumbers bool
	contextLines    int
}

// RunPipeMode runs the application in pipe mode (non-interactive)
func RunPipeMode(input io.Reader, cfg *config.Config, args []string) error {
	// Initialize themes
	if err := themes.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize themes: %w", err)
	}

	// Set theme
	if err := themes.SetTheme(cfg.UI.Theme); err != nil {
		return fmt.Errorf("failed to set theme: %w", err)
	}

	var diffText string
	var err error

	// Get diff text from input or generate from files
	if input != nil {
		// Read from stdin
		data, err := io.ReadAll(input)
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		diffText = string(data)
	} else if len(args) == 2 {
		// Generate diff from two files
		diffText, err = runDiff(args[0], args[1])
		if err != nil {
			return fmt.Errorf("failed to diff files: %w", err)
		}
	} else if len(args) > 0 {
		// Pass args to git diff
		diffText, err = runGitDiff(args)
		if err != nil {
			return fmt.Errorf("failed to run git diff: %w", err)
		}
	} else {
		return fmt.Errorf("no diff input provided")
	}

	// Determine terminal width
	width := getTerminalWidth()

	// Create render options
	opts := diff.RenderOptions{
		Width:           width,
		ShowLineNumbers: cfg.UI.LineNumbers,
		ContextLines:    cfg.Git.DefaultContext,
		TabWidth:        cfg.UI.TabWidth,
	}

	// Format based on view mode
	var output string
	if cfg.UI.DefaultView == "side-by-side" {
		opts.ViewMode = diff.ViewSideBySide
		output, err = diff.FormatSideBySideDiff("", diffText, opts)
	} else {
		opts.ViewMode = diff.ViewUnified
		output, err = diff.FormatUnifiedDiff("", diffText, opts)
	}

	if err != nil {
		return fmt.Errorf("failed to format diff: %w", err)
	}

	// Determine if we should use a pager
	termHeight := getTerminalHeight()
	lineCount := strings.Count(output, "\n")
	
	// Show inline if diff is small enough to fit in terminal
	if lineCount < (termHeight - 5) {
		fmt.Print(output)
		return nil
	}
	
	// Use pager for larger diffs (unless disabled)
	if shouldUsePager() {
		return showWithPager(output)
	}

	fmt.Print(output)
	return nil
}

// RunTUIMode runs the application in TUI mode (interactive)
func RunTUIMode(args []string, cfg *config.Config) error {
	// Initialize themes
	if err := themes.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize themes: %w", err)
	}

	// Set theme
	if err := themes.SetTheme(cfg.UI.Theme); err != nil {
		return fmt.Errorf("failed to set theme: %w", err)
	}

	// Create initial model
	m := Model{
		mode:            ModeDiff,
		config:          cfg,
		showLineNumbers: cfg.UI.LineNumbers,
		contextLines:    cfg.Git.DefaultContext,
		viewMode:        diff.ViewUnified,
	}

	// Handle different input modes
	if len(args) == 0 {
		// No args - try to run git diff in current directory
		diffText, err := runGitDiff([]string{})
		if err != nil {
			return fmt.Errorf("failed to get git diff: %w", err)
		}
		m.diffText = diffText
	} else if len(args) == 2 {
		// Two files - compare them
		diffText, err := runDiff(args[0], args[1])
		if err != nil {
			return fmt.Errorf("failed to diff files: %w", err)
		}
		m.diffText = diffText
		m.filename = args[1]
	} else {
		// Pass args to git diff
		diffText, err := runGitDiff(args)
		if err != nil {
			return fmt.Errorf("failed to run git diff: %w", err)
		}
		m.diffText = diffText
	}

	// Parse diff
	result, err := diff.ParseUnifiedDiff(m.diffText)
	if err != nil {
		return fmt.Errorf("failed to parse diff: %w", err)
	}
	m.diffResult = result

	// Start TUI
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running program: %w", err)
	}

	return nil
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		m.ready = true
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case error:
		m.err = msg
		return m, nil
	}

	return m, nil
}

// View renders the UI
func (m Model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	if m.err != nil {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff0000")).
			Render(fmt.Sprintf("Error: %v", m.err))
	}

	if m.diffResult == nil || len(m.diffResult.Hunks) == 0 {
		return "No changes to display"
	}

	// Render diff
	opts := diff.RenderOptions{
		Width:           m.windowWidth,
		ViewMode:        m.viewMode,
		ShowLineNumbers: m.showLineNumbers,
		ContextLines:    m.contextLines,
		TabWidth:        m.config.UI.TabWidth,
	}

	var output string
	if m.viewMode == diff.ViewSideBySide {
		output = diff.RenderSideBySideDiff(m.diffResult, opts)
	} else {
		output = diff.RenderUnifiedDiff(m.diffResult, opts)
	}

	// Apply scrolling
	lines := strings.Split(output, "\n")
	visibleLines := m.windowHeight - 2 // Leave room for status bar

	if m.scrollOffset >= len(lines) {
		m.scrollOffset = len(lines) - 1
	}
	if m.scrollOffset < 0 {
		m.scrollOffset = 0
	}

	end := m.scrollOffset + visibleLines
	if end > len(lines) {
		end = len(lines)
	}

	visible := strings.Join(lines[m.scrollOffset:end], "\n")

	// Add status bar
	statusBar := m.renderStatusBar()

	return visible + "\n" + statusBar
}

// handleKeyPress handles keyboard input
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "j", "down":
		m.scrollOffset++
		return m, nil

	case "k", "up":
		if m.scrollOffset > 0 {
			m.scrollOffset--
		}
		return m, nil

	case "ctrl+f", "pgdown":
		m.scrollOffset += m.windowHeight - 2
		return m, nil

	case "ctrl+b", "pgup":
		m.scrollOffset -= m.windowHeight - 2
		if m.scrollOffset < 0 {
			m.scrollOffset = 0
		}
		return m, nil

	case "g", "home":
		m.scrollOffset = 0
		return m, nil

	case "G", "end":
		// Scroll to bottom
		totalLines := countLines(m.View())
		m.scrollOffset = totalLines - m.windowHeight + 2
		if m.scrollOffset < 0 {
			m.scrollOffset = 0
		}
		return m, nil

	case "tab":
		// Toggle view mode
		if m.viewMode == diff.ViewUnified {
			m.viewMode = diff.ViewSideBySide
		} else {
			m.viewMode = diff.ViewUnified
		}
		return m, nil

	case "n":
		// Toggle line numbers
		m.showLineNumbers = !m.showLineNumbers
		return m, nil

	case "?":
		// Show help
		m.mode = ModeHelp
		return m, nil
	}

	return m, nil
}

// renderStatusBar renders the bottom status bar
func (m Model) renderStatusBar() string {
	theme := themes.GetCurrentTheme()
	
	style := lipgloss.NewStyle().
		Background(theme.BackgroundPanel).
		Foreground(theme.Text).
		Width(m.windowWidth)

	// Build status text
	var parts []string

	// File info
	if m.diffResult.NewFile != "" {
		parts = append(parts, m.diffResult.NewFile)
	}

	// Stats
	additions, deletions := m.diffResult.CountChanges()
	parts = append(parts, fmt.Sprintf("+%d -%d", additions, deletions))

	// View mode
	viewMode := "Unified"
	if m.viewMode == diff.ViewSideBySide {
		viewMode = "Side-by-Side"
	}
	parts = append(parts, viewMode)

	// Line numbers
	if m.showLineNumbers {
		parts = append(parts, "Lines: ON")
	} else {
		parts = append(parts, "Lines: OFF")
	}

	// Controls hint
	parts = append(parts, "? for help")

	status := strings.Join(parts, " │ ")
	return style.Render(status)
}

// Helper functions

func getTerminalWidth() int {
	cmd := exec.Command("tput", "cols")
	output, err := cmd.Output()
	if err != nil {
		return 80 // Default
	}

	var width int
	fmt.Sscanf(string(output), "%d", &width)
	if width <= 0 {
		return 80
	}
	return width
}

func getTerminalHeight() int {
	cmd := exec.Command("tput", "lines")
	output, err := cmd.Output()
	if err != nil {
		return 24 // Default terminal height
	}

	var height int
	fmt.Sscanf(string(output), "%d", &height)
	if height <= 0 {
		return 24
	}
	return height
}

func shouldUsePager() bool {
	// Check if stdout is a terminal
	fi, _ := os.Stdout.Stat()
	return fi.Mode()&os.ModeCharDevice != 0
}

func showWithPager(content string) error {
	// Try common pagers
	pagers := []string{"less", "more"}

	for _, pager := range pagers {
		if _, err := exec.LookPath(pager); err == nil {
			cmd := exec.Command(pager, "-R") // -R for ANSI colors
			cmd.Stdin = strings.NewReader(content)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err == nil {
				return nil
			}
		}
	}

	// Fallback to direct output
	fmt.Print(content)
	return nil
}

func runGitDiff(args []string) (string, error) {
	cmdArgs := append([]string{"diff", "--no-color", "--no-ext-diff"}, args...)
	cmd := exec.Command("git", cmdArgs...)
	output, err := cmd.Output()
	if err != nil {
		// Check if it's just an empty diff
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return "", nil
		}
		return "", err
	}
	return string(output), nil
}

func runDiff(file1, file2 string) (string, error) {
	cmd := exec.Command("diff", "-u", file1, file2)
	output, err := cmd.Output()
	if err != nil {
		// diff returns exit code 1 when files differ, which is normal
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return string(output), nil
		}
		return "", err
	}
	return string(output), nil
}

func countLines(s string) int {
	return strings.Count(s, "\n") + 1
}