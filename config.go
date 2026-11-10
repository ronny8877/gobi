// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

const maxArrayLength = 10000

func (a *App) loadAppConfig(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}
	next := App{Config: newDefaultConfig()}
	if err := json.Unmarshal(data, &next); err != nil {
		return fmt.Errorf("parse config: %w", err)
	}
	if strings.TrimSpace(string(data)) == "null" {
		return fmt.Errorf("config must be a JSON object")
	}
	if err := next.validate(); err != nil {
		return err
	}
	a.mu.Lock()
	a.Config, a.Ref, a.APIs = next.Config, next.Ref, next.APIs
	a.mu.Unlock()
	return nil
}

func (a *App) validate() error {
	c := &a.Config
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}
	if c.Prefix != "" && (!strings.HasPrefix(c.Prefix, "/") || strings.ContainsAny(c.Prefix, "{}?# ")) {
		return fmt.Errorf("prefix must be a URL path")
	}
	c.Prefix = strings.TrimRight(c.Prefix, "/")
	if c.Latency == nil {
		v := 0
		c.Latency = &v
	}
	if *c.Latency < 0 || *c.Latency > 30000 {
		return fmt.Errorf("latency must be between 0 and 30000 ms")
	}
	if c.FailRate != nil && !validRate(float64(*c.FailRate)) {
		return fmt.Errorf("failRate must be between 0 and 1")
	}
	if c.PocketBase != nil {
		if _, err := c.PocketBase.client(); err != nil {
			return err
		}
	}
	seen := map[string]bool{}
	for i := range a.APIs {
		api := &a.APIs[i]
		api.Method = strings.ToUpper(api.Method)
		switch api.Method {
		case "GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS":
		default:
			return fmt.Errorf("api[%d]: invalid method", i)
		}
		if !strings.HasPrefix(api.Path, "/") || strings.ContainsAny(api.Path, "?# ") {
			return fmt.Errorf("api[%d]: path must start with /", i)
		}
		key := api.Method + " " + c.Prefix + api.Path
		if seen[key] {
			return fmt.Errorf("duplicate route %s", key)
		}
		seen[key] = true
		if c.Prefix+api.Path == "/health" || c.Prefix+api.Path == "/auth" {
			return fmt.Errorf("api[%d]: reserved route", i)
		}
		if api.Latency != nil && (*api.Latency < 0 || *api.Latency > 30000) {
			return fmt.Errorf("api[%d]: invalid latency", i)
		}
		if api.FailRate != nil && !validRate(*api.FailRate) {
			return fmt.Errorf("api[%d]: invalid failRate", i)
		}
		if api.ResponseType != nil {
			if _, err := arrayLength(*api.ResponseType); err != nil {
				return fmt.Errorf("api[%d]: %w", i, err)
			}
		}
		if api.PocketBase != nil {
			if c.PocketBase == nil {
				return fmt.Errorf("api[%d]: pocketbase config required", i)
			}
			if !validSegment(api.PocketBase.Collection) {
				return fmt.Errorf("api[%d]: invalid collection", i)
			}
		}
	}
	return nil
}

func validRate(v float64) bool { return !math.IsNaN(v) && v >= 0 && v <= 1 }

// Watch the directory so atomic editor saves don't detach the watcher.
func (a *App) watchConfigFile(ctx context.Context, filename string) {
	filename, err := filepath.Abs(filename)
	if err != nil {
		logger.err("Config watch: %v", err)
		return
	}
	w, err := fsnotify.NewWatcher()
	if err != nil {
		logger.err("Config watch: %v", err)
		return
	}
	defer w.Close()
	if err = w.Add(filepath.Dir(filename)); err != nil {
		logger.err("Config watch: %v", err)
		return
	}
	timer := time.NewTimer(time.Hour)
	if !timer.Stop() {
		<-timer.C
	}
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case e, ok := <-w.Events:
			if !ok {
				return
			}
			if filepath.Clean(e.Name) == filename && e.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename) != 0 {
				if !timer.Stop() {
					select {
					case <-timer.C:
					default:
					}
				}
				timer.Reset(100 * time.Millisecond)
			}
		case <-timer.C:
			oldPort := a.snapshot().Config.Port
			if err := a.loadAppConfig(filename); err != nil {
				logger.err("Config reload rejected; keeping previous configuration: %v", err)
			} else {
				logger.info("Configuration reloaded")
				if a.snapshot().Config.Port != oldPort {
					logger.warn("Port change takes effect after restart")
				}
			}
		case err, ok := <-w.Errors:
			if !ok {
				return
			}
			logger.err("Config watch: %v", err)
		}
	}
}
