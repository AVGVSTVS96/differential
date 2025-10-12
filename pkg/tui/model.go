package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/avgvstvs96/differential/pkg/config"
	"github.com/avgvstvs96/differential/pkg/diff"
)

// Model represents the TUI state
type Model struct {
	// Diff data
	results    []diff.DiffResult
	hunkIndex  []diff.HunkMetadata
	
	// View state
	currentFile int
	currentHunk int
	scrollOffset int
	
	// Display state
	width  int
	height int
	
	// Configuration
	cfg     *config.Config
	keymap  KeyMap
	
	// Search state
	searchMode   bool
	searchQuery  string
	searchResults []int // Indices of matching hunks
	
	// UI state
	showHelp     bool
	showFileTree bool
	err          error
	quitting     bool
	
	// Rendered content (cached)
	renderedLines []string
}

// Err returns any error that occurred
func (m Model) Err() error {
	return m.err
}

// NewModel creates a new TUI model
func NewModel(results []diff.DiffResult, cfg *config.Config) Model {
	hunkIndex := diff.BuildHunkIndex(results)
	
	return Model{
		results:      results,
		hunkIndex:    hunkIndex,
		currentFile:  0,
		currentHunk:  0,
		scrollOffset: 0,
		cfg:          cfg,
		keymap:       DefaultKeyMap(),
		searchMode:   false,
		showHelp:     false,
		showFileTree: false, // Disable file tree by default until fixed
		quitting:     false,
	}
}

// Init initializes the model (required by Bubble Tea)
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle search mode separately
		if m.searchMode {
			return m.updateSearch(msg)
		}
		
		// Handle normal navigation
		action := m.keymap.HandleKey(msg)
		switch action {
		case "quit":
			m.quitting = true
			return m, tea.Quit
			
		case "help":
			m.showHelp = !m.showHelp
			return m, nil
			
		case "up":
			m.moveUp()
			return m, nil
			
		case "down":
			m.moveDown()
			return m, nil
			
		case "next-file":
			m.nextFile()
			return m, nil
			
		case "prev-file":
			m.prevFile()
			return m, nil
			
		case "next-hunk":
			m.nextHunk()
			return m, nil
			
		case "prev-hunk":
			m.prevHunk()
			return m, nil
			
		case "page-up":
			m.pageUp()
			return m, nil
			
		case "page-down":
			m.pageDown()
			return m, nil
			
		case "top":
			m.goToTop()
			return m, nil
			
		case "bottom":
			m.goToBottom()
			return m, nil
			
		case "search":
			m.searchMode = true
			m.searchQuery = ""
			return m, nil
			
		case "toggle-tree":
			m.showFileTree = !m.showFileTree
			return m, nil
		}
		
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}
	
	return m, nil
}

// updateSearch handles search mode input
func (m Model) updateSearch(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "ctrl+c":
		m.searchMode = false
		m.searchQuery = ""
		return m, nil
		
	case "enter":
		m.performSearch()
		m.searchMode = false
		return m, nil
		
	case "backspace":
		if len(m.searchQuery) > 0 {
			m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
		}
		return m, nil
		
	default:
		// Add character to search query
		if len(msg.Runes) > 0 {
			m.searchQuery += string(msg.Runes)
		}
		return m, nil
	}
}

// Navigation methods

func (m *Model) moveUp() {
	if m.scrollOffset > 0 {
		m.scrollOffset--
	}
}

func (m *Model) moveDown() {
	if m.scrollOffset < len(m.renderedLines)-m.height {
		m.scrollOffset++
	}
}

func (m *Model) nextFile() {
	if m.currentFile < len(m.results)-1 {
		m.currentFile++
		m.currentHunk = m.findFirstHunkInFile(m.currentFile)
		m.scrollToCurrentHunk()
	}
}

func (m *Model) prevFile() {
	if m.currentFile > 0 {
		m.currentFile--
		m.currentHunk = m.findFirstHunkInFile(m.currentFile)
		m.scrollToCurrentHunk()
	}
}

func (m *Model) nextHunk() {
	if m.currentHunk < len(m.hunkIndex)-1 {
		m.currentHunk++
		// Update current file if we moved to a different file
		if m.currentHunk < len(m.hunkIndex) {
			m.currentFile = m.hunkIndex[m.currentHunk].FileIndex
		}
		m.scrollToCurrentHunk()
	}
}

func (m *Model) prevHunk() {
	if m.currentHunk > 0 {
		m.currentHunk--
		// Update current file if we moved to a different file
		if m.currentHunk >= 0 {
			m.currentFile = m.hunkIndex[m.currentHunk].FileIndex
		}
		m.scrollToCurrentHunk()
	}
}

func (m *Model) pageUp() {
	m.scrollOffset -= m.height
	if m.scrollOffset < 0 {
		m.scrollOffset = 0
	}
}

func (m *Model) pageDown() {
	m.scrollOffset += m.height
	maxOffset := len(m.renderedLines) - m.height
	if m.scrollOffset > maxOffset {
		m.scrollOffset = maxOffset
	}
	if m.scrollOffset < 0 {
		m.scrollOffset = 0
	}
}

func (m *Model) goToTop() {
	m.scrollOffset = 0
	m.currentFile = 0
	m.currentHunk = 0
}

func (m *Model) goToBottom() {
	m.scrollOffset = len(m.renderedLines) - m.height
	if m.scrollOffset < 0 {
		m.scrollOffset = 0
	}
	m.currentFile = len(m.results) - 1
	if len(m.hunkIndex) > 0 {
		m.currentHunk = len(m.hunkIndex) - 1
	}
}

// Helper methods

func (m *Model) findFirstHunkInFile(fileIndex int) int {
	for i, hunk := range m.hunkIndex {
		if hunk.FileIndex == fileIndex {
			return i
		}
	}
	return 0
}

func (m *Model) scrollToCurrentHunk() {
	// This would calculate the line offset for the current hunk
	// and adjust scrollOffset to make it visible
	// For now, just ensure we're within bounds
	if m.scrollOffset < 0 {
		m.scrollOffset = 0
	}
}

func (m *Model) performSearch() {
	if m.searchQuery == "" {
		return
	}
	
	m.searchResults = []int{}
	query := strings.ToLower(m.searchQuery)
	
	// Search through all hunks
	for i, hunkMeta := range m.hunkIndex {
		// Check if hunk context matches
		if strings.Contains(strings.ToLower(hunkMeta.Context), query) {
			m.searchResults = append(m.searchResults, i)
			continue
		}
		
		// Check if any line content matches
		fileResult := m.results[hunkMeta.FileIndex]
		if hunkMeta.HunkIndex < len(fileResult.Hunks) {
			hunk := fileResult.Hunks[hunkMeta.HunkIndex]
			for _, line := range hunk.Lines {
				if strings.Contains(strings.ToLower(line.Content), query) {
					m.searchResults = append(m.searchResults, i)
					break
				}
			}
		}
	}
	
	// Jump to first result if found
	if len(m.searchResults) > 0 {
		m.currentHunk = m.searchResults[0]
		m.currentFile = m.hunkIndex[m.currentHunk].FileIndex
		m.scrollToCurrentHunk()
	}
}

// GetStatusLine returns the status line showing current position
func (m Model) GetStatusLine() string {
	if m.searchMode {
		return fmt.Sprintf("Search: %s", m.searchQuery)
	}
	
	totalFiles := len(m.results)
	totalHunks := len(m.hunkIndex)
	
	fileName := ""
	if m.currentFile < len(m.results) {
		fileName = m.results[m.currentFile].NewFile
	}
	
	return fmt.Sprintf("File %d/%d | Hunk %d/%d | %s | Press ? for help",
		m.currentFile+1, totalFiles,
		m.currentHunk+1, totalHunks,
		fileName)
}
