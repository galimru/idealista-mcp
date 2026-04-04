package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/galimru/idealista-mcp/client"
	"github.com/galimru/idealista-mcp/internal/api"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type adStatsResult struct {
	Views        int `json:"views"`
	ContactMails int `json:"contactMails"`
	SentToFriend int `json:"sentToFriend"`
	Favorites    int `json:"favorites"`
}

func fetchAdStats(ctx context.Context, apiClient client.APIClient, adID string) (*mcp.CallToolResult, error) {
	rawURL := fmt.Sprintf("%s/%s/stats", api.StatsURL, adID)
	queryParams := url.Values{"language": {"en"}}

	body, err := apiClient.Get(ctx, rawURL, queryParams)
	if err != nil {
		return nil, fmt.Errorf("get ad stats: %w", err)
	}

	var resp api.AdStats
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse ad stats response: %w", err)
	}

	result := adStatsResult{
		Views:        resp.Views.Value,
		ContactMails: resp.ContactMails.Value,
		SentToFriend: resp.SentToFriend.Value,
		Favorites:    resp.Favorites.Value,
	}
	return structJSON(result)
}

func RegisterStatsTools(s *server.MCPServer, runtime *RuntimeProvider) {
	s.AddTool(
		mcp.NewTool("get_ad_stats",
			mcp.WithDescription("Get engagement statistics for a property listing: views, email contacts, sent-to-friend count, and favourites."),
			mcp.WithString("ad_id",
				mcp.Required(),
				mcp.Description("The property ad ID (numeric code, e.g. 111033504)"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			apiClient, err := runtimeAPIClient(runtime)
			if err != nil {
				return nil, err
			}

			adID, err := req.RequireString("ad_id")
			if err != nil {
				return nil, err
			}
			return fetchAdStats(ctx, apiClient, adID)
		},
	)
}
