package screens

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type MainPageScreen struct {
	options    []string
	focusIndex int
	Done       bool
}

func NewMainPageScreen() *MainPageScreen {
	return &MainPageScreen{
		options: []string{
			"mysecrets",
			"newsecret",
		},
	}
}

func (m MainPageScreen) NextScreen() string {
	return m.options[m.focusIndex]
}

func (m *MainPageScreen) Init() tea.Cmd {
	m.Done = false
	return nil
}

func (m *MainPageScreen) Update(msg tea.Msg) (*MainPageScreen, tea.Cmd) {
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

func (m MainPageScreen) View() string {
	var b strings.Builder

	// title
	b.WriteString(labelStyle.Render("    goph-keeper | main page"))
	b.WriteString("\n\n\n")

	// login
	b.WriteString("    Chose an action")
	b.WriteString("\n")
	b.WriteString(m.buildRow("mysecrets"))
	b.WriteString("\n")
	b.WriteString(m.buildRow("newsecret"))
	b.WriteString("\n")
	b.WriteString("\n\n\n")

	b.WriteString(helpStyle.Render("    Use ↑, ↓ and Enter to navigate, "))
	b.WriteString("\n")

	return b.String()
}

func (m MainPageScreen) buildRow(option string) string {
	cursor := " "
	if option == m.options[m.focusIndex] {
		cursor = ">"
	}

	return fmt.Sprintf("[%s] %s", cursor, m.optionToUIName(option))
}

func (m MainPageScreen) optionToUIName(op string) string {
	switch op {
	case "mysecrets":
		return "my secrets"
	case "newsecret":
		return "new secret"
	default:
		return "ERROR: unknown"
	}
}
