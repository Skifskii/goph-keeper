package screens

import (
	"fmt"
	"strconv"

	"github.com/Skifskii/goph-keeper/cmd/goph-keeper-client/api"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MySecretsScreen struct {
	api     *api.APIClient
	table   table.Model
	loading bool
	err     string
	Done    bool
}

func NewMySecretsScreen(apiClient *api.APIClient) *MySecretsScreen {
	columns := []table.Column{
		{Title: "ID", Width: 10},
		{Title: "Type", Width: 12},
		{Title: "Metadata", Width: 30},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	styles := table.DefaultStyles()
	styles.Header = styles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true)
	styles.Selected = styles.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57"))

	t.SetStyles(styles)

	return &MySecretsScreen{
		api:   apiClient,
		table: t,
	}
}

func (s *MySecretsScreen) Init() tea.Cmd {
	s.Done = false
	s.loading = true
	return s.loadSecrets()
}

type secretsLoadedMsg []api.SecretMeta

func (s *MySecretsScreen) loadSecrets() tea.Cmd {
	return func() tea.Msg {
		secrets, err := s.api.ListSecrets(50, 0)
		if err != nil {
			return err
		}
		return secretsLoadedMsg(secrets)
	}
}

func (s *MySecretsScreen) Update(msg tea.Msg) (*MySecretsScreen, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case secretsLoadedMsg:
		s.loading = false

		rows := make([]table.Row, 0, len(msg))
		for _, secret := range msg {
			rows = append(rows, table.Row{
				strconv.Itoa(secret.ID),
				secret.SecretType,
				secret.Metadata,
			})
		}
		s.table.SetRows(rows)

	case error:
		s.loading = false
		s.err = msg.Error()

	case tea.KeyMsg:
		switch msg.String() {

		case "enter":
			selected := s.table.SelectedRow()
			fmt.Println("open secret:", selected[0])
			s.Done = true

		case "q", "esc":
			s.Done = true
		}
	}

	s.table, cmd = s.table.Update(msg)
	return s, cmd
}

func (s MySecretsScreen) View() string {
	if s.loading {
		return "\n  loading secrets...\n"
	}

	if s.err != "" {
		return errorStyle.Render(s.err)
	}

	return lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Render(s.table.View()) +
		"\n\n" +
		helpStyle.Render("  ↑ ↓ navigate   Enter open   q back")
}
