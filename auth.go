package main

import (
	"errors"
	"net/http"
	"strings"
)

func auth(config *AppConfig, api *API, r *http.Request) (bool, error) {

	if *api.Auth.ProtectedBy == "apiKey" {
		if r.Header.Get("X-API-Key") == *config.Auth.ApiKey || r.URL.Query().Get("apiKey") == *config.Auth.ApiKey {
			return true, nil
		} else {
			return false, errors.New("invalid API Key")
		}
	}

	if *api.Auth.ProtectedBy == "bearer" {
		bearerToken := r.Header.Get("Authorization")
		if bearerToken == "" {
			return false, errors.New("bearer token is missing")
		}

		// Split the token and compare
		parts := strings.Split(bearerToken, " ")
		if len(parts) != 2 || parts[0] != "Bearer" || parts[1] != *config.Auth.BearerToken {
			return false, errors.New("bearer token is invalid")
		}

		return true, nil
	}

	if *api.Auth.ProtectedBy == "cookie" {
		cookies := r.Cookies()
		for _, cookie := range cookies {
			if cookie.Name == "auth" {
				if cookie.Value != strings.Split(*config.Auth.Cookie, "=")[1] {
					return false, errors.New("cookie is invalid")
				}
				// TODO : Add More validation for cookie
				// if cookie.Expires.Before(time.Now()) {
				// 	return false, errors.New("cookie is expired")
				// }
				return true, nil
			}
		}
		return false, errors.New("auth_cookie is missing")
	}

	// Other authentication methods (apiKey, bearer)...
	return false, errors.New("unsupported authentication method")

}
