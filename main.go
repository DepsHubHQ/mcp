package main

import (
	"embed"
	"fmt"
	"os"
)

//go:embed prompts/*
var promptsFS embed.FS

var port string
var transport string
var baseURL string

func main() {
	if len(os.Args) < 2 {
		runServer() // default mode
		return
	}

	cmd := os.Args[1]
	switch cmd {
	case "prompt":
		runPrompt()
	case "start":
		runServer()
	default:
		fmt.Println("Unknown command:", cmd)
		fmt.Println("Usage: depshub-mcp [start|prompt]")
		os.Exit(1)
	}
}
