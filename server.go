package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func startServer(app *App) {
	//So with even a minimal setup, we can check if the server is running
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Powered-By", "Gobi")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})
	//Mock Auth server
	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Powered-By", "Gobi")

		//If no auth is provided, return an error
		if app.Config.Auth == nil {
			http.Error(w, `{"error": "No Auth Configured"}`, http.StatusInternalServerError)
			return
		}

		//Send whatever creds we had in the response
		response := map[string]interface{}{
			"apiKey":      app.Config.Auth.ApiKey,
			"bearerToken": app.Config.Auth.BearerToken,
			"cookie":      app.Config.Auth.Cookie,
		}
		// attach the cookie in the header if it is provided
		if app.Config.Auth.Cookie != nil && *app.Config.Auth.Cookie != "" {
			cookieParts := strings.SplitN(*app.Config.Auth.Cookie, "=", 2)
			if len(cookieParts) == 2 {
				http.SetCookie(w, &http.Cookie{
					Name:  cookieParts[0],
					Value: cookieParts[1],
				})
			} else {
				// Made an oopsie
				http.Error(w, "Invalid cookie format", http.StatusInternalServerError)
				return
			}
		}
		//Attach the api key in the header if it is provided
		if app.Config.Auth.Cookie != nil && *app.Config.Auth.ApiKey != "" {
			w.Header().Set("X-API-Key", *app.Config.Auth.ApiKey)
		}
		json.NewEncoder(w).Encode(response)

	})

	http.HandleFunc(fmt.Sprintf("%s/", app.Config.Prefix), func(w http.ResponseWriter, r *http.Request) {

		// if logging is enabled, log the request
		if app.Config.Logging != nil && *app.Config.Logging {
			fmt.Printf("%s %s %s\n", r.Method, r.URL.Path, r.Proto)
		}

		//Global Latency
		if *app.Config.Latency != 0 {
			logger.debug("Adding Latency %d ms\n", *app.Config.Latency)
			time.Sleep(time.Duration(*app.Config.Latency) * time.Millisecond)
		}
		//Global Fail Rate
		if app.Config.FailRate != nil {
			random := rand.Float64()
			if random < float64(*app.Config.FailRate) {
				http.Error(w, `{"error":"Internal server Error"}`, http.StatusInternalServerError)
				return
			}
		}

		if len(app.APIs) == 0 {
			http.Error(w, `{"error": "No APIs found"}`, http.StatusNotFound)
			return
		}
		//check if the path is in the APIs
		var found = false
		for _, api := range app.APIs {
			if strings.HasPrefix(r.URL.Path, app.Config.Prefix) && matchPath(fmt.Sprint(app.Config.Prefix, api.Path), r) && api.Method == r.Method {
				//parse the path params
				if app.Config.Logging != nil && *app.Config.Logging {
					pathParams, err := parsePathParams(fmt.Sprint(app.Config.Prefix, api.Path), r)
					if err != nil {
						logger.err("Error parsing path params: ", err)
					}
					fmt.Println("Path Params: ", pathParams)
				}
				found = true

				// API-specific latency
				if api.Latency != nil && *api.Latency != 0 {
					logger.debug("Adding Latency %d ms\n", *api.Latency)
					time.Sleep(time.Duration(*api.Latency) * time.Millisecond)
				}

				// API-specific fail rate
				if api.FailRate != nil && *api.FailRate != 0 {
					random := rand.Float64()
					if random < *api.FailRate {
						http.Error(w, `{"error":"Internal server Error"}`, http.StatusInternalServerError)
						return
					}
				}

				// Authentication
				if app.Config.Auth != nil && api.Auth != nil && api.Auth.ProtectedBy != nil && *api.Auth.ProtectedBy != "" {
					_, err := auth(&app.Config, &api, r)
					if err != nil {
						http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
						return
					}
				}

				// Validation
				if api.Validate != nil {
					//Query Validation
					if api.Validate.Query != nil {
						_, err := validateQuery(*api.Validate.Query, r.URL.Query())
						if err != nil {
							http.Error(w, `{"error": "Invalid Query Params"}`, http.StatusBadRequest)
							return
						}
						if app.Config.Logging != nil && *app.Config.Logging {
							logger.debug("Valid Query Params")
						}
					}
					//Body Validation
					if api.Validate.Body != nil {
						_, err := validateBody(*api.Validate.Body, r.Body)
						if err != nil {
							http.Error(w, `{"error": "Invalid Body"}`, http.StatusBadRequest)
							return
						}

					}
				}
				// Response
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("X-Powered-By", "Gobi")
				//Response Type
				if api.ResponseType != nil {
					respType, arg, err := parseValueBwBrackets(*api.ResponseType)
					if err != nil {
						http.Error(w, `{"error": "Invalid Response Type"}`, http.StatusBadRequest)
						return
					}

					response := []interface{}{}
					if respType == "Array" {
						arrLenInt, _ := strconv.Atoi(arg)
						for i := 0; i < arrLenInt; i++ {
							response = append(response, responseBuilder(api.Response))
						}
						json.NewEncoder(w).Encode(response)
						break
					}
				} else {
					response := responseBuilder(api.Response)
					json.NewEncoder(w).Encode(response)
					break
				}
			}
		}
		if !found {
			http.Error(w, `{"error": "API not found"}`, http.StatusNotFound)
		}

	})
}
