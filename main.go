package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/fsnotify/fsnotify"
)

type ProtectedByType string

// Constants for ProtectedByType
const (
	APIKey      ProtectedByType = "apiKey"
	BearerToken ProtectedByType = "BearerToken"
	Cookies     ProtectedByType = "Cookies"
)

type APIAuth struct {
	Protected   *bool            `json:"protected"`
	ProtectedBy *ProtectedByType `json:"protectedBy"`
}

type API struct {
	Auth         *APIAuth               `json:"auth"` // Optional field
	Method       string                 `json:"method"`
	Path         string                 `json:"path"`
	Latency      *int                   `json:"latency"`
	FailRate     *float64               `json:"failRate"`
	Validate     *Validate              `json:"validate"`
	ResponseType *string                `json:"responseType"`
	Response     map[string]interface{} `json:"response"`
}

type Validate struct {
	Query *[]string `json:"query"`
	Body  *[]string `json:"body"`
}

type Auth struct {
	ApiKey      *string `json:"apiKey"`
	BearerToken *string `json:"bearer"`
	Cookie      *string `json:"cookie"`
}

type AppConfig struct {
	Auth     *Auth    `json:"auth"`
	Prefix   string   `json:"prefix"`
	Port     int      `json:"port"`
	Latency  *int     `json:"latency"`
	Logging  *bool    `json:"logging"`
	FailRate *float32 `json:"failRate"`
}

type App struct {
	Config AppConfig               `json:"config"`
	Ref    *map[string]interface{} `json:"ref"`
	APIs   []API                   `json:"api"`
}

var defaultLatency = 0
var defaultFailRate = float32(0.0)
var defaultPort = 8080
var defaultLogging = false
var filename = "api.gobi.json"

var defaultConfig = AppConfig{
	Prefix:   "/api",
	Port:     8080,
	Logging:  &defaultLogging,  // Optional field
	Latency:  &defaultLatency,  // Optional field
	FailRate: &defaultFailRate, // Optional field        // Optional field
}

var logger = Logger(true)

func fileExistsOrCreate() error {
	log.Debug("Checking if file exists...")
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		fmt.Println("File does not exist, creating a new one...")
		file, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("error creating file: %w", err)
		}
		defer file.Close()
		_, err = file.WriteString(defaultSchema)
		if err != nil {
			return fmt.Errorf("error writing to file: %w", err)
		}
	}
	return nil
}

func loadAppConfig(app *App) error {
	// Read the file
	file, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// Check if the file is empty
	if len(file) == 0 {
		return fmt.Errorf("error: %s is empty", filename)
	}

	// Create a new instance of AppInput
	var newApp App

	// Unmarshal the JSON into the new instance
	err = json.Unmarshal(file, &newApp)
	if err != nil {
		return fmt.Errorf("error parsing JSON")
	}

	// Copy the values from the new instance to the existing app instance
	*app = newApp

	// Set default values if necessary
	if int(app.Config.Port) == 0 {
		app.Config.Port = defaultPort
	}

	if app.Config.Prefix == "" {
		app.Config.Prefix = defaultConfig.Prefix
	}

	if app.Config.Latency == nil {
		app.Config.Latency = defaultConfig.Latency
	}
	logger.debug("Configuration Loaded: ")

	return nil
}

func watchConfigFile(filename string, app *App) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Error creating watcher:", err)
		return
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					logger.info("Config file modified, reloading...")
					time.Sleep(100 * time.Millisecond) // Add a small delay before reloading as the it was crashing without it
					if err := loadAppConfig(app); err != nil {
						logger.err("Error reloading config: ", err)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Fatal("Watcher error:", err)
			}
		}
	}()

	err = watcher.Add(filename)
	if err != nil {
		log.Fatal("Error adding file to watcher:", err)
	}
	<-done
}

var app App

func main() {
	ok := startApp()
	if ok != nil {
		log.Error("Error starting the app")
		return
	}

	list, _ := getFilesList(".")
	fmt.Println(list)
	fileExistsOrCreate()
	// Load the configuration
	loadAppConfig(&app)
	go watchConfigFile(filename, &app)

	// Start the server
	go startServer(&app)
	logger.debug("Server is running on http://localhost:%d%s/\n", app.Config.Port, app.Config.Prefix)
	http.ListenAndServe(fmt.Sprintf(":%d", app.Config.Port), nil)

	// Wait for interrupt signal to gracefully shutdown the application
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	fmt.Println("Shutting down...")
}
