package ui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

// MODEL DATA

type corgiTui struct{ text string }

func NewCorgiTui(text string) corgiTui {
	return corgiTui{text: text}
}

func (s corgiTui) Init() tea.Cmd { return nil }

// VIEW

func (s corgiTui) View() tea.View {
	textLen := len(s.text)
	topAndBottomBar := strings.Repeat("*", textLen+4)
	return tea.NewView(fmt.Sprintf(
		"%s\n* %s *\n%s\n\nPress Ctrl+C to exit",
		topAndBottomBar, s.text, topAndBottomBar,
	))
}

// UPDATE

func (s corgiTui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return s, tea.Quit
		}
	}
	return s, nil
}
