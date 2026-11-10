// SPDX-License-Identifier: MIT

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPocketBaseCRUD(t *testing.T) {
	t.Setenv("TEST_PB_TOKEN", "backend-secret")
	var stored map[string]any
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "backend-secret" {
			t.Error("backend credential not forwarded")
		}
		if r.Header.Get("X-API-Key") != "" {
			t.Error("mock auth leaked upstream")
		}
		switch r.Method {
		case "POST":
			if r.URL.Path != "/api/collections/users/records" {
				t.Error(r.URL.Path)
			}
			_ = json.NewDecoder(r.Body).Decode(&stored)
			stored["id"] = "record123"
			writeJSON(w, 200, stored)
		case "GET":
			if strings.HasSuffix(r.URL.Path, "/record123") {
				writeJSON(w, 200, stored)
			} else {
				if r.URL.Query().Get("page") != "2" || r.URL.Query().Get("apiKey") != "" {
					t.Error("incorrect query forwarding")
				}
				writeJSON(w, 200, map[string]any{"items": []any{stored}})
			}
		case "PATCH":
			var updates map[string]any
			_ = json.NewDecoder(r.Body).Decode(&updates)
			for k, v := range updates {
				stored[k] = v
			}
			writeJSON(w, 200, stored)
		case "DELETE":
			stored = nil
			w.WriteHeader(204)
		}
	}))
	defer backend.Close()
	config := fmt.Sprintf(`{"config":{"pocketbase":{"url":%q,"tokenEnv":"TEST_PB_TOKEN"}},"api":[{"method":"POST","path":"/users","validate":{"body":["name"]},"pocketbase":{"collection":"users"}},{"method":"GET","path":"/users","pocketbase":{"collection":"users"}},{"method":"GET","path":"/users/:id","pocketbase":{"collection":"users"}},{"method":"PATCH","path":"/users/:id","pocketbase":{"collection":"users"}},{"method":"DELETE","path":"/users/:id","pocketbase":{"collection":"users"}}]}`, backend.URL)
	a := testApp(t, config)
	for _, tc := range []struct {
		method, path, body, want string
		code                     int
	}{
		{"POST", "/api/users", `{"name":"Ronny"}`, `"id":"record123"`, 200},
		{"GET", "/api/users/record123", "", `"name":"Ronny"`, 200},
		{"PATCH", "/api/users/record123", `{"name":"Updated"}`, `"name":"Updated"`, 200},
		{"GET", "/api/users?page=2&apiKey=client-secret", "", `"items"`, 200},
		{"POST", "/api/users", `{"other":"name"}`, "missing", 400},
		{"PATCH", "/api/users/record123", `[]`, "JSON object", 400},
		{"DELETE", "/api/users/record123", "", "", 204},
	} {
		w := request(a, tc.method, tc.path, tc.body)
		if w.Code != tc.code || !strings.Contains(w.Body.String(), tc.want) {
			t.Fatalf("%s %s: %d %s", tc.method, tc.path, w.Code, w.Body)
		}
	}
	if stored != nil {
		t.Fatal("record not deleted")
	}
}
func TestPocketBaseFailures(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		body   string
		want   int
	}{{"upstream error", 403, `{"message":"Forbidden"}`, 403}, {"invalid JSON", 200, "not json", 502}, {"redirect", 302, "", 502}} {
		t.Run(tc.name, func(t *testing.T) {
			backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Location", "https://example.invalid")
				w.WriteHeader(tc.status)
				_, _ = io.WriteString(w, tc.body)
			}))
			defer backend.Close()
			a := testApp(t, fmt.Sprintf(`{"config":{"pocketbase":{"url":%q}},"api":[{"method":"GET","path":"/users","pocketbase":{"collection":"users"}}]}`, backend.URL))
			w := request(a, "GET", "/api/users", "")
			if w.Code != tc.want {
				t.Fatal(w.Code, w.Body)
			}
		})
	}
}
