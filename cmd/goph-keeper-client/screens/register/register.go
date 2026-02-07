package register_screen

import (
	"fmt"
	"strings"

	"github.com/Skifskii/goph-keeper/cmd/goph-keeper-client/api"
	"github.com/Skifskii/goph-keeper/cmd/goph-keeper-client/screens"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Name is the canonical identifier for the register screen.
var Name = "register"

// Screen implements the TUI state and behavior for the register view.
type Screen struct {
	api        *api.APIClient
	options    []string
	username   textinput.Model
	password   textinput.Model
	focusIndex int
	err        string
	Done       bool
	Switch     bool
}

// NewScreen constructs a register Screen wired with an API client.
func NewScreen(apiClient *api.APIClient) *Screen {
	return &Screen{
		api: apiClient,
		options: []string{
			"username",
			"password",
			"submit",
			"signin",
		},
	}
}

// Init initializes input models and internal state for the register Screen.
func (l *Screen) Init() tea.Cmd {
	l.username = textinput.New()
	l.username.Placeholder = "__________"
	l.username.Prompt = ""
	l.username.Focus()

	l.password = textinput.New()
	l.password.Placeholder = "__________"
	l.password.Prompt = ""
	l.password.EchoMode = textinput.EchoPassword
	l.password.EchoCharacter = '•'

	l.Done = false
	l.Switch = false

	l.focusIndex = 0

	return nil
}

// Update processes incoming messages and updates the register Screen
// state. It returns the possibly modified screen and any command to run.
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
			switch l.focusIndex {

			// submit
			case 2:
				if !l.Done {
					err := l.api.Register(l.username.Value(), l.password.Value())
					if err != nil {
						l.err = err.Error()
						return l, nil
					} else {
						l.Done = true
					}
				}

			// login
			case 3:
				l.Switch = true
			}
		}
	}

	l.err = ""
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

// View renders the register Screen as a string for display by the TUI.
func (l Screen) View() string {
	var b strings.Builder

	// title
	b.WriteString(screens.LabelStyle.Render("    welcome to goph-keeper!"))
	b.WriteString("\n\n\n")

	// login
	b.WriteString("    REGISTER")
	b.WriteString("\n")
	b.WriteString(l.buildRow("username") + ": ")
	b.WriteString(l.username.View())
	b.WriteString("\n")
	b.WriteString(l.buildRow("password") + ": ")
	b.WriteString(l.password.View())
	b.WriteString("\n")
	if l.err != "" {
		b.WriteString(screens.ErrorStyle.Render(l.err))
		b.WriteString("\n")
	}
	b.WriteString(screens.ButtonSubmitStyle.Render(l.buildRow("submit")))
	b.WriteString("\n\n\n")

	// register
	b.WriteString(screens.ItalicStyle.Render("    Already have an account?"))
	b.WriteString("\n")
	b.WriteString(l.buildRow("signin"))
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
	case "signin":
		return "[ sign in ]"
	default:
		return "ERROR: unknown"
	}
}
