package main

import (
	"fmt"
	"net/http"
	"strings"
)

func parseValueBwBrackets(value string) (string, string, error) {
	parts := strings.SplitN(value, "(", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid format: %s", value)
	}
	name := parts[0]
	args := strings.TrimRight(parts[1], ")")
	return name, args, nil
}

func processMap(data map[string]interface{}, funcMap map[string]func(*string) interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range data {
		switch v := value.(type) {
		case string:
			//Getting the value In bw the brackets also removing the brackets
			funcName, args, err := parseValueBwBrackets(v)
			if err != nil {
				// Keep the original value if parsing fails
				result[key] = v
				continue
			}

			if fn, ok := funcMap[funcName]; ok {
				result[key] = fn(&args)
			} else {
				result[key] = v
			}
		case map[string]interface{}:
			result[key] = processMap(v, funcMap)
		default:
			result[key] = v
		}
	}
	return result
}

func parsePathParams(path string, r *http.Request) (map[string]string, error) {
	pathParams := make(map[string]string)
	pathParts := strings.Split(path, "/")
	requestParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) != len(requestParts) {
		return nil, fmt.Errorf("path and request path length mismatch")
	}
	for i, part := range pathParts {
		if strings.HasPrefix(part, ":") {
			pathParams[part[1:]] = requestParts[i]
		}
	}
	return pathParams, nil
}

func matchPath(path string, r *http.Request) bool {
	pathParts := strings.Split(path, "/")
	requestParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) != len(requestParts) {
		return false
	}
	for i, part := range pathParts {
		if strings.HasPrefix(part, ":") {
			continue
		}
		if part != requestParts[i] {
			return false
		}
	}
	return true
}
func formatAvatarURL(seed, avatarType string) string {
	return fmt.Sprintf("https://api.dicebear.com/9.x/%s/svg?seed=%s", avatarType, seed)
}

// func generateUuid() string {
// 	return uuid.New().String()
// }

// func generateBearerToken() (string, error) {
// 	return jwt.New(jwt.SigningMethodHS256).SignedString([]byte("secret"))
// }

// func createCookie(auth_id string) *http.Cookie {
// 	// Craft a new cookie with the greeting "Hello world!" and additional attributes.
// 	newCookie := http.Cookie{
// 		Name:     "auth_id",
// 		Value:    auth_id,
// 		Path:     "/",
// 		MaxAge:   3600,
// 		HttpOnly: true,
// 		Secure:   true,
// 		SameSite: http.SameSiteLaxMode,
// 	}
// 	return &newCookie
// 	// Dispatch the cookie to the client using the http.SetCookie() method.
// }
