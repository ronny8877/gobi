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
var defaultLogging = true

var defaultConfig = AppConfig{
	Prefix:   "/api",
	Port:     8080,
	Logging:  &defaultLogging,  // Optional field
	Latency:  &defaultLatency,  // Optional field
	FailRate: &defaultFailRate, // Optional field        // Optional field
}

var app App = App{
	Config: defaultConfig,
	APIs:   []API{},
}
var logger = app.Logger()

func (app *App) loadAppConfig(filename string) error {
	// Read the file
	file, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// Check if the file is empty
	if len(file) == 0 {
		return fmt.Errorf("error: %s is empty", filename)
	}

	// Unmarshal the JSON into the new instance
	err = json.Unmarshal(file, &app)

	if err != nil {
		return fmt.Errorf("error parsing JSON")
	}

	logger.info("config loaded successfully")

	return nil
}

func (app *App) watchConfigFile(filename string) {
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
					app.Config = defaultConfig
					app.APIs = []API{}
					app.Ref = nil
					if err := app.loadAppConfig(filename); err != nil {
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

func main() {
	// Channel to signal when the Bubble Tea interface has completed
	done := make(chan error, 1)

	// Run the Bubble Tea interface in a separate goroutine
	go func() {
		done <- startApp()
	}()

	// Wait for the Bubble Tea interface to complete
	if err := <-done; err != nil {
		log.Error("Error starting the app")
		return
	}

	if config.Active == "" {
		log.Error("No config file provided")
		return
	}
	serverSetup()
	// Wait for interrupt signal to gracefully shutdown the application
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	fmt.Println("Shutting down...")
}

func serverSetup() {
	filename := "api.gobi.json"
	path := config.Active
	if path == "" {
		path = path + "/api.gobi.json"
		filename = path
	}
	filename = path
	// Load the configuration
	if err := app.loadAppConfig(filename); err != nil {
		logger.err("Error loading config: ", err)
		return
	}
	go app.watchConfigFile(filename)

	go startServer(&app)
	logger.info("Server is running on http://localhost:%d%s/\n", app.Config.Port, app.Config.Prefix)
	http.ListenAndServe(fmt.Sprintf(":%d", app.Config.Port), nil)
}
