package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	var transport string
	flag.StringVar(&transport, "t", "stdio", "Transport type (stdio or http)")
	flag.StringVar(&transport, "transport", "stdio", "Transport type (stdio or http)")
	flag.Parse()

	s := NewMCPServer()

	// Only check for "http" since stdio is the default
	if transport == "http" {
		httpServer := server.NewStreamableHTTPServer(s)
		log.Printf("HTTP server listening on :8080/mcp")
		if err := httpServer.Start(":8080"); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	} else {
		if err := server.ServeStdio(s); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}
}

func NewMCPServer() *server.MCPServer {
	hooks := &server.Hooks{}
	hooks.AddBeforeAny(func(ctx context.Context, id any, method mcp.MCPMethod, message any) {
		log.Printf("beforeAny: %s, %v, %v\n", method, id, message)
	})
	hooks.AddOnSuccess(func(ctx context.Context, id any, method mcp.MCPMethod, message any, result any) {
		log.Printf("onSuccess: %s, %v, %v, %v\n", method, id, message, result)
	})
	hooks.AddOnError(func(ctx context.Context, id any, method mcp.MCPMethod, message any, err error) {
		log.Printf("onError: %s, %v, %v, %v\n", method, id, message, err)
	})
	hooks.AddBeforeInitialize(func(ctx context.Context, id any, message *mcp.InitializeRequest) {
		log.Printf("beforeInitialize: %v, %v\n", id, message)
	})
	hooks.AddOnRequestInitialization(func(ctx context.Context, id any, message any) error {
		log.Printf("AddOnRequestInitialization: %v, %v\n", id, message)
		// authorization verification and other preprocessing tasks are performed.
		return nil
	})
	hooks.AddAfterInitialize(func(ctx context.Context, id any, message *mcp.InitializeRequest, result *mcp.InitializeResult) {
		log.Printf("afterInitialize: %v, %v, %v\n", id, message, result)
	})
	hooks.AddAfterCallTool(func(ctx context.Context, id any, message *mcp.CallToolRequest, result *mcp.CallToolResult) {
		log.Printf("afterCallTool: %v, %v, %v\n", id, message, result)
	})
	hooks.AddBeforeCallTool(func(ctx context.Context, id any, message *mcp.CallToolRequest) {
		log.Printf("beforeCallTool: %v, %v\n", id, message)
	})
	s := server.NewMCPServer(
		"DepsHub",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	listVersionsTool := mcp.NewTool("update_details",
		mcp.WithDescription("Returns the information about the potential update of a library or package. You have to include the name, ecosystem (eg NPM), and the current version of the library or package to get information about."),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Title("Name"),
			mcp.Description("Name of the library or package to get information about"),
		),
		mcp.WithString("ecosystem",
			mcp.Required(),
			mcp.Title("Ecosystem"),
			mcp.Description("Name of the ecosystem. Can be one of NPM, GO, RUBYGEMS, CARGO, PYPI"),
		),
		mcp.WithString("version",
			mcp.Required(),
			mcp.Title("Version"),
			mcp.Description("Current version of the library or package to get information about. You have to take the version from the project dependencies manifest file (eg package.json)"),
		),
	)

	// Add tool handler
	s.AddTool(listVersionsTool, handleUpdateDetailsTool)

	return s
}

func handleUpdateDetailsTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	name := args["name"].(string)
	ecosystem := args["ecosystem"].(string)
	version := args["version"].(string)
	// Call the actual endpoint
	baseURL := "http://localhost:8080/update-details" // or wherever your Gin server is running
	params := url.Values{}
	params.Set("name", name)
	params.Set("ecosystem", ecosystem)
	params.Set("version", version)
	resp, err := http.Get(fmt.Sprintf("%s?%s", baseURL, params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to call backend: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("backend error: %s", body)
	}
	// Parse and format the response
	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("invalid JSON from backend: %w", err)
	}
	// Optionally convert result to a readable string
	resultBytes, _ := json.MarshalIndent(result, "", "  ")
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(resultBytes),
			},
		},
	}, nil
}
