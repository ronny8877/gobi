package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/common-nighthawk/go-figure"
)

func welcomeMessage() {
	asciiArt := figure.NewFigure("Welcome to Gobi!", "", true).String()
	var myCuteBorder = lipgloss.Border{
		Top:         "🥬",
		Bottom:      "🐪",
		Left:        "🌵",
		Right:       "🌵",
		TopLeft:     "❤️",
		TopRight:    "❤️",
		BottomLeft:  "❣️",
		BottomRight: "❣️",
	}
	var style = lipgloss.NewStyle().
		Bold(true).
		BorderStyle(myCuteBorder).
		BorderForeground(lipgloss.Color("228")).
		Padding(2)

	// Print the styled ASCII art
	fmt.Println(style.Render(asciiArt))
}

func showFilesList() {
	var style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	fmt.Println(style.Render("Files in the current directory:"))
	list, _ := getFilesList(".")
	fmt.Println(list)
}
