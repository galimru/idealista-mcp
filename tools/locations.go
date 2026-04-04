package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/galimru/idealista-mcp/internal/api"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterLocationTools adds the search_locations tool to the MCP server.
func RegisterLocationTools(s *server.MCPServer, runtime *RuntimeProvider) {
	s.AddTool(
		mcp.NewTool("search_locations",
			mcp.WithDescription("Search Idealista locations by name prefix to obtain a locationId for property searches"),
			mcp.WithString("prefix",
				mcp.Required(),
				mcp.Description("Location name prefix to search for (e.g. 'Valencia', 'Madrid')"),
			),
			mcp.WithString("property_type",
				mcp.Required(),
				mcp.Description("Type of property: homes, offices, premises, garages, bedrooms, newDevelopments"),
				mcp.Enum("homes", "offices", "premises", "garages", "bedrooms", "newDevelopments"),
			),
			mcp.WithString("operation",
				mcp.Required(),
				mcp.Description("Type of operation: sale or rent"),
				mcp.Enum("sale", "rent"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			apiClient, err := runtimeAPIClient(runtime)
			if err != nil {
				return nil, err
			}

			prefix, err := req.RequireString("prefix")
			if err != nil {
				return nil, err
			}
			propertyType, err := req.RequireString("property_type")
			if err != nil {
				return nil, err
			}
			operation, err := req.RequireString("operation")
			if err != nil {
				return nil, err
			}

			queryParams := url.Values{
				"prefix":       {prefix},
				"propertyType": {propertyType},
				"operation":    {operation},
			}

			body, err := apiClient.Get(ctx, api.LocationsURL, queryParams)
			if err != nil {
				return nil, fmt.Errorf("search locations: %w", err)
			}

			var resp api.LocationsResponse
			if err := json.Unmarshal(body, &resp); err != nil {
				return nil, fmt.Errorf("parse locations response: %w", err)
			}
			return structJSON(resp)
		},
	)
}
