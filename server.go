package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"
)

func runServer() {
	log.Printf("Starting MCP server...")

	port = os.Getenv("PORT")
	transport = os.Getenv("TRANSPORT")
	baseURL = os.Getenv("BASE_URL")

	log.Printf("Using base URL: %s", baseURL)
	log.Printf("Using transport: %s", transport)

	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	flag.Parse()

	s := NewMCPServer()

	if transport == "http" {
		httpServer := server.NewStreamableHTTPServer(s)
		log.Printf("HTTP server listening on :%s/mcp", port)
		if err := httpServer.Start(fmt.Sprintf(":%s", port)); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	} else {
		log.Printf("Stdio server is starting...")
		if err := server.ServeStdio(s); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}
}
