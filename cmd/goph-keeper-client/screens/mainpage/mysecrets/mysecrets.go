package mysecrets_screen

import (
	"fmt"
	"strconv"

	"github.com/Skifskii/goph-keeper/cmd/goph-keeper-client/api"
	"github.com/Skifskii/goph-keeper/cmd/goph-keeper-client/screens"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var Name = "mysecrets"

type Screen struct {
	api        *api.APIClient
	table      table.Model
	loading    bool
	err        string
	SelectedID string
	Done       bool

	limit   int
	offset  int
	hasMore bool
}

func NewScreen(apiClient *api.APIClient) *Screen {
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

	return &Screen{
		api:   apiClient,
		table: t,
	}
}

func (s *Screen) Init() tea.Cmd {
	s.SelectedID = ""
	s.Done = false
	s.loading = true
	s.err = ""

	s.limit = 10
	s.offset = 0
	s.hasMore = true

	s.table.SetRows(nil)
	s.table.SetCursor(0)

	return s.loadSecrets(s.offset)
}

type secretsLoadedMsg struct {
	secrets []api.SecretMeta
}

func (s *Screen) loadSecrets(offset int) tea.Cmd {
	return func() tea.Msg {
		secrets, err := s.api.ListSecrets(s.limit, offset)
		if err != nil {
			return err
		}
		return secretsLoadedMsg{secrets: secrets}
	}
}

func (s *Screen) Update(msg tea.Msg) (*Screen, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case secretsLoadedMsg:
		s.loading = false

		if len(msg.secrets) < s.limit {
			s.hasMore = false
		}

		rows := s.table.Rows()
		for _, secret := range msg.secrets {
			rows = append(rows, table.Row{
				strconv.Itoa(secret.ID),
				secret.SecretType,
				secret.Metadata,
			})
		}

		s.table.SetRows(rows)
		s.offset += len(msg.secrets)

	case error:
		s.loading = false
		s.err = msg.Error()

	case tea.KeyMsg:
		switch msg.String() {

		case "enter":
			selected := s.table.SelectedRow()
			fmt.Println("open secret:", selected[0])
			s.SelectedID = selected[0]
			s.Done = true

		case "esc":
			s.Done = true
		}
	}

	prevCursor := s.table.Cursor()
	s.table, cmd = s.table.Update(msg)

	if s.hasMore &&
		!s.loading &&
		s.table.Cursor() == len(s.table.Rows())-1 &&
		s.table.Cursor() != prevCursor {

		s.loading = true
		return s, s.loadSecrets(s.offset)
	}

	return s, cmd
}

func (s Screen) View() string {
	if s.loading && s.offset == 0 {
		return "\n  loading secrets...\n"
	}

	if s.err != "" {
		return screens.ErrorStyle.Render(s.err)
	}

	view := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Render(s.table.View())

	if s.loading && s.hasMore {
		view += "\n  loading more..."
	}

	view += "\n\n" +
		screens.HelpStyle.Render("    Use '↑', '↓' and 'Enter' to navigate, 'q' to exit")

	return view
}
