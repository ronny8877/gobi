// SPDX-License-Identifier: MIT

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"unicode"
)

// PocketBaseConfig connects to a separately managed PocketBase instance.
// Tokens come from the environment and are never returned by Gobi's mock /auth.
type PocketBaseConfig struct {
	URL      string `json:"url"`
	TokenEnv string `json:"tokenEnv,omitempty"`
}
type PocketBaseRoute struct {
	Collection string `json:"collection"`
}

func validSegment(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '-' {
			return false
		}
	}
	return true
}
func (c *PocketBaseConfig) client() (*url.URL, error) {
	u, err := url.Parse(c.URL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return nil, fmt.Errorf("pocketbase.url must be an http(s) URL without credentials, query or fragment")
	}
	if c.TokenEnv != "" && os.Getenv(c.TokenEnv) == "" {
		return nil, fmt.Errorf("PocketBase token environment variable %s is empty", c.TokenEnv)
	}
	return u, nil
}

func (a *App) servePocketBase(w http.ResponseWriter, r *http.Request, api *API) {
	cfg := a.Config.PocketBase
	if cfg == nil {
		writeError(w, 500, "PocketBase is not configured")
		return
	}
	u, err := cfg.client()
	if err != nil {
		writeError(w, 500, "PocketBase configuration is invalid")
		return
	}
	params, _ := parsePathParams(a.Config.Prefix+api.Path, r)
	id := params["id"]
	if strings.Contains(api.Path, ":id") && !validSegment(id) {
		writeError(w, 400, "invalid record id")
		return
	}
	switch r.Method {
	case "GET":
	case "POST":
		if id != "" {
			writeError(w, 405, "POST requires a collection route")
			return
		}
	case "PATCH", "DELETE":
		if id == "" {
			writeError(w, 400, "record id required")
			return
		}
	default:
		writeError(w, 405, "PocketBase supports GET, POST, PATCH and DELETE")
		return
	}
	u.Path = strings.TrimRight(u.Path, "/") + "/api/collections/" + api.PocketBase.Collection + "/records"
	if id != "" {
		u.Path += "/" + id
	}
	// Only forward documented list parameters, never client authorization headers.
	query := url.Values{}
	for _, key := range []string{"page", "perPage", "sort", "filter", "expand", "fields", "skipTotal"} {
		if values, ok := r.URL.Query()[key]; ok {
			query[key] = values
		}
	}
	u.RawQuery = query.Encode()
	var body io.Reader
	if r.Method == "POST" || r.Method == "PATCH" {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, 413, "request body too large")
			return
		}
		var object map[string]json.RawMessage
		if json.Unmarshal(data, &object) != nil || object == nil {
			writeError(w, 400, "body must be a JSON object")
			return
		}
		body = bytes.NewReader(data)
	}
	request, err := http.NewRequestWithContext(r.Context(), r.Method, u.String(), body)
	if err != nil {
		writeError(w, 500, "cannot construct PocketBase request")
		return
	}
	request.Header.Set("Content-Type", "application/json")
	if cfg.TokenEnv != "" {
		request.Header.Set("Authorization", os.Getenv(cfg.TokenEnv))
	}
	client := &http.Client{Timeout: 15 * time.Second, CheckRedirect: func(_ *http.Request, _ []*http.Request) error { return http.ErrUseLastResponse }}
	response, err := client.Do(request)
	if err != nil {
		writeError(w, 502, "PocketBase is unavailable")
		return
	}
	defer response.Body.Close()
	if response.StatusCode >= 300 && response.StatusCode < 400 {
		writeError(w, 502, "unexpected PocketBase redirect")
		return
	}
	data, err := io.ReadAll(io.LimitReader(response.Body, (8<<20)+1))
	if err != nil || len(data) > 8<<20 {
		writeError(w, 502, "invalid PocketBase response")
		return
	}
	if response.StatusCode != http.StatusNoContent && !json.Valid(data) {
		writeError(w, 502, "invalid PocketBase JSON response")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Powered-By", "Gobi")
	w.WriteHeader(response.StatusCode)
	_, _ = w.Write(data)
}
