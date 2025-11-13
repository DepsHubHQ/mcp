package main

import (
	"context"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Set during build using -ldflags "-X main.version=x.y.z"
var version string

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
		version,
		server.WithToolCapabilities(true),
	)

	analyzeDependenciesTool := mcp.NewTool("analyze_dependencies",
		mcp.WithDescription("Returns the information about all the available versions, security vulnerabilities, and release dates for the given libraries or packages."),
		mcp.WithArray("packages",
			mcp.Description("List of package PURLs to analyze"),
			mcp.Items(map[string]any{"purl": "string"}),
		),
	)

	getUpdateInsights := mcp.NewTool("get_update_insights",
		mcp.WithDescription("Returns update insights about specific library versions, including their changelogs, release dates, and migration guides. Use this to get information about a particular versions once you have the list of available versions from the analyze_dependencies tool."),
		mcp.WithObject("package",
			mcp.Description("Current version PURL and the update version PURL (including specific versions) to analyze"),
			mcp.Properties(map[string]any{
				"current_purl": map[string]any{
					"type":        "string",
					"description": "Current version PURL",
				},
				"update_purl": map[string]any{
					"type":        "string",
					"description": "Update version PURL",
				},
			}),
		),
	)

	s.AddTool(analyzeDependenciesTool, handleAnalyzeDependenciesTool)
	s.AddTool(getUpdateInsights, handleGetUpdateInsights)

	return s
}
