package newtext_screen

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Skifskii/goph-keeper/cmd/goph-keeper-client/api"
	"github.com/Skifskii/goph-keeper/cmd/goph-keeper-client/screens"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Name is the canonical identifier for the text creation screen.
var Name = "newtext"

// Screen implements the UI for creating a new text secret.
type Screen struct {
	api        *api.APIClient
	options    []string
	metadata   textinput.Model
	text       textinput.Model
	focusIndex int
	err        string
	Done       bool
	Quit       bool
}

// NewScreen constructs a new text creation Screen wired with the API
// client.
func NewScreen(apiClient *api.APIClient) *Screen {
	return &Screen{
		api: apiClient,
		options: []string{
			"metadata",
			"text",
			"submit",
		},
	}
}

// Init prepares input models and resets internal state.
func (m *Screen) Init() tea.Cmd {
	m.metadata = textinput.New()
	m.metadata.Placeholder = "__________"
	m.metadata.Prompt = ""
	m.metadata.Focus()

	m.text = textinput.New()
	m.text.Placeholder = "__________"
	m.text.Prompt = ""

	m.focusIndex = 0
	m.err = ""
	m.Done = false
	m.Quit = false

	return nil
}

// Update handles input events for the text creation Screen and may
// trigger secret creation.
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
			if m.focusIndex == 2 {
				if !m.Done {
					err := m.CreateTextSecret(m.text.Value(), m.metadata.Value())
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

// TextPayload represents the JSON structure for text secrets.
type TextPayload struct {
	Text string `json:"text"`
}

// CreateTextSecret builds and sends a text secret using the API client.
func (m *Screen) CreateTextSecret(text, metadata string) error {
	payload := TextPayload{
		Text: text,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	secret := api.Secret{
		SecretType: "text",
		Metadata:   metadata,
		Payload:    jsonPayload,
	}

	_, err = m.api.CreateSecret(secret)
	if err != nil {
		return fmt.Errorf("failed to create secret: %w", err)
	}

	return nil
}

func (m *Screen) updateFocus() {
	m.metadata.Blur()
	m.text.Blur()

	switch m.focusIndex {
	case 0:
		m.metadata.Focus()
	case 1:
		m.text.Focus()
	}
}

func (m *Screen) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 2)

	// Only update text inputs if they're focused
	switch m.focusIndex {
	case 0:
		m.metadata, cmds[0] = m.metadata.Update(msg)
	case 1:
		m.text, cmds[0] = m.text.Update(msg)
	}

	return tea.Batch(cmds...)
}

// View renders the create-text Screen for display.
func (m Screen) View() string {
	var b strings.Builder

	// title
	b.WriteString(screens.LabelStyle.Render("    goph-keeper / main page / new secret / text"))
	b.WriteString("\n\n\n")

	// metadata
	b.WriteString(m.buildRow("metadata", 0) + ": ")
	b.WriteString(m.metadata.View())
	b.WriteString("\n\n\n")

	// text
	b.WriteString(m.buildRow("text", 1) + ": ")
	b.WriteString(m.text.View())
	b.WriteString("\n")
	if m.err != "" {
		b.WriteString(screens.ErrorStyle.Render(m.err))
		b.WriteString("\n")
	}
	b.WriteString(screens.ButtonSubmitStyle.Render(m.buildRow("[ submit ]", 2)))
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
