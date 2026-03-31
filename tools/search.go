package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/galimru/idealista-mcp/client"
	"github.com/galimru/idealista-mcp/internal/api"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterSearchTools adds the search_ads tool to the MCP server.
func RegisterSearchTools(s *server.MCPServer, c client.APIClient) {
	s.AddTool(
		mcp.NewTool("search_ads",
			mcp.WithDescription("Search Idealista property listings with optional price, size, and room filters"),
			mcp.WithString("location_id",
				mcp.Required(),
				mcp.Description("Location ID obtained from search_locations (e.g. '0-EU-ES-46-02-002-250')"),
			),
			mcp.WithString("operation",
				mcp.Required(),
				mcp.Description("Type of operation: sale or rent"),
				mcp.Enum("sale", "rent"),
			),
			mcp.WithString("property_type",
				mcp.Required(),
				mcp.Description("Type of property: homes, offices, premises, garages, bedrooms, newDevelopments"),
				mcp.Enum("homes", "offices", "premises", "garages", "bedrooms", "newDevelopments"),
			),
			mcp.WithNumber("num_page",
				mcp.Description("Page number to retrieve (default: 1)"),
			),
			mcp.WithNumber("max_items",
				mcp.Description("Maximum number of results per page (default: 20, max: 50)"),
			),
			mcp.WithString("sort",
				mcp.Description("Sort order: desc (newest first) or asc (oldest first)"),
				mcp.Enum("desc", "asc"),
			),
			mcp.WithNumber("min_price",
				mcp.Description("Minimum price filter"),
			),
			mcp.WithNumber("max_price",
				mcp.Description("Maximum price filter"),
			),
			mcp.WithNumber("min_rooms",
				mcp.Description("Minimum number of rooms"),
			),
			mcp.WithNumber("max_rooms",
				mcp.Description("Maximum number of rooms"),
			),
			mcp.WithNumber("min_size",
				mcp.Description("Minimum size in square metres"),
			),
			mcp.WithNumber("max_size",
				mcp.Description("Maximum size in square metres"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			locationID, err := req.RequireString("location_id")
			if err != nil {
				return nil, err
			}
			operation, err := req.RequireString("operation")
			if err != nil {
				return nil, err
			}
			propertyType, err := req.RequireString("property_type")
			if err != nil {
				return nil, err
			}

			bodyParams := url.Values{
				"locationId":   {locationID},
				"operation":    {operation},
				"propertyType": {propertyType},
			}

			numPage := int(req.GetFloat("num_page", 1))
			bodyParams.Set("numPage", strconv.Itoa(numPage))

			maxItems := int(req.GetFloat("max_items", 20))
			if maxItems > 50 {
				maxItems = 50
			}
			bodyParams.Set("maxItems", strconv.Itoa(maxItems))

			if sort := req.GetString("sort", ""); sort != "" {
				bodyParams.Set("sort", sort)
			}
			if v := req.GetFloat("min_price", 0); v > 0 {
				bodyParams.Set("minPrice", strconv.FormatFloat(v, 'f', 0, 64))
			}
			if v := req.GetFloat("max_price", 0); v > 0 {
				bodyParams.Set("maxPrice", strconv.FormatFloat(v, 'f', 0, 64))
			}
			if v := int(req.GetFloat("min_rooms", 0)); v > 0 {
				bodyParams.Set("minRooms", strconv.Itoa(v))
			}
			if v := int(req.GetFloat("max_rooms", 0)); v > 0 {
				bodyParams.Set("maxRooms", strconv.Itoa(v))
			}
			if v := req.GetFloat("min_size", 0); v > 0 {
				bodyParams.Set("minSize", strconv.FormatFloat(v, 'f', 0, 64))
			}
			if v := req.GetFloat("max_size", 0); v > 0 {
				bodyParams.Set("maxSize", strconv.FormatFloat(v, 'f', 0, 64))
			}

			respBody, err := c.Post(ctx, api.SearchURL, bodyParams)
			if err != nil {
				return nil, fmt.Errorf("search ads: %w", err)
			}

			var resp api.SearchResponse
			if err := json.Unmarshal(respBody, &resp); err != nil {
				return nil, fmt.Errorf("parse search response: %w", err)
			}
			return structJSON(resp)
		},
	)
}
