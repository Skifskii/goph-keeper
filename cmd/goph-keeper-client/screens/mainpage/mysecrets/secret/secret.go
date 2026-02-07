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

var Name = "secret"

type Screen struct {
	apiClient *api.APIClient
	secretID  int

	secretType string
	metadata   string
	payload    json.RawMessage
	err        string
	Done       bool
}

func NewScreen(apiClient *api.APIClient) *Screen {
	return &Screen{
		apiClient: apiClient,
	}
}

func (s *Screen) SetSecretID(secretID int) tea.Cmd {
	s.secretID = secretID

	s.secretType = ""
	s.metadata = ""
	s.payload = nil
	s.err = ""
	s.Done = false

	return s.loadSecret(secretID)
}

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

func (s *Screen) Update(msg tea.Msg) (*Screen, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case secretLoadedMsg:
		s.secretType = msg.SecretType
		s.metadata = msg.Metadata
		s.payload = msg.Payload

	case error:
		s.err = msg.Error()

	case tea.KeyMsg:
		switch msg.String() {

		case "q", "esc":
			s.Done = true
		}
	}

	return s, cmd
}

func (s Screen) View() string {
	var b strings.Builder

	// title
	b.WriteString(screens.LabelStyle.Render("    goph-keeper / main page / my secrets / secret"))
	b.WriteString("\n\n\n")

	// info
	b.WriteString(fmt.Sprintf("type: %s", s.secretType))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("metadata: %s", s.metadata))
	b.WriteString("\n\n\n")

	// secret
	var pretty bytes.Buffer
	if err := json.Indent(&pretty, s.payload, "", "\t"); err != nil {
		b.WriteString(string(s.payload))
	} else {
		b.WriteString(pretty.String())
	}

	b.WriteString("\n\n\n")

	b.WriteString(screens.HelpStyle.Render("    Use 'q' to exit"))
	b.WriteString("\n")

	return b.String()
}
