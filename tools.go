package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/mark3labs/mcp-go/mcp"
)

func handleUpdateDetailsTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	name := args["name"].(string)
	ecosystem := args["ecosystem"].(string)
	version := args["version"].(string)
	// Call the actual endpoint
	baseURL := "http://localhost:8080/update-details"
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
	type Response struct {
		Changelog_diff          string   `json:"changelog_diff"`
		Current_vulnerabilities []string `json:"current_vulnerabilities"`
		Recommended_version     string   `json:"recommended_version"`
		Newer_versions          []string `json:"newer_versions"`
	}

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("invalid JSON from backend: %w", err)
	}

	var responseTemplate = `
		Here is a summary of the update details:

		Newer versions: %s
		Recommended new version: %s
		Current vulnerabilities: %s
		Current changelog diff: %s
		`
	var responseMessage = fmt.Sprintf(responseTemplate, result.Newer_versions, result.Recommended_version, result.Current_vulnerabilities, result.Changelog_diff)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: responseMessage,
			},
		},
	}, nil
}
