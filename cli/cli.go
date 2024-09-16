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

type model struct {
	Choices      []string
	Index        int
	Chosen       string
	Loading      bool
	ConfigLoaded bool
	Quitting     bool
	Progress     float64
	Help         help.Model
	Error        error
	Spinner      spinner.Model
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
				Loading: false,
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
	return model{
		Help:    h,
		Spinner: s,
		Chosen:  "welcome",
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(loadConfig, m.Spinner.Tick)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// fmt.Println("Update", msg, m)
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.Quitting = true
			return m, tea.Quit
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

	var cmd tea.Cmd
	m.Spinner, cmd = m.Spinner.Update(spinner.TickMsg{})
	return m, cmd

}

func errorView(err error) string {
	//timer := time.NewTimer(5 * time.Second)

	return fmt.Sprintf("❌ Error: %v", err)
}

func (m model) View() string {
	if m.Chosen == "welcome" {
		s := welcomeMessage()
		s += "\n\n"
		s += "Press q to quit"
		return s
	}
	if m.Loading {
		return m.Spinner.View() + " Loading..."
	}
	if m.Error != nil {
		return errorView(m.Error)
	}

	if !m.ConfigLoaded {
		return errorView(fmt.Errorf("Config not loaded Something went wrong"))
	}

	return "Press q to quit"
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
func createFile(path string) error {
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

func renderAppMenu() {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true).Margin(1).Align(lipgloss.Center)
	itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	fmt.Println(style.Render("Menu"))
	fmt.Println(itemStyle.Render("1. Create new api"))
	fmt.Println(itemStyle.Render("2. Load from path"))
	fmt.Println(itemStyle.Render("-------------------------------"))
	readFilesList()
	textarea := textinput.New()
	textarea.Placeholder = "Enter the path to the file (leave empty for default)"
	textarea.Focus()
	textarea.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	textarea.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	textarea.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("99"))
	//Wait for user input
	var choice int
	fmt.Scanln(&choice)
	switch choice {
	case 1:

		fmt.Println("Creating new api...")
		createFile("api.gobi.json")
		updateFilesList("api.gobi.json")
	case 2:
		fmt.Println("Loading from path...")
		loadFile("api.gobi.json")
	default:
		fmt.Println("Invalid choice")
	}
}
