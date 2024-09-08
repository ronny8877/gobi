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
	case "bearer":
		return validateBearerToken(config, r)
	case "cookie":
		return validateCookie(config, r)
	default:
		return false, errors.New("unsupported authentication method")
	}

}

func validateAPIKey(config *AppConfig, r *http.Request) (bool, error) {
	if config.Auth == nil || config.Auth.ApiKey == nil {
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
	if config.Auth == nil || config.Auth.BearerToken == nil {
		return false, errors.New("bearer token not configured")
	}
	bearerToken := r.Header.Get("Authorization")
	if bearerToken == "" {
		return false, errors.New("bearer token is missing")
	}
	parts := strings.Split(bearerToken, " ")
	if len(parts) != 2 || parts[0] != "Bearer" || parts[1] != *config.Auth.BearerToken {
		return false, errors.New("bearer token is invalid")
	}
	return true, nil
}

func validateCookie(config *AppConfig, r *http.Request) (bool, error) {
	if config.Auth == nil || config.Auth.Cookie == nil {
		return false, errors.New("cookie not configured")
	}
	expectedCookieValue := strings.Split(*config.Auth.Cookie, "=")[1]
	for _, cookie := range r.Cookies() {
		if cookie.Name == "auth" {
			if cookie.Value != expectedCookieValue {
				return false, errors.New("cookie is invalid")
			}
			// TODO: Add more validation for cookie
			// if cookie.Expires.Before(time.Now()) {
			// 	return false, errors.New("cookie is expired")
			// }
			return true, nil
		}
	}
	return false, errors.New("auth cookie is missing")
}
