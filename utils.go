// SPDX-License-Identifier: MIT

package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func parseValueBwBrackets(value string) (string, string, error) {
	// Trim any leading or trailing whitespace from the input value
	value = strings.TrimSpace(value)

	// Split the value into parts using the first occurrence of '('
	parts := strings.SplitN(value, "(", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid format: %s", value)
	}
	// Trim any leading or trailing whitespace from the name part
	name := strings.TrimSpace(parts[0])
	// Trim any trailing whitespace and the closing ')' from the args part
	var args string
	if strings.HasSuffix(parts[1], ")") {
		args = strings.TrimSpace(parts[1][:len(parts[1])-1])
	} else {
		return "", "", fmt.Errorf("missing closing parenthesis: %s", value)
	}

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
		case []interface{}:
			result[key] = processSlice(v, funcMap)
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
			if requestParts[i] == "" {
				return false
			}
			continue
		}
		if part != requestParts[i] {
			return false
		}
	}
	return true
}
func formatAvatarURL(seed, avatarType string) string {
	return fmt.Sprintf("https://api.dicebear.com/9.x/%s/svg?seed=%s", avatarType, url.QueryEscape(seed))
}

// parseArguments splits only at top-level commas and the first equals sign.
func parseArguments(args string) (map[string]string, error) {
	result := make(map[string]string)
	if strings.TrimSpace(args) == "" {
		return result, nil
	}
	depth, start := 0, 0
	add := func(part string) error {
		key, value, ok := strings.Cut(part, "=")
		key = strings.TrimSpace(key)
		if !ok || key == "" {
			return errors.New("arguments must use name=value")
		}
		if _, exists := result[key]; exists {
			return fmt.Errorf("duplicate argument %s", key)
		}
		result[key] = strings.TrimSpace(value)
		return nil
	}
	for i, r := range args {
		switch r {
		case '(':
			depth++
		case ')':
			depth--
			if depth < 0 {
				return nil, errors.New("mismatched brackets in argument string")
			}
		case ',':
			if depth == 0 {
				if err := add(args[start:i]); err != nil {
					return nil, err
				}
				start = i + 1
			}
		}
	}
	if depth != 0 {
		return nil, errors.New("mismatched brackets in argument string")
	}
	if err := add(args[start:]); err != nil {
		return nil, err
	}
	return result, nil
}

func processSlice(values []interface{}, funcs map[string]func(*string) interface{}) []interface{} {
	result := make([]interface{}, len(values))
	for i, value := range values {
		switch v := value.(type) {
		case map[string]interface{}:
			result[i] = processMap(v, funcs)
		case []interface{}:
			result[i] = processSlice(v, funcs)
		case string:
			result[i] = processMap(map[string]interface{}{"value": v}, funcs)["value"]
		default:
			result[i] = value
		}
	}
	return result
}
