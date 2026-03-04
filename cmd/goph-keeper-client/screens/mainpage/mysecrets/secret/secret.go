package secret_screen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Skifskii/goph-keeper/cmd/goph-keeper-client/api"
	"github.com/Skifskii/goph-keeper/cmd/goph-keeper-client/screens"
	tea "github.com/charmbracelet/bubbletea"
)

// Name is the canonical identifier for the secret detail screen.
var Name = "secret"

// Screen implements the UI for viewing a single secret, including
// actions to edit or delete the secret.
type Screen struct {
	apiClient *api.APIClient
	SecretID  int

	secretType string
	Metadata   string
	Payload    json.RawMessage
	err        string
	Done       bool
	Edit       bool
	Deleted    bool
}

// NewScreen constructs a secret detail Screen wired with the API client.
func NewScreen(apiClient *api.APIClient) *Screen {
	return &Screen{
		apiClient: apiClient,
	}
}

// SetSecretID sets the secret ID to display and starts loading the
// secret data asynchronously.
func (s *Screen) SetSecretID(secretID int) tea.Cmd {
	s.SecretID = secretID

	s.secretType = ""
	s.Metadata = ""
	s.Payload = nil
	s.err = ""
	s.Done = false
	s.Edit = false
	s.Deleted = false

	return s.loadSecret(secretID)
}

// Init is a noop for the secret detail Screen and implements the
// Bubble Tea lifecycle.
func (s *Screen) Init() tea.Cmd {
	return nil
}

type secretLoadedMsg api.Secret

func (s *Screen) loadSecret(secretID int) tea.Cmd {
	return func() tea.Msg {
		secret, err := s.apiClient.GetSecret(secretID)
		if err != nil {
			return err
		}
		return secretLoadedMsg(secret)
	}
}

// Update processes incoming messages, including a loaded secret, and
// updates the screen state accordingly.
func (s *Screen) Update(msg tea.Msg) (*Screen, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case secretLoadedMsg:
		s.secretType = msg.SecretType
		s.Metadata = msg.Metadata
		s.Payload = msg.Payload

	case error:
		s.err = msg.Error()

	case tea.KeyMsg:
		switch msg.String() {

		case "e":
			s.Edit = true
			s.Done = true
			return s, nil

		case "d":
			err := s.apiClient.DeleteSecret(s.SecretID)
			if err != nil {
				s.err = err.Error()
				return s, nil
			} else {
				s.Deleted = true
			}

			return s, nil

		case "esc":
			s.Done = true
		}
	}

	return s, cmd
}

// View renders the secret details or deletion state for display.
func (s Screen) View() string {
	var b strings.Builder

	// title
	b.WriteString(screens.LabelStyle.Render(fmt.Sprintf("    goph-keeper / main page / my secrets / secret '%d'", s.SecretID)))
	b.WriteString("\n\n\n")

	if s.Deleted {
		b.WriteString("DELETED")
		b.WriteString("\n\n\n")

		b.WriteString(screens.HelpStyle.Render("    Use 'Esc' to exit"))
		return b.String()
	}

	// info
	b.WriteString(fmt.Sprintf("type: %s", s.secretType))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("metadata: %s", s.Metadata))
	b.WriteString("\n\n\n")

	// secret
	var pretty bytes.Buffer
	if err := json.Indent(&pretty, s.Payload, "", "\t"); err != nil {
		b.WriteString(string(s.Payload))
	} else {
		b.WriteString(pretty.String())
	}

	b.WriteString("\n\n\n")

	b.WriteString(screens.HelpStyle.Render("    Use 'e' to edit, 'd' to delete, 'Esc' to exit"))
	b.WriteString("\n")

	return b.String()
}
