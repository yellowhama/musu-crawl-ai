package agent

import (
	"context"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var globalOrchestrator *Orchestrator

// StartMCPServer starts a Model Context Protocol server on stdio.
func StartMCPServer(version string) error {
	globalOrchestrator = NewOrchestrator("./wiki")

	s := server.NewMCPServer(
		"musu-crawl-ai",
		version,
	)

	// Tool: Fetch
	fetchTool := mcp.NewTool("fetch",
		mcp.WithDescription("Fetch and parse content from YouTube, GitHub, Arxiv, or Web into clean Markdown"),
	)
	s.AddTool(fetchTool, handleFetch)

	// Tool: Search
	searchTool := mcp.NewTool("search",
		mcp.WithDescription("Search the local knowledge base using keywords or semantics"),
	)
	s.AddTool(searchTool, handleSearch)

	// Tool: Research
	researchTool := mcp.NewTool("research",
		mcp.WithDescription("Perform autonomous, recursive research on a topic"),
	)
	s.AddTool(researchTool, handleResearch)

	fmt.Fprintf(os.Stderr, "🚀 musu-crawl-ai MCP Server %s started on stdio\n", version)

	return server.ServeStdio(s)
}

func handleFetch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments.(map[string]interface{})
	source, _ := args["source"].(string)
	url, _ := args["url"].(string)
	project, ok := args["project"].(string)
	if !ok { project = "default" }

	res, err := globalOrchestrator.FetchAction(source, url, project, "en", "llama3")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(res), nil
}

func handleSearch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments.(map[string]interface{})
	query, _ := args["query"].(string)
	project, ok := args["project"].(string)
	if !ok { project = "all" }
	semantic, _ := args["semantic"].(bool)

	entries, err := globalOrchestrator.SearchAction(query, project, semantic, "nomic-embed-text")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	resText := fmt.Sprintf("Found %d results:\n", len(entries))
	for i, e := range entries {
		resText += fmt.Sprintf("%d. [%s] %s (ID: %s, Project: %s)\nSummary: %s\n\n", i+1, e.Source, e.Title, e.ID, e.Project, e.Summary)
	}

	return mcp.NewToolResultText(resText), nil
}

func handleResearch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments.(map[string]interface{})
	question, _ := args["question"].(string)
	project, _ := args["project"].(string)
	depthVal, ok := args["depth"].(float64)
	depth := 2
	if ok { depth = int(depthVal) }

	res, err := globalOrchestrator.ResearchAction(question, project, depth, "llama3")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(res), nil
}
