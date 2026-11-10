// SPDX-License-Identifier: MIT

package main

import (
	"errors"
	"net/http"
	"strings"
)

func auth(config *AppConfig, api *API, r *http.Request) (bool, error) {
	if api.Auth == nil || api.Auth.ProtectedBy == nil {
		return false, errors.New("authentication method not specified")
	}

	switch *api.Auth.ProtectedBy {
	case "apiKey":
		return validateAPIKey(config, r)
	case "bearer", "BearerToken":
		return validateBearerToken(config, r)
	case "cookie", "Cookies":
		return validateCookie(config, r)
	default:
		return false, errors.New("unsupported authentication method")
	}

}

func validateAPIKey(config *AppConfig, r *http.Request) (bool, error) {
	if config.Auth == nil || config.Auth.ApiKey == nil || *config.Auth.ApiKey == "" {
		return false, errors.New("API key not configured")
	}
	apiKey := r.Header.Get("X-API-Key")
	if apiKey == "" {
		apiKey = r.URL.Query().Get("apiKey")
	}
	if apiKey == *config.Auth.ApiKey {
		return true, nil
	}
	return false, errors.New("invalid API key")
}

func validateBearerToken(config *AppConfig, r *http.Request) (bool, error) {
	if config.Auth == nil || config.Auth.BearerToken == nil || *config.Auth.BearerToken == "" {
		return false, errors.New("bearer token not configured")
	}
	bearerToken := r.Header.Get("Authorization")
	if bearerToken == "" {
		return false, errors.New("bearer token is missing")
	}
	parts := strings.Fields(bearerToken)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] != *config.Auth.BearerToken {
		return false, errors.New("bearer token is invalid")
	}
	return true, nil
}

func validateCookie(config *AppConfig, r *http.Request) (bool, error) {
	if config.Auth == nil || config.Auth.Cookie == nil {
		return false, errors.New("cookie not configured")
	}
	parts := strings.SplitN(*config.Auth.Cookie, "=", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return false, errors.New("invalid cookie configuration")
	}
	cookie, err := r.Cookie(parts[0])
	if err != nil || cookie.Value != parts[1] {
		return false, errors.New("cookie is invalid or missing")
	}
	return true, nil
}
