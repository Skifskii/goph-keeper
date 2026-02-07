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
	editcred_screen "github.com/Skifskii/goph-keeper/cmd/goph-keeper-client/screens/mainpage/mysecrets/secret/editcred"
	newsecret_screen "github.com/Skifskii/goph-keeper/cmd/goph-keeper-client/screens/mainpage/newsecret"
	newcred_screen "github.com/Skifskii/goph-keeper/cmd/goph-keeper-client/screens/mainpage/newsecret/newcred"
	newtext_screen "github.com/Skifskii/goph-keeper/cmd/goph-keeper-client/screens/mainpage/newsecret/newtext"
	register_screen "github.com/Skifskii/goph-keeper/cmd/goph-keeper-client/screens/register"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	currentScreen string

	login    *login_screen.Screen
	register *register_screen.Screen

	mainPage  *mainpage_screen.Screen
	mySecrets *mysecrets_screen.Screen
	secret    *secret_screen.Screen

	newSecret *newsecret_screen.Screen
	newCred   *newcred_screen.Screen
	newText   *newtext_screen.Screen

	editCred *editcred_screen.Screen
}

func initialModel() model {
	apiClient := api.NewClient("http://localhost:8080")

	login := login_screen.NewScreen(apiClient)
	login.Init()

	return model{
		currentScreen: login_screen.Name,

		login:    login,
		register: register_screen.NewScreen(apiClient),

		mainPage:  mainpage_screen.NewScreen(),
		mySecrets: mysecrets_screen.NewScreen(apiClient),
		secret:    secret_screen.NewScreen(apiClient),

		newSecret: newsecret_screen.NewScreen(),
		newCred:   newcred_screen.NewScreen(apiClient),
		newText:   newtext_screen.NewScreen(apiClient),

		editCred: editcred_screen.NewScreen(apiClient),
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

		if m.login.Switch {
			// change screen to "register"
			m.currentScreen = register_screen.Name
			return m, m.register.Init()
		}

		if m.login.Done {
			// change screen to "main page"
			m.currentScreen = mainpage_screen.Name
			return m, m.mainPage.Init()
		}

	// register
	case register_screen.Name:
		m.register, cmd = m.register.Update(msg)

		if m.register.Switch || m.register.Done {
			// change screen to "login"
			m.currentScreen = login_screen.Name
			return m, m.login.Init()
		}

	// main page
	case mainpage_screen.Name:
		m.mainPage, cmd = m.mainPage.Update(msg)
		if m.mainPage.Quit {
			// change screen to "login"
			m.currentScreen = login_screen.Name
			return m, m.login.Init()
		}

		if m.mainPage.Done {
			switch m.mainPage.NextScreen() {
			case mysecrets_screen.Name:
				// change screen to "my secrets"
				m.currentScreen = mysecrets_screen.Name
				return m, m.mySecrets.Init()
			case newsecret_screen.Name:
				// change screen to "new secret"
				m.currentScreen = newsecret_screen.Name
				return m, m.newSecret.Init()
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

	// secret
	case secret_screen.Name:
		m.secret, cmd = m.secret.Update(msg)

		if m.secret.Edit {
			// edit
			m.currentScreen = editcred_screen.Name
			return m, m.editCred.InitWitParams(m.secret.SecretID, m.secret.Metadata, m.secret.Payload)
		}

		if m.secret.Done {
			// quit
			m.currentScreen = mysecrets_screen.Name
			return m, m.mySecrets.Init()
		}

	// new secret
	case newsecret_screen.Name:
		m.newSecret, cmd = m.newSecret.Update(msg)

		if m.newSecret.Done {
			if m.newSecret.Quit {
				m.currentScreen = mainpage_screen.Name
				return m, m.mainPage.Init()
			}

			switch m.newSecret.NextScreen() {
			case newcred_screen.Name:
				// change screen to "new cred"
				m.currentScreen = newcred_screen.Name
				return m, m.newCred.Init()

			case newtext_screen.Name:
				// change screen to "new text"
				m.currentScreen = newtext_screen.Name
				return m, m.newText.Init()
			}
		}

	// new cred
	case newcred_screen.Name:
		m.newCred, cmd = m.newCred.Update(msg)

		if m.newCred.Done {
			// change screen to "new cred"
			m.currentScreen = mysecrets_screen.Name
			return m, m.mySecrets.Init()
		}

	// new text
	case newtext_screen.Name:
		m.newText, cmd = m.newText.Update(msg)

		if m.newText.Done {
			// change screen to "new cred"
			m.currentScreen = mysecrets_screen.Name
			return m, m.mySecrets.Init()
		}

	// edit cred
	case editcred_screen.Name:
		m.editCred, cmd = m.editCred.Update(msg)

		if m.editCred.Done {
			// change screen to "new cred"
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
	case register_screen.Name:
		return m.register.View()

	case mainpage_screen.Name:
		return m.mainPage.View()
	case mysecrets_screen.Name:
		return m.mySecrets.View()
	case secret_screen.Name:
		return m.secret.View()

	case newsecret_screen.Name:
		return m.newSecret.View()
	case newcred_screen.Name:
		return m.newCred.View()
	case newtext_screen.Name:
		return m.newText.View()

	case editcred_screen.Name:
		return m.editCred.View()

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
