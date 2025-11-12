package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func runPrompt() {
	mainPrompt, err := promptsFS.ReadFile("prompts/prompt.txt")
	if err != nil {
		log.Fatalf("Failed to read main prompt: %v", err)
	}

	finalPrompt := string(mainPrompt)

	if isGitHubCI() {
		log.Println("Detected GitHub Actions environment. Including GitHub-specific prompts...")

		if githubPrompt, err := promptsFS.ReadFile("prompts/github.txt"); err == nil {
			finalPrompt += "\n\n" + string(githubPrompt)
		}

		if prTemplate, err := promptsFS.ReadFile("prompts/pr_template.txt"); err == nil {
			finalPrompt += "\n\n" + string(prTemplate)
		}
	}

	fmt.Print(strings.TrimSpace(finalPrompt) + "\n")
}

func isGitHubCI() bool {
	return os.Getenv("GITHUB_ACTIONS") == "true"
}
