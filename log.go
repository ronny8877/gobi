package main

import (
	"fmt"

	"github.com/mattn/go-colorable"
)

// LoggerStruct defines the structure for the logger
type LoggerStruct struct {
	info  func(string, ...interface{})
	warn  func(string, ...interface{})
	err   func(string, ...interface{})
	debug func(string, ...interface{})
}

// Logger creates a new LoggerStruct instance
func Logger(condition bool) *LoggerStruct {
	return &LoggerStruct{
		info: func(msg string, args ...interface{}) {
			if condition {
				colorable.NewColorableStdout().Write([]byte(fmt.Sprintf("\033[1;34mINFO: "+msg+"\033[0m\n", args...))) // Blue
			}
		},
		warn: func(msg string, args ...interface{}) {
			if condition {
				colorable.NewColorableStdout().Write([]byte(fmt.Sprintf("\033[1;33mWARN: "+msg+"\033[0m\n", args...))) // Yellow
			}
		},
		err: func(msg string, args ...interface{}) {
			if condition {
				colorable.NewColorableStderr().Write([]byte(fmt.Sprintf("\033[1;31mERROR: "+msg+"\033[0m\n", args...))) // Red
			}
		},
		debug: func(msg string, args ...interface{}) {
			if condition {
				colorable.NewColorableStdout().Write([]byte(fmt.Sprintf("\033[1;36mDEBUG: "+msg+"\033[0m\n", args...))) // Cyan
			}
		},
	}
}
