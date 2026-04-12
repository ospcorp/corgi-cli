package ui

import (
	"bytes"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"fmt"
)

// MODEL DATA

type Message struct {
	Author  string
	Content string
}

type corgiTui struct {
	text      string
	textInput textinput.Model
	messages  []Message
	quitting  bool
}

func NewCorgiTui(text string) corgiTui {
	ti := textinput.New()
	ti.Placeholder = "What's on your mind?"
	ti.SetVirtualCursor(false)
	ti.Focus()
	ti.CharLimit = 156
	ti.SetWidth(100)

	return corgiTui{text: text, textInput: ti}
}

func (s corgiTui) Init() tea.Cmd { return textinput.Blink }

// VIEW

func (s corgiTui) View() tea.View {
	var b bytes.Buffer
	for _, v := range s.messages {
		styledAuthor := authorStyle.Render(string(v.Author))
		b.WriteString(fmt.Sprintf("%s: %s\n", styledAuthor, v.Content))
	}

	var c *tea.Cursor
	if !s.textInput.VirtualCursor() {
		c = s.textInput.Cursor()
		c.Y += lipgloss.Height(greeting) + lipgloss.Height(b.String())
	}

	inputBlock := lipgloss.JoinVertical(lipgloss.Top, s.textInput.View(), s.footerView())
	if s.quitting {
		inputBlock += "\n"
	}

	v := tea.NewView(lipgloss.JoinVertical(lipgloss.Top, greeting, b.String(), inputBlock))
	v.Cursor = c
	return v
}

// UPDATE

func (s corgiTui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			s.quitting = true
			return s, tea.Quit
		case "enter":
			s.messages = append(s.messages, Message{
				Author:  "Agamemnon",
				Content: s.textInput.Value(),
			})
			s.textInput.Reset()
			return s, cmd
		}
	}

	s.textInput, cmd = s.textInput.Update(msg)
	return s, cmd
}

func (s corgiTui) footerView() string { return "\n(esc to quit)" }
