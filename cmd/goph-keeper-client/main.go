package main

import (
	"fmt"
	"os"

	"github.com/Skifskii/goph-keeper/cmd/goph-keeper-client/api"
	"github.com/Skifskii/goph-keeper/cmd/goph-keeper-client/screens"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	currentScreen   string
	loginScreen     *screens.LoginScreen
	mainPageScreen  *screens.MainPageScreen
	mySecretsScreen *screens.MySecretsScreen
}

func initialModel() model {
	apiClient := api.NewClient("http://localhost:8080")

	return model{
		currentScreen:   "login",
		loginScreen:     screens.NewLoginScreen(apiClient),
		mainPageScreen:  screens.NewMainPageScreen(),
		mySecretsScreen: screens.NewMySecretsScreen(apiClient),
	}
}

func (m model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "ctrl+c" {
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	switch m.currentScreen {
	case screens.ScreenLogin:
		m.loginScreen, cmd = m.loginScreen.Update(msg)
		if m.loginScreen.Done {
			m.currentScreen = screens.ScreenMainPage
			return m, m.mainPageScreen.Init()
		}
	case screens.ScreenMainPage:
		m.mainPageScreen, cmd = m.mainPageScreen.Update(msg)
		if m.mainPageScreen.Done {
			nextScreen := m.mainPageScreen.NextScreen()
			m.currentScreen = nextScreen
			if nextScreen == screens.ScreenMySecrets {
				return m, m.mySecretsScreen.Init()
			}
		}
	case screens.ScreenMySecrets:
		m.mySecretsScreen, cmd = m.mySecretsScreen.Update(msg)
		if m.mySecretsScreen.Done {
			m.currentScreen = screens.ScreenMainPage
			return m, m.mainPageScreen.Init()
		}
	}

	return m, cmd
}

func (m model) View() string {
	switch m.currentScreen {
	case screens.ScreenLogin:
		return m.loginScreen.View()
	case screens.ScreenMainPage:
		return m.mainPageScreen.View()
	case screens.ScreenMySecrets:
		return m.mySecretsScreen.View()
	default:
		return fmt.Sprintf("unknown screen: '%s'", m.currentScreen)
	}
}

func main() {
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
