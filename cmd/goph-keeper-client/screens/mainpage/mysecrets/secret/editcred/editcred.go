package editcred_screen

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Skifskii/goph-keeper/cmd/goph-keeper-client/api"
	"github.com/Skifskii/goph-keeper/cmd/goph-keeper-client/screens"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Name is the canonical identifier for the edit credential screen.
var Name = "editcred"

// Screen implements the UI for editing a credential secret.
type Screen struct {
	api        *api.APIClient
	secretID   int
	options    []string
	metadata   textinput.Model
	login      textinput.Model
	password   textinput.Model
	focusIndex int
	err        string
	Done       bool
	Quit       bool
}

// NewScreen constructs an edit credential Screen wired with the API client.
func NewScreen(apiClient *api.APIClient) *Screen {
	return &Screen{
		api: apiClient,
		options: []string{
			"metadata",
			"login",
			"password",
			"submit",
		},
	}
}

// InitWitParams initializes the edit screen state with existing secret
// parameters and payload. It prepares input fields for editing.
func (m *Screen) InitWitParams(secretID int, metadata string, payload json.RawMessage) tea.Cmd {
	var p CredentialPayload
	json.Unmarshal(payload, &p)

	m.secretID = secretID

	m.metadata = textinput.New()
	m.metadata.Placeholder = "__________"
	m.metadata.Prompt = ""
	m.metadata.SetValue(metadata)
	m.metadata.Focus()

	m.login = textinput.New()
	m.login.Placeholder = "__________"
	m.login.Prompt = ""
	m.login.SetValue(p.Login)

	m.password = textinput.New()
	m.password.Placeholder = "__________"
	m.password.Prompt = ""
	m.password.SetValue(p.Password)

	m.focusIndex = 0
	m.err = ""
	m.Done = false
	m.Quit = false

	return nil
}

// Init implements the Bubble Tea Init hook for the edit screen.
func (m *Screen) Init() tea.Cmd {
	return nil
}

// Update processes events for the edit screen and triggers updates
// to the secret when requested.
func (m *Screen) Update(msg tea.Msg) (*Screen, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "up":
			if m.focusIndex > 0 {
				m.focusIndex--
			}
			m.updateFocus()

		case "down":
			if m.focusIndex < len(m.options)-1 {
				m.focusIndex++
			}
			m.updateFocus()

		case "enter":
			// submit
			if m.focusIndex == 3 {
				if !m.Done {
					err := m.EditCredentialSecret(m.login.Value(), m.password.Value(), m.metadata.Value())
					if err != nil {
						m.err = err.Error()
						return m, nil
					} else {
						m.Done = true
					}
				}
			}

		case "esc":
			m.Done = true
			m.Quit = true
			return m, nil
		}
	}

	m.err = ""
	cmd := m.updateInputs(msg)

	return m, cmd
}

// CredentialPayload represents the JSON structure used for credential
// secrets in edit operations.
type CredentialPayload struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// EditCredentialSecret builds and sends an update for the given secret
// id using the API client.
func (m *Screen) EditCredentialSecret(login, password, metadata string) error {
	payload := CredentialPayload{
		Login:    login,
		Password: password,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	secret := api.Secret{
		Metadata: metadata,
		Payload:  jsonPayload,
	}

	err = m.api.UpdateSecret(m.secretID, secret)
	if err != nil {
		return fmt.Errorf("failed to create secret: %w", err)
	}

	return nil
}

func (m *Screen) updateFocus() {
	m.metadata.Blur()
	m.login.Blur()
	m.password.Blur()

	switch m.focusIndex {
	case 0:
		m.metadata.Focus()
	case 1:
		m.login.Focus()
	case 2:
		m.password.Focus()
	}
}

func (m *Screen) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 3)

	// Only update text inputs if they're focused
	switch m.focusIndex {
	case 0:
		m.metadata, cmds[0] = m.metadata.Update(msg)
	case 1:
		m.login, cmds[0] = m.login.Update(msg)
	case 2:
		m.password, cmds[1] = m.password.Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m Screen) View() string {
	var b strings.Builder

	// title
	b.WriteString(screens.LabelStyle.Render(fmt.Sprintf("    goph-keeper / main page / my secrets / secret '%d' / edit", m.secretID)))
	b.WriteString("\n\n\n")

	// metadata
	b.WriteString(m.buildRow("metadata", 0) + ": ")
	b.WriteString(m.metadata.View())
	b.WriteString("\n\n\n")

	// credentials
	b.WriteString(m.buildRow("login", 1) + ": ")
	b.WriteString(m.login.View())
	b.WriteString("\n")
	b.WriteString(m.buildRow("password", 2) + ": ")
	b.WriteString(m.password.View())
	b.WriteString("\n")
	if m.err != "" {
		b.WriteString(screens.ErrorStyle.Render(m.err))
		b.WriteString("\n")
	}
	b.WriteString(screens.ButtonSubmitStyle.Render(m.buildRow("[ submit ]", 3)))
	b.WriteString("\n\n\n")

	b.WriteString(screens.HelpStyle.Render("    Use '↑', '↓' and 'Enter' to navigate, 'Esc' to exit"))
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
