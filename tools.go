package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

func handleAnalyzeDependenciesTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	log.Printf("handleAnalyzeDependenciesTool called with args: %v", args)

	// Extract the array of PURLs
	rawPkgs, ok := args["packages"].([]any)
	log.Printf("Extracted raw packages: %v", rawPkgs)
	if !ok {
		return nil, fmt.Errorf("invalid input: expected 'packages' array")
	}

	type PackageInput struct {
		Purl string `json:"purl"`
	}

	var packages []PackageInput
	for _, p := range rawPkgs {
		if s, ok := p.(string); ok {
			packages = append(packages, PackageInput{Purl: s})
		}
	}
	log.Printf("Parsed packages: %v", packages)

	if len(packages) == 0 {
		return nil, fmt.Errorf("no packages provided")
	}

	// Prepare request body
	payload := map[string]any{
		"packages": packages,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request body: %w", err)
	}
	log.Printf("Request payload: %s", string(bodyBytes))

	endpoint := fmt.Sprintf("%s/analyze-dependencies", baseURL)
	log.Printf("Calling backend endpoint: %s", endpoint)
	resp, err := http.Post(endpoint, "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to call backend: %w", err)
	}
	log.Printf("Backend response status: %s", resp.Status)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	log.Printf("Backend response body: %s", string(respBody))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("backend error: %s", respBody)
	}

	// Define new backend response type
	type Vulnerability struct {
		ID string `json:"id"`
	}

	type PackageAnalysisData struct {
		Purl              string          `json:"purl"`
		AvailableVersions []string        `json:"available_versions"`
		Vulnerabilities   []Vulnerability `json:"vulnerabilities"`
		ReleaseDate       string          `json:"release_date"`
	}

	type AnalyzeDependenciesResponse struct {
		Results []PackageAnalysisData `json:"results"`
	}

	var results AnalyzeDependenciesResponse
	if err := json.Unmarshal(respBody, &results); err != nil {
		return nil, fmt.Errorf("invalid JSON from backend: %w", err)
	}

	// Format results for each package
	var summaries []string
	for _, r := range results.Results {
		var vulnIDs []string
		for _, v := range r.Vulnerabilities {
			vulnIDs = append(vulnIDs, v.ID)
		}
		vulnSummary := "None"
		if len(vulnIDs) > 0 {
			vulnSummary = strings.Join(vulnIDs, ", ")
		}

		summary := fmt.Sprintf(`
*%s*:
- Release date: %s
- Newer available versions: %s
- Vulnerabilities: %s`,
			r.Purl,
			r.ReleaseDate,
			strings.Join(r.AvailableVersions, ", "),
			vulnSummary,
		)
		summaries = append(summaries, strings.TrimSpace(summary))
	}

	finalMessage := "Here is the analysis of the requested packages:\n\n" + strings.Join(summaries, "\n\n")

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: finalMessage,
			},
		},
	}, nil
}

func handleGetUpdateInsights(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	type PackageInput struct {
		CurrentPurl string `json:"current_purl"`
		UpdatePurl  string `json:"update_purl"`
	}

	rawPkg, ok := args["package"]
	if !ok {
		return nil, fmt.Errorf("invalid input: expected 'package' object")
	}

	var p PackageInput
	b, err := json.Marshal(rawPkg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal package argument: %w", err)
	}
	if err := json.Unmarshal(b, &p); err != nil {
		return nil, fmt.Errorf("invalid package argument: %w", err)
	}

	// Prepare request body
	payload := map[string]any{
		"package": p,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request body: %w", err)
	}

	endpoint := fmt.Sprintf("%s/get-update-insights", baseURL)
	resp, err := http.Post(endpoint, "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to call backend: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("backend error: %s", respBody)
	}

	// Parse backend response
	type PackageUpdateInsight struct {
		CurrentPurl  string `json:"current_purl"`
		UpdatePurl   string `json:"update_purl"`
		ReleaseNotes string `json:"release_notes,omitempty"`
	}

	type GetUpdateInsightsResponse struct {
		Errors  []string             `json:"errors"`
		Results PackageUpdateInsight `json:"results"`
	}

	var result GetUpdateInsightsResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("invalid JSON from backend: %w", err)
	}

	// Format results for human-readable output

	summary := fmt.Sprintf(`
Update from *%s → %s*

Versions changelog between current and target update:
%s
`, result.Results.CurrentPurl, result.Results.UpdatePurl, result.Results.ReleaseNotes)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: summary,
			},
		},
	}, nil
}
