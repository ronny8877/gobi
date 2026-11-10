// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestServerGracefulShutdown(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	_ = listener.Close()
	filename := filepath.Join(t.TempDir(), "config.json")
	data := fmt.Sprintf(`{"config":{"port":%d,"logging":false},"api":[{"method":"GET","path":"/slow","latency":300,"response":{"ok":true}}]}`, port)
	if err := os.WriteFile(filename, []byte(data), 0600); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- runServer(ctx, filename) }()
	client := &http.Client{Timeout: 5 * time.Second, Transport: &http.Transport{DisableKeepAlives: true}}
	url := fmt.Sprintf("http://127.0.0.1:%d", port)
	deadline := time.Now().Add(5 * time.Second)
	for {
		r, err := client.Get(url + "/health")
		if err == nil {
			_ = r.Body.Close()
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("server did not start", err)
		}
		time.Sleep(10 * time.Millisecond)
	}
	response := make(chan error, 1)
	go func() {
		r, err := client.Get(url + "/api/slow")
		if err == nil {
			defer r.Body.Close()
			_, err = io.ReadAll(r.Body)
			if r.StatusCode != 200 {
				err = fmt.Errorf("status %d", r.StatusCode)
			}
		}
		response <- err
	}()
	// Wait for the slow request to be accepted before initiating shutdown.
	time.Sleep(100 * time.Millisecond)
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(12 * time.Second):
		t.Fatal("shutdown timed out")
	}
	if err := <-response; err != nil {
		t.Fatal("in-flight request failed", err)
	}
}
