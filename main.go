// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type ProtectedByType string

// Constants for ProtectedByType
const (
	APIKey      ProtectedByType = "apiKey"
	BearerToken ProtectedByType = "bearer"
	Cookies     ProtectedByType = "cookie"
)

type APIAuth struct {
	Protected   *bool            `json:"protected"`
	ProtectedBy *ProtectedByType `json:"protectedBy"`
}

type API struct {
	PocketBase   *PocketBaseRoute       `json:"pocketbase,omitempty"`
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
	PocketBase *PocketBaseConfig `json:"pocketbase,omitempty"`
	Auth       *Auth             `json:"auth"`
	Prefix     string            `json:"prefix"`
	Port       int               `json:"port"`
	Latency    *int              `json:"latency"`
	Logging    *bool             `json:"logging"`
	FailRate   *float32          `json:"failRate"`
}

type App struct {
	mu     sync.RWMutex
	Config AppConfig               `json:"config"`
	Ref    *map[string]interface{} `json:"ref"`
	APIs   []API                   `json:"api"`
}

var app = App{Config: newDefaultConfig()}
var logger = app.Logger()

func newDefaultConfig() AppConfig {
	latency, rate, logging := 0, float32(0), true
	return AppConfig{Prefix: "/api", Port: 8080, Latency: &latency, FailRate: &rate, Logging: &logging}
}

func (a *App) snapshot() *App {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return &App{Config: a.Config, Ref: a.Ref, APIs: a.APIs}
}

func main() {
	filename := flag.String("config", "", "path to a Gobi JSON configuration (skips the interactive menu)")
	flag.Parse()
	if *filename == "" {
		if err := startApp(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		*filename = config.Active
	}
	if *filename == "" {
		return
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := runServer(ctx, *filename); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runServer(ctx context.Context, filename string) error {
	ctx, cancelWatch := context.WithCancel(ctx)
	defer cancelWatch()
	if err := app.loadAppConfig(filename); err != nil {
		return err
	}
	cfg := app.snapshot().Config
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	server := &http.Server{Handler: newHandler(&app), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 30 * time.Second, WriteTimeout: 90 * time.Second, IdleTimeout: 60 * time.Second}
	go app.watchConfigFile(ctx, filename)
	stopped := make(chan struct{})
	shutdownDone := make(chan struct{})
	go func() {
		defer close(shutdownDone)
		select {
		case <-ctx.Done():
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if err := server.Shutdown(shutdownCtx); err != nil {
				_ = server.Close()
			}
		case <-stopped:
		}
	}()
	logger.info("Server listening on %s%s", listener.Addr(), cfg.Prefix)
	serveErr := server.Serve(listener)
	close(stopped)
	<-shutdownDone
	if serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
		return serveErr
	}
	return nil
}
