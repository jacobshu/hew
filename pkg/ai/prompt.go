package ai

import (
	"flag"
	"fmt"
	"log"
	"strings"
)

func combineElements(prompt, task, context string) string {
	var combined strings.Builder

	combined.WriteString("<Prompt>\n")
	combined.WriteString(prompt)
	combined.WriteString("\n</Prompt>\n\n")

	combined.WriteString("<Task>\n")
	combined.WriteString(task)
	combined.WriteString("\n</Task>\n\n")

	combined.WriteString("<Context>\n")
	combined.WriteString(context)
	combined.WriteString("\n</Context>")

	return combined.String()
}

func NewPrompt() {
	// Parse command-line arguments
	promptFile := flag.String("prompt", "", "Path to a .prompt file")
	// chatID := flag.String("chat", "", "ID of an existing chat to continue")
	task := flag.String("task", "", "Task to be added to the prompt")
	apiKey := flag.String("api-key", "", "API key for Claude API")
	flag.Parse()

	if *apiKey == "" {
		log.Fatal("Please provide an API key using the -api-key flag")
	}

	// Get the list of files/directories from remaining arguments
	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("Please provide at least one file or directory")
	}

	// Initialize file handler
	fh := NewFileHandler()

	// Process files and directories
	context, err := fh.ProcessPaths(args)
	if err != nil {
		log.Fatalf("Error processing files: %v", err)
	}

	// Load prompt if specified
	var prompt string
	if *promptFile != "" {
		prompt, err = fh.LoadPrompt(*promptFile)
		if err != nil {
			log.Fatalf("Error loading prompt: %v", err)
		}
	}

	// Combine prompt, task, and context
	combinedMessage := combineElements(prompt, *task, context)

	// Initialize Claude API client
	client := NewAPIClient(*apiKey)

	// Send request to Claude API and handle response
	response, err := client.SendMessage(combinedMessage)
	if err != nil {
		log.Fatalf("Error communicating with Claude API: %v", err)
	}

	fmt.Println("Claude's response:")
	fmt.Println(response)

	// TODO: Implement chat management
}
