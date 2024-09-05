package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
)

type API struct {
	Method       string                 `json:"method"`
	Path         string                 `json:"path"`
	Latency      *int                   `json:"latency"`
	FailRate     *float64               `json:"failRate"`
	Validate     *Validate              `json:"validate"`
	ResponseType *string                `json:"responseType"`
	Response     map[string]interface{} `json:"response"`
}

type Validate struct {
	Query  []string               `json:"query"`
	Body   map[string]interface{} `json:"body"`
	Params map[string]interface{} `json:"params"`
}

type AppInput struct {
	Config AppConfig `json:"config"`
	APIs   []API     `json:"api"`
}

type AppConfig struct {
	Prefix      string   `json:"prefix"`
	HealthCheck *string  `json:"healthCheck"`
	Port        int      `json:"port"`
	Latency     *int     `json:"latency"`
	Logging     *bool    `json:"logging"`
	FailRate    *float32 `json:"failRate"`
	Timeout     *int32   `json:"timeout"`
	ApiKey      *string  `json:"apiKey"`
}

type App struct {
	Config AppConfig `json:"config"`
	APIs   []API     `json:"api"`
}

var defaultLatency = 0
var defaultFailRate = float32(0.4)
var defaultTimeout = int32(5000)
var healthCheck = "/health"
var defaultPort = 8080
var defaultLogging = false

var defaultConfig = AppConfig{
	Prefix:      "/api",
	HealthCheck: &healthCheck, // Optional field
	Port:        8080,
	Logging:     &defaultLogging,  // Optional field
	Latency:     &defaultLatency,  // Optional field
	FailRate:    &defaultFailRate, // Optional field
	Timeout:     &defaultTimeout,  // Optional field
	ApiKey:      nil,              // Optional field
}

var logger = Logger(true)

func loadConfig(app *AppInput) {
	// Read the file
	file, err := os.ReadFile("input.json")
	if err != nil {
		logger.err("Error reading file:", err)
	}

	// Check if the file is empty
	if len(file) == 0 {
		fmt.Println(string(file))
		logger.err("Error: input.json is empty")
	}

	// Create a new instance of AppInput
	//This is a workaround as i had issues where if the key was removed from input the value was not getting updated
	var newApp AppInput

	// Unmarshal the JSON into the new instance
	err = json.Unmarshal(file, &newApp)
	if err != nil {
		log.Fatal("Error unmarshalling JSON:", err)
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

	if app.Config.HealthCheck == nil {
		app.Config.HealthCheck = defaultConfig.HealthCheck
	}

	if app.Config.Latency == nil {
		app.Config.Latency = defaultConfig.Latency
	}

	logger.debug("Configuration Loaded: ")
}

func watchConfigFile(filename string, app *AppInput) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Error creating watcher:", err)
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
					loadConfig(app)
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

// TODO Implement the timeout function
// TODO Implement auth	 using API key cookie or bearer token
func main() {
	var app AppInput
	loadConfig(&app)
	go watchConfigFile("input.json", &app)
	validateBody()
	//starting a server with the port
	go startServer(&app)
	logger.debug("Server is running on http://localhost:%d%s/\n", app.Config.Port, app.Config.Prefix)
	http.ListenAndServe(fmt.Sprintf(":%d", app.Config.Port), nil)
	// Wait for interrupt signal to gracefully shutdown the application
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	fmt.Println("Shutting down...")
}

// func validateQueryParams(query map[string]interface{}) bool {
// 	return true
// }

func startServer(app *AppInput) {
	http.HandleFunc(fmt.Sprintf("%s/", app.Config.Prefix), func(w http.ResponseWriter, r *http.Request) {
		// if logging is enabled, log the request
		if app.Config.Logging != nil && *app.Config.Logging {
			fmt.Printf("%s %s %s\n", r.Method, r.URL.Path, r.Proto)
		}

		//If global latency is provided, sleep for that time
		if app.Config.Latency != nil || *app.Config.Latency != 0 {
			fmt.Printf("Adding Latency %d ms\n", *app.Config.Latency)
			time.Sleep(time.Duration(*app.Config.Latency) * time.Millisecond)
		}
		//If global fail rate is provided, return an error if the random number is less than the fail rate
		if app.Config.FailRate != nil {
			random := rand.Float64()
			if random < float64(*app.Config.FailRate) {
				http.Error(w, `{"error":"Internal server Error"}`, http.StatusInternalServerError)
				return
			}
		}

		var path = r.URL.Path
		//check if the path is in the APIs
		var found = false

		for _, api := range app.APIs {
			if strings.HasPrefix(path, app.Config.Prefix) && api.Path == strings.TrimPrefix(path, app.Config.Prefix) && api.Method == r.Method {
				found = true
				//check if the latency is provided
				if api.Latency != nil && *api.Latency != 0 {
					//sleep for the latency
					fmt.Printf("Adding Latency %d ms\n", *api.Latency)
					time.Sleep(time.Duration(*api.Latency) * time.Millisecond)

				}
				//check if the fail rate is provided
				if api.FailRate != nil && *api.FailRate != 0 {
					//generate a random number between 0 and 1
					//if the number is less than the fail rate, return an error
					random := rand.Float64()
					if random < *api.FailRate {
						http.Error(w, `{"error":"Internal server Error"}`, http.StatusInternalServerError)
						return
					}
				}

				//check if the validate is provided
				if api.Validate != nil {
					//check if the query parameters are provided
					if api.Validate.Query != nil {
						// fmt.Println("Query Params: ", api.Validate.Query)
						// fmt.Println("Query Params: ", r.URL.Query())
						_, err := validateQueryParam(api.Validate.Query, r.URL.Query())
						if err != nil {
							http.Error(w, `{"error": "Invalid Query Params"}`, http.StatusBadRequest)
							return
						}
						if app.Config.Logging != nil && *app.Config.Logging {
							fmt.Println("Valid Query Params")
						}

					}
					//check if the body is provided
					if api.Validate.Body != nil {
						//check if the body is correct
					}
					//check if the params are provided
					if api.Validate.Params != nil {
						//check if the params are correct.
					}
				}
				//return the response
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("X-Powered-By", "Gobi")
				// fmt.Println("Response Type: ", string(*api.ResponseType))
				if api.ResponseType != nil {
					respType, arg := parseValueBwBrackets(*api.ResponseType)
					response := []interface{}{}
					if respType == "Array" {
						arrLenInt, _ := strconv.Atoi(arg)
						for i := 0; i < arrLenInt; i++ {
							response = append(response, ResponseBuilder(api.Response))
						}
						json.NewEncoder(w).Encode(response)
						break
					}
				} else {
					response := ResponseBuilder(api.Response)
					json.NewEncoder(w).Encode(response)
					break
				}
				// fmt.Println("Response: ", response)
				//Send the response and stop the loop

			}
		}
		if !found {
			http.Error(w, `{"error": "API not found"}`, http.StatusNotFound)
		}

	})
}

func parseValueBwBrackets(value string) (string, string) {
	var name = strings.Split(value, "(")[0]
	var args = strings.Split(value, "(")[1]
	args = strings.TrimRight(args, ")")
	return name, args
}

func validateQueryParam(query []string, rQuery map[string][]string) (bool, error) {
	//  if both dose not have the same length then return false
	if len(query) != len(rQuery) {
		return false, fmt.Errorf("mismatch in query params")
	}
	for _, q := range query {
		if _, ok := rQuery[q]; !ok {
			return false, fmt.Errorf("mismatch in query params")
		}
	}
	return true, nil
}
