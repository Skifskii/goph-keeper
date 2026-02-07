package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Skifskii/goph-keeper/cmd/goph-keeper-client/api"
	login_screen "github.com/Skifskii/goph-keeper/cmd/goph-keeper-client/screens/login"
	mainpage_screen "github.com/Skifskii/goph-keeper/cmd/goph-keeper-client/screens/mainpage"
	mysecrets_screen "github.com/Skifskii/goph-keeper/cmd/goph-keeper-client/screens/mainpage/mysecrets"
	secret_screen "github.com/Skifskii/goph-keeper/cmd/goph-keeper-client/screens/mainpage/mysecrets/secret"
	newsecret_screen "github.com/Skifskii/goph-keeper/cmd/goph-keeper-client/screens/mainpage/newsecret"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	currentScreen string
	login         *login_screen.Screen
	mainPage      *mainpage_screen.Screen
	mySecrets     *mysecrets_screen.Screen
	secret        *secret_screen.Screen
}

func initialModel() model {
	apiClient := api.NewClient("http://localhost:8080")

	return model{
		currentScreen: login_screen.Name,

		login:     login_screen.NewScreen(apiClient),
		mainPage:  mainpage_screen.NewScreen(),
		mySecrets: mysecrets_screen.NewScreen(apiClient),
		secret:    secret_screen.NewScreen(apiClient),
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

	// login
	case login_screen.Name:
		m.login, cmd = m.login.Update(msg)

		if m.login.Done {
			// change screen to "main page"
			m.currentScreen = mainpage_screen.Name
			return m, m.mainPage.Init()
		}

	// main page
	case mainpage_screen.Name:
		m.mainPage, cmd = m.mainPage.Update(msg)

		if m.mainPage.Done {
			switch m.mainPage.NextScreen() {
			case mysecrets_screen.Name:
				// change screen to "my secrets"
				m.currentScreen = mysecrets_screen.Name
				return m, m.mySecrets.Init()
			case newsecret_screen.Name:
				// change screen to "new secret"
				return m, nil // TODO:
			}
		}

	// my secrets
	case mysecrets_screen.Name:
		m.mySecrets, cmd = m.mySecrets.Update(msg)

		if m.mySecrets.Done {
			// quit
			if m.mySecrets.SelectedID == "" {
				m.currentScreen = mainpage_screen.Name
				return m, m.mainPage.Init()
			}

			// open secret
			selectedID, _ := strconv.Atoi(m.mySecrets.SelectedID) // TODO: check error
			m.currentScreen = secret_screen.Name
			return m, m.secret.SetSecretID(selectedID)
		}

	case secret_screen.Name:
		m.secret, cmd = m.secret.Update(msg)

		if m.secret.Done {
			// quit
			m.currentScreen = mysecrets_screen.Name
			return m, m.mySecrets.Init()
		}
	}

	return m, cmd
}

func (m model) View() string {
	switch m.currentScreen {
	case login_screen.Name:
		return m.login.View()
	case mainpage_screen.Name:
		return m.mainPage.View()
	case mysecrets_screen.Name:
		return m.mySecrets.View()
	case secret_screen.Name:
		return m.secret.View()

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
