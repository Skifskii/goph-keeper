package login_screen

import (
	"fmt"
	"strings"

	"github.com/Skifskii/goph-keeper/cmd/goph-keeper-client/api"
	"github.com/Skifskii/goph-keeper/cmd/goph-keeper-client/screens"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var Name = "login"

type Screen struct {
	api        *api.APIClient
	options    []string
	username   textinput.Model
	password   textinput.Model
	focusIndex int
	loginErr   string
	Done       bool
}

func NewScreen(apiClient *api.APIClient) *Screen {
	username := textinput.New()
	username.Placeholder = "__________"
	username.Prompt = ""
	username.Focus()

	password := textinput.New()
	password.Placeholder = "__________"
	password.Prompt = ""
	password.EchoMode = textinput.EchoPassword
	password.EchoCharacter = '•'

	return &Screen{
		api:      apiClient,
		username: username,
		password: password,
		options: []string{
			"username",
			"password",
			"submit",
			"signup",
		},
	}
}

func (l Screen) Init() tea.Cmd {
	return nil
}

func (l *Screen) Update(msg tea.Msg) (*Screen, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "up":
			if l.focusIndex > 0 {
				l.focusIndex--
			}
			l.updateFocus()

		case "down":
			if l.focusIndex < len(l.options)-1 {
				l.focusIndex++
			}
			l.updateFocus()

		case "enter":
			// submit
			if l.focusIndex == 2 {
				if !l.Done {
					err := l.api.Login(l.username.Value(), l.password.Value())
					if err != nil {
						l.loginErr = err.Error()
						return l, nil
					} else {
						l.Done = true
					}
				}
			}
		}
	}

	l.loginErr = ""
	cmd := l.updateInputs(msg)

	return l, cmd
}

func (l *Screen) updateFocus() {
	l.username.Blur()
	l.password.Blur()

	switch l.focusIndex {
	case 0:
		l.username.Focus()
	case 1:
		l.password.Focus()
	}
}

func (l *Screen) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 2)

	// Only update text inputs if they're focused
	switch l.focusIndex {
	case 0:
		l.username, cmds[0] = l.username.Update(msg)
	case 1:
		l.password, cmds[1] = l.password.Update(msg)
	}

	return tea.Batch(cmds...)
}

func (l Screen) View() string {
	var b strings.Builder

	// title
	b.WriteString(screens.LabelStyle.Render("    welcome to goph-keeper!"))
	b.WriteString("\n\n\n")

	// login
	b.WriteString("    LOGIN")
	b.WriteString("\n")
	b.WriteString(l.buildRow("username") + ": ")
	b.WriteString(l.username.View())
	b.WriteString("\n")
	b.WriteString(l.buildRow("password") + ": ")
	b.WriteString(l.password.View())
	b.WriteString("\n")
	if l.loginErr != "" {
		b.WriteString(screens.ErrorStyle.Render(l.loginErr))
		b.WriteString("\n")
	}
	b.WriteString(l.buildRow("submit"))
	b.WriteString("\n\n\n")

	// register
	b.WriteString(screens.ItalicStyle.Render("    Don't have an account?"))
	b.WriteString("\n")
	b.WriteString(l.buildRow("signup"))
	b.WriteString("\n\n\n")

	b.WriteString(screens.HelpStyle.Render("    Use '↑', '↓' and 'Enter' to navigate"))
	b.WriteString("\n")

	return b.String()
}

func (l Screen) buildRow(option string) string {
	cursor := " "
	if option == l.options[l.focusIndex] {
		cursor = ">"
	}

	return fmt.Sprintf("[%s] %s", cursor, l.optionToUIName(option))
}

func (l Screen) optionToUIName(op string) string {
	switch op {
	case "username":
		return "username"
	case "password":
		return "password"
	case "submit":
		return "[ submit ]"
	case "signup":
		return "[ sign up ]"
	default:
		return "ERROR: unknown"
	}
}
