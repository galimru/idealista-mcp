package tools

import (
	"encoding/json"
	"fmt"

	"github.com/galimru/idealista-mcp/client"
	"github.com/mark3labs/mcp-go/mcp"
)

// prettyJSON validates body as JSON and returns it formatted with two-space indentation.
func prettyJSON(body []byte) (*mcp.CallToolResult, error) {
	var raw json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("invalid JSON from API: %w", err)
	}
	pretty, _ := json.MarshalIndent(raw, "", "  ")
	return mcp.NewToolResultText(string(pretty)), nil
}

// structJSON marshals v and returns it as a tool result.
func structJSON(v any) (*mcp.CallToolResult, error) {
	out, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("marshal result: %w", err)
	}
	return mcp.NewToolResultText(string(out)), nil
}

func runtimeAPIClient(runtime *RuntimeProvider) (client.APIClient, error) {
	apiClient, err := runtime.APIClient()
	if err != nil {
		return nil, fmt.Errorf("initialize runtime: %w", err)
	}
	return apiClient, nil
}
