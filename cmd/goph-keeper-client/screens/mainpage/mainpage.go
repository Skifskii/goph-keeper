package mainpage_screen

import (
	"fmt"
	"strings"

	"github.com/Skifskii/goph-keeper/cmd/goph-keeper-client/screens"
	mysecrets_screen "github.com/Skifskii/goph-keeper/cmd/goph-keeper-client/screens/mainpage/mysecrets"
	newsecret_screen "github.com/Skifskii/goph-keeper/cmd/goph-keeper-client/screens/mainpage/newsecret"
	tea "github.com/charmbracelet/bubbletea"
)

var Name = "mainpage"

type Screen struct {
	options    []string
	focusIndex int
	Done       bool
}

func NewScreen() *Screen {
	return &Screen{
		options: []string{
			mysecrets_screen.Name,
			newsecret_screen.Name,
		},
	}
}

func (m Screen) NextScreen() string {
	return m.options[m.focusIndex]
}

func (m *Screen) Init() tea.Cmd {
	m.focusIndex = 0
	m.Done = false
	return nil
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
		}
	}

	return m, nil
}

func (m Screen) View() string {
	var b strings.Builder

	// title
	b.WriteString(screens.LabelStyle.Render("    goph-keeper | main page"))
	b.WriteString("\n\n\n")

	// login
	b.WriteString("    Chose an action")
	b.WriteString("\n")
	b.WriteString(m.buildRow(mysecrets_screen.Name))
	b.WriteString("\n")
	b.WriteString(m.buildRow(newsecret_screen.Name))
	b.WriteString("\n")
	b.WriteString("\n\n\n")

	b.WriteString(screens.HelpStyle.Render("    Use '↑', '↓' and 'Enter' to navigate"))
	b.WriteString("\n")

	return b.String()
}

func (m Screen) buildRow(option string) string {
	cursor := " "
	if option == m.options[m.focusIndex] {
		cursor = ">"
	}

	return fmt.Sprintf("[%s] %s", cursor, option)
}
