package tui

import tea "github.com/charmbracelet/bubbletea"

// KeyMap defines the keybindings for the TUI
type KeyMap struct {
	Up           []string
	Down         []string
	NextFile     []string
	PrevFile     []string
	NextHunk     []string
	PrevHunk     []string
	PageUp       []string
	PageDown     []string
	Top          []string
	Bottom       []string
	Search       []string
	ToggleTree   []string
	Quit         []string
	Help         []string
}

// DefaultKeyMap returns the default keybindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up:         []string{"k", "up"},
		Down:       []string{"j", "down"},
		NextFile:   []string{"n", "N"},
		PrevFile:   []string{"p", "P"},
		NextHunk:   []string{"]", "tab"},
		PrevHunk:   []string{"[", "shift+tab"},
		PageUp:     []string{"ctrl+b", "pgup"},
		PageDown:   []string{"ctrl+f", "pgdown"},
		Top:        []string{"g", "home"},
		Bottom:     []string{"G", "end"},
		Search:     []string{"/"},
		ToggleTree: []string{"t", "T"},
		Quit:       []string{"q", "Q", "esc", "ctrl+c"},
		Help:       []string{"?", "h"},
	}
}

// KeyMsg is a custom message type for key presses
type KeyMsg struct {
	Key string
}

// HandleKey processes a key press and returns the appropriate action
func (km KeyMap) HandleKey(msg tea.KeyMsg) string {
	key := msg.String()
	
	for _, k := range km.Up {
		if key == k {
			return "up"
		}
	}
	for _, k := range km.Down {
		if key == k {
			return "down"
		}
	}
	for _, k := range km.NextFile {
		if key == k {
			return "next-file"
		}
	}
	for _, k := range km.PrevFile {
		if key == k {
			return "prev-file"
		}
	}
	for _, k := range km.NextHunk {
		if key == k {
			return "next-hunk"
		}
	}
	for _, k := range km.PrevHunk {
		if key == k {
			return "prev-hunk"
		}
	}
	for _, k := range km.PageUp {
		if key == k {
			return "page-up"
		}
	}
	for _, k := range km.PageDown {
		if key == k {
			return "page-down"
		}
	}
	for _, k := range km.Top {
		if key == k {
			return "top"
		}
	}
	for _, k := range km.Bottom {
		if key == k {
			return "bottom"
		}
	}
	for _, k := range km.Search {
		if key == k {
			return "search"
		}
	}
	for _, k := range km.ToggleTree {
		if key == k {
			return "toggle-tree"
		}
	}
	for _, k := range km.Quit {
		if key == k {
			return "quit"
		}
	}
	for _, k := range km.Help {
		if key == k {
			return "help"
		}
	}
	
	return ""
}
