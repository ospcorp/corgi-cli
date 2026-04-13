package ui

import (
	"bytes"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"fmt"
	"regexp"
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
	content   string
	quitting  bool
	ready     bool
	viewport  viewport.Model
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
	var v tea.View
	var b bytes.Buffer
	for _, v := range s.messages {
		styledAuthor := authorStyle.Render(string(v.Author))
		b.WriteString(fmt.Sprintf("%s: %s\n", styledAuthor, v.Content))
	}

	var c *tea.Cursor
	if !s.textInput.VirtualCursor() {
		c = s.textInput.Cursor()
		c.Y += lipgloss.Height(greeting) + lipgloss.Height(s.viewport.View()) + lipgloss.Height(b.String())
	}

	if !s.ready {
		v.SetContent("\n Initializing...")
	} else {
		v.SetContent(fmt.Sprintf("%s\n%s\n%s", greeting, s.viewport.View(), s.footerView()))
	}
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	v.Cursor = c

	return v
}

// UPDATE

func (s corgiTui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(greeting)
		footerHeight := lipgloss.Height(s.footerView())
		verticalMarginHeight := headerHeight + footerHeight
		if !s.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			s.viewport = viewport.New(viewport.WithWidth(msg.Width), viewport.WithHeight(msg.Height-verticalMarginHeight))
			s.viewport.YPosition = headerHeight
			s.viewport.HighlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Background(lipgloss.Color("34"))
			s.viewport.SelectedHighlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Background(lipgloss.Color("47"))
			s.viewport.SetContent(s.content)
			s.viewport.SetHighlights(regexp.MustCompile("artichoke").FindAllStringIndex(s.content, -1))
			s.viewport.HighlightNext()
			s.ready = true
		} else {
			s.viewport.SetWidth(msg.Width)
			s.viewport.SetHeight(msg.Height - verticalMarginHeight)
		}

		return s, cmd

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
			var b bytes.Buffer

			for _, v := range s.messages {
				styledAuthor := authorStyle.Render(string(v.Author))
				b.WriteString(fmt.Sprintf("%s: %s\n", styledAuthor, v.Content))
			}

			s.content = b.String()
			s.viewport.SetContent(s.content)
			s.viewport.GotoBottom()

			return s, cmd
		}
	}

	s.textInput, cmd = s.textInput.Update(msg)
	return s, cmd
}

func (s corgiTui) footerView() string {
	return lipgloss.JoinVertical(lipgloss.Top, s.textInput.View())
}
