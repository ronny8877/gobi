// SPDX-License-Identifier: MIT

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func ptr[T any](v T) *T { return &v }
func testApp(t *testing.T, data string) *App {
	t.Helper()
	filename := filepath.Join(t.TempDir(), "test.json")
	if err := os.WriteFile(filename, []byte(data), 0600); err != nil {
		t.Fatal(err)
	}
	a := &App{}
	if err := a.loadAppConfig(filename); err != nil {
		t.Fatal(err)
	}
	return a
}
func request(a *App, method, path, body string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, path, strings.NewReader(body))
	w := httptest.NewRecorder()
	newHandler(a).ServeHTTP(w, r)
	return w
}
func TestRoutesAndDefaults(t *testing.T) {
	a := testApp(t, `{"api":[{"method":"get","path":"/users/:id","response":{"ok":true}}]}`)
	for _, tc := range []struct {
		method, path string
		status       int
	}{{"GET", "/health", 200}, {"GET", "/api/users/123", 200}, {"POST", "/api/users/123", 405}, {"GET", "/api/users", 404}, {"GET", "/unknown", 404}} {
		t.Run(tc.method+tc.path, func(t *testing.T) {
			w := request(a, tc.method, tc.path, "")
			if w.Code != tc.status {
				t.Fatalf("%d: %s", w.Code, w.Body)
			}
			if !json.Valid(w.Body.Bytes()) {
				t.Fatal("invalid JSON")
			}
		})
	}
	if a.Config.Port != 8080 || a.Config.Latency == nil {
		t.Fatal("missing defaults")
	}
}
func TestAuthRegression(t *testing.T) {
	for _, tc := range []struct {
		name    string
		config  *Auth
		method  ProtectedByType
		headers map[string]string
		want    int
	}{
		{"api key without cookie", &Auth{ApiKey: ptr("key")}, APIKey, map[string]string{"X-API-Key": "key"}, 200},
		{"missing config", nil, APIKey, nil, 401},
		{"empty configured key", &Auth{ApiKey: ptr("")}, APIKey, nil, 401},
		{"bearer", &Auth{BearerToken: ptr("token")}, BearerToken, map[string]string{"Authorization": "Bearer token"}, 200},
		{"legacy bearer", &Auth{BearerToken: ptr("token")}, ProtectedByType("BearerToken"), map[string]string{"Authorization": "Bearer token"}, 200},
		{"custom cookie name", &Auth{Cookie: ptr("session=value=tail")}, Cookies, map[string]string{"Cookie": "session=value=tail"}, 200},
		{"malformed cookie", &Auth{Cookie: ptr("broken")}, Cookies, nil, 401},
		{"wrong cookie", &Auth{Cookie: ptr("session=value")}, Cookies, map[string]string{"Cookie": "auth=value"}, 401},
	} {
		t.Run(tc.name, func(t *testing.T) {
			a := &App{Config: newDefaultConfig(), APIs: []API{{Method: "GET", Path: "/private", Auth: &APIAuth{Protected: ptr(true), ProtectedBy: &tc.method}}}}
			a.Config.Auth = tc.config
			r := httptest.NewRequest("GET", "/api/private", nil)
			for k, v := range tc.headers {
				r.Header.Set(k, v)
			}
			w := httptest.NewRecorder()
			newHandler(a).ServeHTTP(w, r)
			if w.Code != tc.want {
				t.Fatalf("%d: %s", w.Code, w.Body)
			}
		})
	}
	a := testApp(t, `{"api":[{"method":"GET","path":"/public","auth":{"protected":false,"protectedBy":"apiKey"}}]}`)
	if w := request(a, "GET", "/api/public", ""); w.Code != 200 {
		t.Fatal(w.Code)
	}
	a.Config.Auth = &Auth{ApiKey: ptr("key")}
	if w := request(a, "GET", "/auth", ""); w.Code != 200 || w.Header().Get("X-API-Key") != "key" {
		t.Fatal(w.Code, w.Body)
	}
	a.Config.Auth = &Auth{Cookie: ptr("broken")}
	if w := request(a, "GET", "/auth", ""); w.Code != 500 {
		t.Fatal(w.Code)
	}
}
func TestBodyValidation(t *testing.T) {
	a := testApp(t, `{"api":[{"method":"POST","path":"/users","validate":{"body":["name"]},"response":{"ok":true}}]}`)
	for _, tc := range []struct {
		body   string
		status int
	}{{`{"name":"Ronny"}`, 200}, {`{"name":null}`, 200}, {`{"other":"name"}`, 400}, {`{"nested":{"name":"x"}}`, 400}, {`[]`, 400}, {`null`, 400}, {`{"name":1} {}`, 400}, {`not json`, 400}, {`{"name":"` + strings.Repeat("x", maxBodyBytes) + `"}`, 413}} {
		if w := request(a, "POST", "/api/users", tc.body); w.Code != tc.status {
			t.Fatalf("body length %d: %d %s", len(tc.body), w.Code, w.Body)
		}
	}
}
func TestReloadIsAtomicAndConcurrent(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "config.json")
	good := []byte(`{"config":{"prefix":"/v2"},"api":[{"method":"GET","path":"/users","response":{"id":"Uuid()"}}]}`)
	if err := os.WriteFile(filename, good, 0600); err != nil {
		t.Fatal(err)
	}
	a := &App{}
	if err := a.loadAppConfig(filename); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filename, []byte(`{"api":`), 0600); err != nil {
		t.Fatal(err)
	}
	if err := a.loadAppConfig(filename); err == nil {
		t.Fatal("expected invalid config")
	}
	if w := request(a, "GET", "/v2/users", ""); w.Code != 200 {
		t.Fatal("lost working config")
	}
	if err := os.WriteFile(filename, good, 0600); err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 30; j++ {
				if w := request(a, "GET", "/v2/users", ""); w.Code != 200 {
					t.Error(w.Code)
				}
			}
		}()
	}
	for i := 0; i < 30; i++ {
		if err := a.loadAppConfig(filename); err != nil {
			t.Error(err)
		}
	}
	wg.Wait()
}
func TestInvalidConfig(t *testing.T) {
	for _, data := range []string{`null`, `[]`, `{"config":{"port":0}}`, `{"config":{"latency":-1}}`, `{"config":{"failRate":2}}`, `{"config":{"prefix":"bad"}}`, `{"api":[{"method":"GET","path":"/x","responseType":"Array(-1)"}]}`, `{"api":[{"method":"GET","path":"/x","responseType":"Array(10001)"}]}`, `{"api":[{"method":"GET","path":"x"}]}`, `{"api":[{"method":"GET","path":"/x"},{"method":"GET","path":"/x"}]}`, `{"api":[{"method":"GET","path":"/x","pocketbase":{"collection":"users"}}]}`} {
		filename := filepath.Join(t.TempDir(), "config.json")
		_ = os.WriteFile(filename, []byte(data), 0600)
		if err := new(App).loadAppConfig(filename); err == nil {
			t.Errorf("accepted %s", data)
		}
	}
}
func TestResponseArraysAndReferences(t *testing.T) {
	a := testApp(t, `{"ref":{"nested":{"user":{"id":"Uuid()"}}},"api":[{"method":"GET","path":"/x","response":{"ints":"Array(len=3,type=Int(min=7,max=7))","bools":"Array(len=2,type=Bool(true))","literal":["Uuid()",{"name":"User()"}],"ref":"Ref(ref.nested.user)","birthday":"User(birthday)"}}]}`)
	w := request(a, "GET", "/api/x", "")
	var result map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	for _, v := range result["ints"].([]any) {
		if v != float64(7) {
			t.Fatalf("lost numeric type: %v", v)
		}
	}
	for _, v := range result["bools"].([]any) {
		if v != true {
			t.Fatalf("lost bool type: %v", v)
		}
	}
	if result["literal"].([]any)[0] == "Uuid()" {
		t.Fatal("unexpanded literal array")
	}
	if result["ref"].(map[string]any)["id"] == "Uuid()" {
		t.Fatal("unexpanded reference")
	}
	birthday, err := time.Parse("2006-01-02", result["birthday"].(string))
	if err != nil || !birthday.Before(time.Now()) {
		t.Fatal("invalid birthday", result["birthday"])
	}
	a.APIs[0].ResponseType = ptr("Array(0)")
	if w = request(a, "GET", "/api/x", ""); strings.TrimSpace(w.Body.String()) != "[]" {
		t.Fatal(w.Body)
	}
	cycle := map[string]any{"user": map[string]any{"self": "Ref(ref.user)"}}
	encoded, err := json.Marshal(responseBuilderWithRefs(map[string]any{"user": "Ref(ref.user)"}, &cycle, 0))
	if err != nil || !bytes.Contains(encoded, []byte("maximum generator depth")) {
		t.Fatal(string(encoded), err)
	}
}
func TestCancelledLatency(t *testing.T) {
	a := &App{Config: newDefaultConfig()}
	a.Config.Latency = ptr(30000)
	r := httptest.NewRequest("GET", "/api/x", nil)
	ctx, cancel := context.WithCancel(r.Context())
	cancel()
	started := time.Now()
	newHandler(a).ServeHTTP(httptest.NewRecorder(), r.WithContext(ctx))
	if time.Since(started) > time.Second {
		t.Fatal("ignored cancellation")
	}
}
func TestExamples(t *testing.T) {
	t.Setenv("POCKETBASE_TOKEN", "test-token")
	paths, _ := filepath.Glob("examples/*.json")
	paths = append(paths, "api.gobi.json")
	for _, filename := range paths {
		t.Run(filename, func(t *testing.T) {
			if err := new(App).loadAppConfig(filename); err != nil {
				t.Fatal(err)
			}
		})
	}
}
func TestValidateBodyReadError(t *testing.T) {
	_, err := validateBody([]string{"name"}, io.NopCloser(strings.NewReader("")))
	if err == nil {
		t.Fatal("accepted empty body")
	}
}
