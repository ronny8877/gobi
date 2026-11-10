// SPDX-License-Identifier: MIT

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const maxBodyBytes = 1 << 20

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Powered-By", "Gobi")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func arrayLength(value string) (int, error) {
	name, args, err := parseValueBwBrackets(value)
	if err != nil || name != "Array" {
		return 0, fmt.Errorf("responseType must be Array(n)")
	}
	n, err := strconv.Atoi(args)
	if err != nil || n < 0 || n > maxArrayLength {
		return 0, fmt.Errorf("array length must be between 0 and %d", maxArrayLength)
	}
	return n, nil
}
func delay(r *http.Request, ms *int) bool {
	if ms == nil || *ms <= 0 {
		return true
	}
	timer := time.NewTimer(time.Duration(*ms) * time.Millisecond)
	defer timer.Stop()
	select {
	case <-timer.C:
		return true
	case <-r.Context().Done():
		return false
	}
}

func newHandler(a *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		snapshot := a.snapshot()
		if r.URL.Path == "/health" {
			writeJSON(w, 200, map[string]string{"status": "ok"})
			return
		}
		if r.URL.Path == "/auth" {
			serveMockAuth(w, r, &snapshot.Config)
			return
		}
		snapshot.serveAPI(w, r)
	})
}

func serveMockAuth(w http.ResponseWriter, r *http.Request, c *AppConfig) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		writeError(w, 405, "method not allowed")
		return
	}
	// This endpoint deliberately returns mock credentials; never configure real secrets here.
	if c.Auth == nil {
		writeError(w, 500, "no mock auth configured")
		return
	}
	if c.Auth.Cookie != nil && *c.Auth.Cookie != "" {
		parts := strings.SplitN(*c.Auth.Cookie, "=", 2)
		if len(parts) != 2 || parts[0] == "" {
			writeError(w, 500, "invalid cookie configuration")
			return
		}
		cookie := &http.Cookie{Name: parts[0], Value: parts[1], Path: "/", HttpOnly: true, SameSite: http.SameSiteLaxMode}
		if err := cookie.Valid(); err != nil {
			writeError(w, 500, "invalid cookie configuration")
			return
		}
		http.SetCookie(w, cookie)
	}
	if c.Auth.ApiKey != nil {
		w.Header().Set("X-API-Key", *c.Auth.ApiKey)
	}
	writeJSON(w, 200, map[string]any{"apiKey": c.Auth.ApiKey, "bearerToken": c.Auth.BearerToken, "cookie": c.Auth.Cookie})
}

func (a *App) serveAPI(w http.ResponseWriter, r *http.Request) {
	a.Logger().info("%s %s", r.Method, r.URL.Path)
	if !delay(r, a.Config.Latency) {
		return
	}
	if a.Config.FailRate != nil && rand.Float64() < float64(*a.Config.FailRate) {
		writeError(w, 500, "simulated failure")
		return
	}
	allowed := []string{}
	for i := range a.APIs {
		api := &a.APIs[i]
		if !matchPath(a.Config.Prefix+api.Path, r) {
			continue
		}
		if api.Method != r.Method {
			allowed = append(allowed, api.Method)
			continue
		}
		if !delay(r, api.Latency) {
			return
		}
		if api.FailRate != nil && rand.Float64() < *api.FailRate {
			writeError(w, 500, "simulated failure")
			return
		}
		if api.Auth != nil && (api.Auth.Protected == nil || *api.Auth.Protected) {
			ok, err := auth(&a.Config, api, r)
			if err != nil || !ok {
				writeError(w, 401, "unauthorized")
				return
			}
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
		if api.Validate != nil {
			if api.Validate.Query != nil {
				if _, err := validateQuery(*api.Validate.Query, r.URL.Query()); err != nil {
					writeError(w, 400, err.Error())
					return
				}
			}
			if api.Validate.Body != nil {
				body, err := io.ReadAll(r.Body)
				if err != nil {
					writeError(w, 413, "request body too large")
					return
				}
				if _, err := validateBody(*api.Validate.Body, io.NopCloser(bytes.NewReader(body))); err != nil {
					writeError(w, 400, err.Error())
					return
				}
				r.Body = io.NopCloser(bytes.NewReader(body))
			}
		}
		if api.PocketBase != nil {
			a.servePocketBase(w, r, api)
			return
		}
		if api.ResponseType != nil {
			n, err := arrayLength(*api.ResponseType)
			if err != nil {
				writeError(w, 500, err.Error())
				return
			}
			result := make([]any, 0, n)
			for j := 0; j < n; j++ {
				result = append(result, responseBuilderWithRefs(api.Response, a.Ref, 0))
			}
			writeJSON(w, 200, result)
			return
		}
		writeJSON(w, 200, responseBuilderWithRefs(api.Response, a.Ref, 0))
		return
	}
	if len(allowed) > 0 {
		w.Header().Set("Allow", strings.Join(allowed, ", "))
		writeError(w, 405, "method not allowed")
		return
	}
	writeError(w, 404, "API not found")
}
