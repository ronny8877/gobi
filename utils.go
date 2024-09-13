package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
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
		args = strings.TrimSpace(parts[1])
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

func parseArguments(args string) (map[string]string, error) {
	result := make(map[string]string)
	stack := []string{}
	currentKey := ""
	currentValue := ""
	inBrackets := false

	for i := 0; i < len(args); i++ {
		char := args[i]

		switch char {
		case '=':
			if !inBrackets {
				currentKey = strings.TrimSpace(currentValue)
				currentValue = ""
			} else {
				currentValue += string(char)
			}
		case ',':
			if !inBrackets {
				result[currentKey] = strings.TrimSpace(currentValue)
				currentKey = ""
				currentValue = ""
			} else {
				currentValue += string(char)
			}
		case '(':
			inBrackets = true
			stack = append(stack, currentValue)
			currentValue += string(char)
		case ')':
			currentValue += string(char)
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				inBrackets = false
			}
		default:
			currentValue += string(char)
		}
	}

	if currentKey != "" {
		result[currentKey] = strings.TrimSpace(currentValue)
	}

	if len(stack) > 0 {
		return nil, errors.New("mismatched brackets in argument string")
	}

	return result, nil
}

func getFilesList(path string) ([]string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var fileList []string
	for _, file := range files {
		if strings.Contains(file.Name(), "gobi") && strings.HasSuffix(file.Name(), ".json") {
			fileList = append(fileList, file.Name())
		}
	}
	return fileList, nil
}
