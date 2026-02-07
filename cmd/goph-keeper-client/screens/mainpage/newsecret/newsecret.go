package newsecret_screen

import (
	"fmt"
	"strings"

	"github.com/Skifskii/goph-keeper/cmd/goph-keeper-client/screens"
	newcred_screen "github.com/Skifskii/goph-keeper/cmd/goph-keeper-client/screens/mainpage/newsecret/newcred"
	tea "github.com/charmbracelet/bubbletea"
)

var Name = "newsecret"

type Screen struct {
	options    []string
	focusIndex int
	Done       bool
	Quit       bool
}

func NewScreen() *Screen {
	return &Screen{
		options: []string{
			newcred_screen.Name,
		},
	}
}

func (m *Screen) Init() tea.Cmd {
	m.focusIndex = 0
	m.Done = false
	m.Quit = false
	return nil
}

func (m Screen) NextScreen() string {
	return m.options[m.focusIndex]
}

func (m *Screen) Update(msg tea.Msg) (*Screen, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "up":
			if m.focusIndex > 0 {
				m.focusIndex--
			}

		case "down":
			if m.focusIndex < len(m.options)-1 {
				m.focusIndex++
			}

		case "enter":
			m.Done = true
			return m, nil

		case "q":
			m.Done = true
			m.Quit = true
			return m, nil
		}
	}

	return m, nil
}

func (m Screen) View() string {
	var b strings.Builder

	// title
	b.WriteString(screens.LabelStyle.Render("    goph-keeper / main page / new secret"))
	b.WriteString("\n\n\n")

	// login
	b.WriteString("    Chose a type of secret")
	b.WriteString("\n")
	b.WriteString(m.buildRow("credential", 0))
	b.WriteString("\n")
	b.WriteString("\n\n\n")

	b.WriteString(screens.HelpStyle.Render("    Use '↑', '↓' and 'Enter' to navigate, 'q' to exit"))
	b.WriteString("\n")

	return b.String()
}

func (m Screen) buildRow(option string, optionIndex int) string {
	cursor := " "
	if optionIndex == m.focusIndex {
		cursor = ">"
	}

	return fmt.Sprintf("[%s] %s", cursor, option)
}
