package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/common-nighthawk/go-figure"
)

var defaultConfigFilename = "gobi.json"

const (
	padding  = 2
	maxWidth = 80
)

var (
	textStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	mainMenuStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	focusedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	checkMark     = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
	crossMark     = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).SetString("✗")
	errorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
)

type model struct {
	Screen           string
	Choices          []string
	Index            int
	Chosen           string
	Loading          bool
	ConfigLoaded     bool
	Quitting         bool
	Progress         float64
	Help             help.Model
	Error            error
	Spinner          spinner.Model
	Cursor           int
	TextInput        textinput.Model
	IsInputValid     bool
	CreateNewOptions []string
}

type Config struct {
	Active string   `json:"active"`
	Files  []string `json:"files"`
}

var config Config

type loadError struct {
	err error
}
type loadSuccess struct{}

func initialModel() model {
	_, err := os.Stat(defaultConfigFilename)
	if os.IsNotExist(err) {
		file, ok := os.Create(defaultConfigFilename)
		if ok != nil {
			fmt.Println("Error creating config file:", ok)
			return model{
				Loading: true,
				Error:   ok,
			}
		}
		defer file.Close()
		file.WriteString(`{"files": []}`)
	}

	h := help.New()
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 60
	ti.CursorStyle = focusedStyle

	return model{
		Choices:          []string{"Create new", "Open Folder", "Select from List"},
		CreateNewOptions: []string{"Empty Project", "Project With Config", "Project With Example"},
		Cursor:           1,
		Help:             h,
		Spinner:          s,
		Screen:           "main",
		TextInput:        ti,
		IsInputValid:     true,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(loadConfig, m.Spinner.Tick, textinput.Blink)
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
	if strings.HasSuffix(fileName, ".gobi.json") && len(strings.Split(fileName, ".")) == 3 {
		return true
	}

	return false
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
				}
				return m, nil
			})
		}
		if m.Screen == "createChoices" {
			return updateMenu(msg.String(), m, m.CreateNewOptions, func(m model) (model, tea.Cmd) {
				m.Chosen = m.CreateNewOptions[m.Cursor]
				switch m.Chosen {
				case "Empty Project":
					m.Loading = true
					// createFile(m.TextInput.Value())
					updateFilesList(m.TextInput.Value())
					m.Loading = false
					return m, nil
				}
				return m, nil
			})
		}

		if strings.ToLower(m.Screen) == "create" {
			if k == "enter" {
				if isValidFileName(m.TextInput.Value()) {
					m.Screen = "createChoices"
					m.Cursor = 0
					m.IsInputValid = true
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

	if msg, ok := msg.(loadError); ok {
		m.Loading = false
		m.Error = msg.err
		return m, nil
	}
	if _, ok := msg.(loadSuccess); ok {
		m.Loading = false
		m.ConfigLoaded = true
		return m, nil
	}

	m.Spinner, cmd = m.Spinner.Update(spinner.TickMsg{})
	return m, cmd

}

func errorView(err error) string {
	return fmt.Sprintf("❌ Error: %v", err)
}

func createNewApi(m model) string {
	s := ""
	if len(m.TextInput.Value()) > 2 && !m.IsInputValid {
		s = errorStyle.Render("Invalid file name " + crossMark.Render())
	}
	m.TextInput.Placeholder = "api.gobi.json"
	return fmt.Sprintf(
		"Name the new Project?\n Either D://Code//api.gobi.json or api.gobi.json \n\n%s\n\n%s\n\n%s\n%s",
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

// func mainView(m model) string {
// 	s := "Hi, How would you like to proceed?\n\n"

// 	// Iterate over our choices
// 	for i, choice := range m.Choices {

// 		// Is the cursor pointing at this choice?
// 		cursor := " " // no cursor
// 		if m.Cursor == i {
// 			cursor = ">" // cursor!
// 		}

// 		// Render the row
// 		s += fmt.Sprintf("%s  %s\n", cursor, choice)
// 	}

// 	// Send the UI for rendering
// 	return s

// }

func (m model) View() string {

	if m.Loading {
		return m.Spinner.View() + textStyle.Render(" Loading...")
	}
	if m.Screen == "main" {
		s := welcomeMessage() + "\n\n" + listView("Hi, How would you like to proceed?\n\n", m.Choices, m)
		return mainMenuStyle.Render(s)
	}
	if strings.ToLower(m.Screen) == "create" {
		return createNewApi(m)
	}
	if m.Screen == "createChoices" {
		s := m.TextInput.Value() + " " + checkMark.Render() + "\n\n" + listView("Get Started:\n\n", m.CreateNewOptions, m)
		return mainMenuStyle.Render(s)
	}
	if m.Error != nil {
		return errorView(m.Error)
	}

	if !m.ConfigLoaded {
		return errorView(fmt.Errorf("Config not loaded Something went wrong"))
	}

	return "\n\nPress q to quit"
}

func StartApp() {
	tea.NewProgram(initialModel()).Run()
}

func loadConfig() tea.Msg {

	// Read the config file
	file, ok := os.ReadFile(defaultConfigFilename)
	if ok != nil {
		fmt.Println("Error reading config file:", ok)
		return loadError{ok}
	}

	// Check if the file is empty
	if len(file) == 0 {
		fmt.Println("Error: Config file is empty")
		return loadError{fmt.Errorf("config file is empty")}
	}

	err := json.Unmarshal(file, &config)
	if err != nil {
		fmt.Println("Error unmarshalling config file:")
		return loadError{err}
	}

	return loadSuccess{}

}

func formatFileNames() []string {
	var files []string
	for _, file := range config.Files {
		//Extreact file from path
		file = strings.Split(file, "/")[1]
		files = append(files, file)
	}
	return files
}

func renderFileList() {
	enumeratorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("99")).MarginRight(1)
	itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("212")).MarginRight(1)
	list := list.New(formatFileNames()).Enumerator(list.Bullet).EnumeratorStyle(enumeratorStyle).ItemStyle(itemStyle)
	var style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	fmt.Println(style.Render("Files in the current directory:"))
	fmt.Println(list)
}

// Function to create a file
func createDefaultConfigFile(path string) error {
	if path == "" {
		path = "./" + defaultConfigFilename
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

// Function to load a file (for demonstration purposes, just reads the content)
func loadFile(path string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	fmt.Println("File content:\n", string(content))
	return nil
}

// Function to update the list of files
func updateFilesList(path string) {
	files, err := readFilesList()
	if err != nil {
		files = []string{}
	}
	files = append(files, path)
	saveFilesList(files)
}

// Read files list from config.json
func readFilesList() ([]string, error) {
	file, err := os.Open("config.json")
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}
	defer file.Close()
	var files []string
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&files); err != nil {
		return nil, err
	}
	return files, nil
}

// Save files list to config.json
func saveFilesList(files []string) {
	file, err := os.Create("config.json")
	if err != nil {
		fmt.Println("Error saving files list:", err)
		return
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(files); err != nil {
		fmt.Println("Error encoding files list:", err)
	}
}

func welcomeMessage() string {
	asciiArt := figure.NewFigure("Welcome to Gobi!", "", true).String()
	var style = lipgloss.NewStyle().
		Bold(true).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("228")).
		Padding(2)

	// Print the styled ASCII art
	return style.Render(asciiArt)
}
