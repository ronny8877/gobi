package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/common-nighthawk/go-figure"
)

var defaultConfigFilename = "gobi.config.json"

var (
	textStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	mainMenuStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	focusedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	checkMark     = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
	crossMark     = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).SetString("✗")
	errorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	tableStyle    = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))
)

type model struct {
	Screen            string
	Choices           []string
	Index             int
	Chosen            string
	Loading           bool
	ConfigLoaded      bool
	Quitting          bool
	Help              help.Model
	Error             error
	Spinner           spinner.Model
	Cursor            int
	TextInput         textinput.Model
	IsInputValid      bool
	CreateNewOptions  []string
	SelectedFile      string
	StartServer       bool
	InputErrorMessage string
	Table             table.Model
}

type Config struct {
	Active string
	Files  []string `json:"files"`
}

var config Config = Config{
	Active: "",
	Files:  []string{},
}

const (
	defaultSchema = `{
        "config": {
            "Port": 8080
        },
        "ref": {},
        "api": []
    }`

	defaultSchemaWithExample = `{
        "config": {
            "Port": 8080
        },
        "ref": {},
        "api": [
            {
                "method": "GET",
                "path": "/example",
                "response": {
                    "message": "Hello, World!",
                    "name": "User(username)"
                }
            }
        ]
    }`

	defaultSchemaWithConfig = `{
        "config": {
            "Port": 8080,
            "latency": 200,
            "logging": true,
            "failRate": 0.5,
            "prefix": "/v2/api",
            "auth": {}
        },
        "ref": {},
        "api": []
    }`
)

type loadError struct {
	err error
}
type loadSuccess struct{}
type serverSuccessMsg struct {
	msg string
}

func initialModel() model {
	h := help.New()
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 60
	ti.CursorStyle = focusedStyle

	columns := []table.Column{
		{Title: "Id", Width: 4},
		{Title: "File", Width: 25},
		{Title: "Status", Width: 15},
		{Title: "Path", Width: 35},
	}

	t := table.New(table.WithColumns(columns), table.WithFocused(true))
	style := table.DefaultStyles()
	style.Header = style.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	style.Selected = style.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(style)

	return model{
		Choices:          []string{"Create new", "Open Folder", "Select from List"},
		CreateNewOptions: []string{"Empty Project", "Project With Config", "Project With Example"},
		Cursor:           1,
		Help:             h,
		Spinner:          s,
		Screen:           "main",
		TextInput:        ti,
		IsInputValid:     true,
		Table:            t,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.Spinner.Tick, textinput.Blink, loadConfig)
}

func navigateMenu(msg tea.Msg, m model, choices []string) (model, tea.Cmd) {
	switch msg {
	case "q", "esc", "ctrl+c":
		m.Quitting = true
		return m, tea.Quit
	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
		}

	// The "down" and "j" keys move the cursor down
	case "down", "j":
		if m.Cursor < len(choices)-1 {
			m.Cursor++
		}
	}
	return m, nil
}
func updateMenu(msg tea.Msg, m model, choices []string, selectFunction func(model) (model, tea.Cmd)) (model, tea.Cmd) {
	m, cmd := navigateMenu(msg, m, choices)
	switch msg {
	// The "enter" key and the spacebar (a literal space) toggle
	// the selected state for the item that the cursor is pointing at.
	case "enter", " ":
		return selectFunction(m)
	}
	return m, cmd
}

// checkIfPathExists checks if a given path exists.
func checkIfPathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func createApiFile(choice int, m model) error {
	file, err := os.Create(m.TextInput.Value())
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	switch choice {
	case 0:
		file.WriteString(defaultSchema)
	case 1:
		file.WriteString(defaultSchemaWithConfig)
	case 2:
		file.WriteString(defaultSchemaWithExample)
	}
	return nil
}

// isValidFileName verifies if the file name follows the convention anyname.gobi.json
// and if the provided path exists.
func isValidFileName(name string) bool {
	if name == "" {
		return false
	}

	// Split the name to separate the path and the file name
	pathParts := strings.Split(name, "/")
	if len(pathParts) == 1 {
		pathParts = strings.Split(name, "\\")
	}

	// Extract the file name and the path
	fileName := pathParts[len(pathParts)-1]
	path := strings.Join(pathParts[:len(pathParts)-1], "/")

	// Check if the path exists
	if path != "" && !checkIfPathExists(path) {
		return false
	}

	// Verify if the file name follows the convention anyname.gobi.json
	if strings.HasSuffix(fileName, ".gobi.json") && len(strings.Split(fileName, ".")) >= 3 {
		return true
	}

	return false
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func getFileName(path string) string {
	pathParts := strings.Split(path, "/")
	if len(pathParts) == 1 {
		pathParts = strings.Split(path, "\\")
	}
	return pathParts[len(pathParts)-1]
}

func makeTableData(model) []table.Row {
	var data []table.Row
	if len(config.Files) == 0 {
		return append(data, table.Row{"", "No files found", "", ""})
	}

	for i, file := range config.Files {
		exist := checkIfPathExists(file)
		status := checkMark.Render()
		fileName := getFileName(file)
		if !exist {
			status = crossMark.Render()
		}

		data = append(data, table.Row{strconv.Itoa(i), fileName, status, file})
	}
	return data
}

func closeCli() tea.Msg {
	return serverSuccessMsg{"Server started successfully"}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	// fmt.Println("Update", msg, m)
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.Quitting = true
			return m, tea.Quit
		}
		if m.Screen == "main" {
			return updateMenu(msg.String(), m, m.Choices, func(m model) (model, tea.Cmd) {
				m.Chosen = m.Choices[m.Cursor]
				switch m.Chosen {
				case "Create new":
					m.Screen = "create"
					return m, nil
				case "Open Folder":
					m.Screen = "open"
				case "Select from List":
					m.Screen = "list"
				}
				return m, tea.Batch(loadConfig)
			})
		}
		if m.Screen == "createChoices" {
			return updateMenu(msg.String(), m, m.CreateNewOptions, func(m model) (model, tea.Cmd) {
				m.Chosen = m.CreateNewOptions[m.Cursor]
				m.Loading = true
				createApiFile(m.Cursor, m)
				updateFilesList(m.TextInput.Value())
				m.Loading = false
				m.StartServer = true
				config.Active = m.TextInput.Value()
				m.TextInput.SetValue("")
				return m, tea.Batch(closeCli)
			})
		}
		if m.Screen == "open" {
			if k == "enter" {
				if isValidFileName(m.TextInput.Value()) {
					m.Loading = true
					if !fileExists(m.TextInput.Value()) {
						m.Loading = false
						m.IsInputValid = false
						// m.Error = fmt.Errorf("file does not exist")
						return m, nil
					}
					updateFilesList(m.TextInput.Value())
					m.SelectedFile = m.TextInput.Value()
					m.Cursor = 0
					m.IsInputValid = true
					config.Active = m.SelectedFile
					m.TextInput.SetValue("")
					m.StartServer = true
					m.Screen = ""
					return m, tea.Batch(closeCli)
				} else {
					m.IsInputValid = false
				}
				return m, nil
			}
			m.TextInput, cmd = m.TextInput.Update(msg)
			return m, cmd
		}

		if m.Screen == "list" {
			m.Table, cmd = m.Table.Update(msg)
			if k == "q" {
				m.Screen = "main"
			}

			if k == "enter" {
				if m.Table.SelectedRow()[0] == "" {
					return m, nil
				}
				m.SelectedFile = m.Table.SelectedRow()[3]
				m.Loading = true
				m.StartServer = true
				config.Active = m.SelectedFile
				return m, tea.Batch(
					tea.Printf("Starting with %s!", m.Table.SelectedRow()[1]),
					closeCli,
				)
			}

			return m, cmd

		}

		if m.StartServer {
			m.Loading = false
			m.StartServer = false

			return m, tea.Quit
		}

		if strings.ToLower(m.Screen) == "create" {
			if k == "enter" {
				if isValidFileName(m.TextInput.Value()) {
					m.Screen = "createChoices"
					m.Cursor = 0
					m.IsInputValid = true
					updateFilesList(m.TextInput.Value())
				} else {
					m.IsInputValid = false
				}

				return m, nil
			}
			m.TextInput, cmd = m.TextInput.Update(msg)
			return m, cmd
		}
	}
	// m.Error = fmt.Errorf("Unknown message: %v", msg)
	if m.Quitting {
		return m, tea.Quit
	}

	if m.Error != nil {
		m.Quitting = true
		return m, nil
	}

	if _, ok := msg.(serverSuccessMsg); ok {
		return m, tea.Batch(tea.ClearScreen, tea.Quit)
	}
	if msg, ok := msg.(loadError); ok {
		m.Loading = false
		m.Error = msg.err
		return m, nil
	}
	if _, ok := msg.(loadSuccess); ok {
		m.Loading = false
		m.ConfigLoaded = true
		m.Table.SetRows(makeTableData(m))
		return m, nil
	}

	m.Spinner, cmd = m.Spinner.Update(spinner.TickMsg{})
	return m, cmd

}

func errorView(err error) string {
	return fmt.Sprintf("❌ Error: %v", err)
}

func createNewApi(m model, heading string, hints string) string {
	s := ""
	if len(m.TextInput.Value()) > 2 && !m.IsInputValid {
		s = errorStyle.Render("Invalid file name or path" + crossMark.Render())
	}
	m.TextInput.Placeholder = "api.gobi.json"
	return fmt.Sprintf(
		"%s\n %s \n\n%s\n\n%s\n\n%s\n%s",
		heading,
		hints,
		m.TextInput.View(),
		textStyle.Render("Press enter to continue"),
		s,
		"\n(esc to quit)",
	) + "\n"

}

func listView(s string, l []string, m model) string {
	for i, choice := range l {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.Cursor == i {
			cursor = ">" // cursor!
		}

		s += fmt.Sprintf("%s  %s\n", cursor, choice)
	}

	return s
}

func (m model) View() string {
	if m.StartServer {
		return m.Spinner.View() + textStyle.Render(" Starting server Please wait...")
	}
	if m.Loading {
		return m.Spinner.View() + textStyle.Render(" Loading...")
	}

	if m.Screen == "main" {
		s := welcomeMessage() + "\n\n" + listView("Hi, How would you like to proceed?\n\n", m.Choices, m)
		return mainMenuStyle.Render(s)
	}
	if strings.ToLower(m.Screen) == "create" {
		return createNewApi(m, "Enter the name of the new API file:", "File name should be in the format anyname.gobi.json")
	}
	if m.Screen == "createChoices" {
		s := m.TextInput.Value() + " " + checkMark.Render() + "\n\n" + listView("Get Started:\n\n", m.CreateNewOptions, m)
		return mainMenuStyle.Render(s)
	}
	if m.Screen == "open" {
		return createNewApi(m, "Enter the Filepath to open:", "File path should be in the format ./anyname.gobi.json d:/path/anyname.gobi.json")
	}

	if m.Screen == "list" {
		return tableStyle.Render(m.Table.View()) + "\n\n" + textStyle.Render("Press q to quit")
	}

	if m.Error != nil {
		return errorView(m.Error)
	}
	if m.StartServer {
		log.Debug("Starting server Please wait...")
	}
	if !m.ConfigLoaded {
		return errorView(fmt.Errorf("Config not loaded Something went wrong"))
	}

	return "\n\nPress q to quit"
}

func startApp() error {
	_, err := tea.NewProgram(initialModel()).Run()
	if err != nil {
		return err
	}
	return nil
}

func loadConfig() tea.Msg {
	//default cache location
	homeDir, ok := os.UserHomeDir()

	if ok != nil {
		fmt.Println("Error getting home directory:", ok)
		return loadError{ok}
	}
	cache := filepath.Join(homeDir, ".cache", "gobi")
	// Read the config file

	// Ensure the cache directory exists
	makeErr := os.MkdirAll(cache, 0755)
	if makeErr != nil {
		tea.Println("Error creating cache directory:", makeErr)
		return loadError{makeErr}
	}
	// Define the config file path
	configFilePath := filepath.Join(cache, defaultConfigFilename)

	//CHECK IF CONFIG FILE EXISTS
	if !fileExists(configFilePath) {
		// Create the config file if it doesn't exist
		createErr := createDefaultConfigFile(configFilePath)
		if createErr != nil {
			tea.Println("Error creating config file:", createErr)
			return loadError{createErr}
		}
	}

	file, ok := os.ReadFile(configFilePath)
	if ok != nil {
		tea.Println("Error reading config file:", ok)
		return loadError{ok}
	}

	// Check if the file is empty
	if len(file) == 0 {
		tea.Println("Error: Config file is empty")
		return loadError{}
	}

	err := json.Unmarshal(file, &config)
	if err != nil {
		fmt.Println("Error unmarshalling config file:")
		return loadError{err}
	}

	return loadSuccess{}

}

// Function to create a file
func createDefaultConfigFile(path string) error {
	defaultConfig := Config{
		Active: "",
		Files:  []string{},
	}
	// Marshal the default config to JSON
	data, err := json.MarshalIndent(defaultConfig, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling default config:", err)
		return err
	}

	// Write the default config to the file
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		fmt.Println("Error writing default config file:", err)
	}
	return nil
}

// Function to update the list of files
func updateFilesList(path string) {
	//check if path is just a filename or a path
	//If you simply do api.gobi.json without providing a path we will have to get the current working directory
	// I did not plan ahead for this
	// It was supposed to be a mock server Now as i use the server i feel i should have this or that
	//Resulting in a lot of changes Well who cares :)
	if !filepath.IsAbs(path) {
		// Get the current working directory
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error getting current working directory:", err)
			return
		}
		path = filepath.Join(cwd, path)
	}

	// Normalize the path to handle different path separators
	normalizedPath := filepath.Clean(path)

	homeDir, err := os.UserHomeDir()

	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return
	}
	configFilePath := filepath.Join(homeDir, ".cache", "gobi", defaultConfigFilename)

	// Read the file as by this point we know it exists
	//TODO: HANDEL BETTER
	data, err := os.ReadFile(configFilePath)
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("Error reading config file:", err)
		return
	}

	// Unmarshal the data
	var config Config
	if len(data) > 0 {
		if err := json.Unmarshal(data, &config); err != nil {
			fmt.Println("Error unmarshalling config file:", err)
			return
		}
	}

	// Check if the path already exists in the list
	for _, existingPath := range config.Files {
		if existingPath == normalizedPath {
			fmt.Println("Path already exists in the config file")
			return
		}
	}

	// Append the new file to the list
	config.Files = append(config.Files, normalizedPath)

	// Marshal the data
	newData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling config file:", err)
		return
	}

	// Write the data back to the file
	if err := os.WriteFile(configFilePath, newData, 0755); err != nil {
		fmt.Println("Error writing config file:", err)
		return
	}

}

func welcomeMessage() string {
	asciiArt := figure.NewFigure("Gobi!", "", true).String()
	var style = lipgloss.NewStyle().
		Bold(true).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("228")).
		Padding(2)

	// Print the styled ASCII art
	return style.Render(asciiArt)
}
