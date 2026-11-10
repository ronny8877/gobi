// SPDX-License-Identifier: MIT

package main

import (
	"encoding/json"
	"fmt"
	"io"
)

func validateQuery(query []string, values map[string][]string) (bool, error) {
	for _, key := range query {
		if _, ok := values[key]; !ok {
			return false, fmt.Errorf("query field %s is missing", key)
		}
	}
	return true, nil
}
func validateBody(required []string, body io.ReadCloser) (bool, error) {
	defer body.Close()
	var object map[string]json.RawMessage
	decoder := json.NewDecoder(io.LimitReader(body, maxBodyBytes+1))
	if err := decoder.Decode(&object); err != nil || object == nil {
		return false, fmt.Errorf("body must be a JSON object")
	}
	if decoder.Decode(new(any)) != io.EOF {
		return false, fmt.Errorf("body must contain one JSON object")
	}
	for _, key := range required {
		if _, ok := object[key]; !ok {
			return false, fmt.Errorf("body field %s is missing", key)
		}
	}
	return true, nil
}
