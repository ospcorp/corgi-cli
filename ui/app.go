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

type message struct {
	author  string
	content string
}

type model struct {
	text      string
	textInput textinput.Model
	messages  []message
	content   string
	quitting  bool
	ready     bool
	viewport  viewport.Model
}

func NewCorgiTui(text string) model {
	ti := textinput.New()
	ti.Placeholder = "What's on your mind?"
	ti.SetVirtualCursor(false)
	ti.Focus()
	ti.CharLimit = 156
	ti.SetWidth(100)

	return model{text: text, textInput: ti}
}

func (m model) Init() tea.Cmd { return textinput.Blink }

// VIEW

func (m model) View() tea.View {
	var v tea.View
	var b bytes.Buffer
	for _, v := range m.messages {
		styledauthor := authorStyle.Render(string(v.author))
		b.WriteString(fmt.Sprintf("%s: %s\n", styledauthor, v.content))
	}

	var c *tea.Cursor
	if !m.textInput.VirtualCursor() {
		c = m.textInput.Cursor()
		c.Y += lipgloss.Height(greeting) + lipgloss.Height(m.viewport.View()) + lipgloss.Height(b.String())
	}

	if !m.ready {
		v.SetContent("\n Initializing...")
	} else {
		v.SetContent(fmt.Sprintf("%s\n%s\n%s", greeting, m.viewport.View(), m.footerView()))
	}
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	v.Cursor = c

	return v
}

// UPDATE

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(greeting)
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight
		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(viewport.WithWidth(msg.Width), viewport.WithHeight(msg.Height-verticalMarginHeight))
			m.viewport.YPosition = headerHeight
			m.viewport.HighlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Background(lipgloss.Color("34"))
			m.viewport.SelectedHighlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Background(lipgloss.Color("47"))
			m.viewport.SetContent(m.content)
			m.viewport.SetHighlights(regexp.MustCompile("artichoke").FindAllStringIndex(m.content, -1))
			m.viewport.HighlightNext()
			m.ready = true
		} else {
			m.viewport.SetWidth(msg.Width)
			m.viewport.SetHeight(msg.Height - verticalMarginHeight)
		}

		return m, cmd

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			m.messages = append(m.messages, message{
				author:  "Agamemnon",
				content: m.textInput.Value(),
			})
			m.textInput.Reset()
			var b bytes.Buffer

			for _, v := range m.messages {
				styledauthor := authorStyle.Render(string(v.author))
				b.WriteString(fmt.Sprintf("%s: %s\n", styledauthor, v.content))
			}

			m.content = b.String()
			m.viewport.SetContent(m.content)
			m.viewport.GotoBottom()

			return m, cmd
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) footerView() string {
	return lipgloss.JoinVertical(lipgloss.Top, m.textInput.View())
}
