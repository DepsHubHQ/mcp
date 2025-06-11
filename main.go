package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"
)

var port string
var transport string
var baseURL string

func main() {
	log.Printf("Starting MCP server...")

	flag.StringVar(&transport, "t", "stdio", "Transport type (stdio or http)")
	flag.StringVar(&transport, "transport", "stdio", "Transport type (stdio or http)")

	flag.StringVar(&port, "p", "8080", "Port for HTTP transport")
	flag.StringVar(&port, "port", "8080", "Port for HTTP transport")

	baseURL := os.Getenv("BASE_URL")

	log.Printf("Using base URL: %s", baseURL)
	log.Printf("Using transport: %s", transport)
	log.Printf("Using port: %s", port)

	if baseURL == "" {
		baseURL = "https://mcp-api.depshub.com"
	}

	flag.Parse()

	s := NewMCPServer()

	// Only check for "http" since stdio is the default
	if transport == "http" {
		httpServer := server.NewStreamableHTTPServer(s)
		log.Printf("HTTP server listening on :%s/mcp", port)
		if err := httpServer.Start(fmt.Sprintf(":%s", port)); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	} else {
		if err := server.ServeStdio(s); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}
}
