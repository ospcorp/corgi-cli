package ui

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"fmt"
	"strings"
)

// MODEL DATA

type corgiTui struct {
	text      string
	textInput textinput.Model
	quitting  bool
}

func NewCorgiTui(text string) corgiTui {
	ti := textinput.New()
	ti.Placeholder = "What's on your mind?"
	ti.SetVirtualCursor(false)
	ti.Focus()
	ti.CharLimit = 156
	ti.SetWidth(20)

	return corgiTui{text: text, textInput: ti}
}

func (s corgiTui) Init() tea.Cmd { return textinput.Blink }

// VIEW

func (s corgiTui) View() tea.View {
	textLen := len(s.text)
	topAndBottomBar := strings.Repeat("*", textLen+4)
	greeting := fmt.Sprintf("%s\n* %s *\n%s", topAndBottomBar, s.text, topAndBottomBar)

	var c *tea.Cursor
	if !s.textInput.VirtualCursor() {
		c = s.textInput.Cursor()
		c.Y += lipgloss.Height(greeting)
	}

	inputBlock := lipgloss.JoinVertical(lipgloss.Top, s.textInput.View(), s.footerView())
	if s.quitting {
		inputBlock += "\n"
	}

	v := tea.NewView(lipgloss.JoinVertical(lipgloss.Top, greeting, inputBlock))
	v.Cursor = c
	return v
}

// UPDATE

func (s corgiTui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "enter", "esc":
			s.quitting = true
			return s, tea.Quit
		}
	}

	s.textInput, cmd = s.textInput.Update(msg)
	return s, cmd
}

func (s corgiTui) footerView() string { return "\n(esc to quit)" }
