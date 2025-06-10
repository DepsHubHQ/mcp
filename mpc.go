package main

import (
	"context"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

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

	s.AddTool(listVersionsTool, handleUpdateDetailsTool)

	return s
}
