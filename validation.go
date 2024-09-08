package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

func validateQuery(query []string, rQuery map[string][]string) (bool, error) {
	for _, q := range query {
		if _, ok := rQuery[q]; !ok {
			return false, fmt.Errorf("object %s is missing", q)
		}
	}
	return true, nil
}

func validateBody(body []string, rBody io.ReadCloser) (bool, error) {
	defer rBody.Close()
	//TODO: Add validators Like
	// UUID()
	//String()
	//Email()
	//URL()
	// Read the body content
	bodyContent, err := ioutil.ReadAll(rBody)
	if err != nil {
		return false, fmt.Errorf("failed to read request body: %v", err)
	}

	bodyString := string(bodyContent)

	// Check if each required field is present in the body content
	for _, b := range body {
		if !strings.Contains(bodyString, b) {
			return false, fmt.Errorf("object %s is missing", b)
		}
	}

	return true, nil
}
