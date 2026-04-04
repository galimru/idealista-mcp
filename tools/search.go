package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/galimru/idealista-mcp/internal/api"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type searchPriceDropResult struct {
	FormerPrice    float64 `json:"formerPrice"`
	DropValue      int     `json:"dropValue"`
	DropPercentage int     `json:"dropPercentage"`
}

type searchResultItem struct {
	AdID         string                 `json:"adId"`
	PropertyType string                 `json:"propertyType"`
	Operation    string                 `json:"operation"`
	Price        float64                `json:"price"`
	Currency     string                 `json:"currency"`
	PricePerM2   float64                `json:"pricePerM2,omitempty"`
	PriceDrop    *searchPriceDropResult `json:"priceDrop,omitempty"`

	SizeM2    float64 `json:"sizeM2"`
	Rooms     int     `json:"rooms"`
	Bathrooms int     `json:"bathrooms"`
	Floor     string  `json:"floor,omitempty"`
	Exterior  bool    `json:"exterior"`

	Location struct {
		Address       string `json:"address"`
		Neighborhood  string `json:"neighborhood,omitempty"`
		District      string `json:"district,omitempty"`
		Municipality  string `json:"municipality"`
		GoogleMapsURL string `json:"googleMapsUrl"`
		Distance      string `json:"distance,omitempty"`
	} `json:"location"`

	Features struct {
		Lift         bool `json:"lift"`
		AirCon       bool `json:"airConditioning"`
		Terrace      bool `json:"terrace,omitempty"`
		SwimmingPool bool `json:"swimmingPool,omitempty"`
		Garden       bool `json:"garden,omitempty"`
		BoxRoom      bool `json:"boxRoom,omitempty"`
		Parking      *struct {
			IncludedInPrice bool `json:"includedInPrice"`
		} `json:"parking,omitempty"`
	} `json:"features"`

	Contact struct {
		Name     string `json:"name"`
		Phone    string `json:"phone,omitempty"`
		UserType string `json:"userType"`
	} `json:"contact"`

	Media struct {
		HasVideo  bool `json:"hasVideo,omitempty"`
		Has3DTour bool `json:"has3DTour,omitempty"`
		Has360    bool `json:"has360,omitempty"`
	} `json:"media,omitempty"`

	WebLink string `json:"webLink"`
}

type searchResult struct {
	Items      []searchResultItem `json:"items"`
	Total      int                `json:"total"`
	TotalPages int                `json:"totalPages"`
	ActualPage int                `json:"actualPage"`
}

type searchScope struct {
	LocationID  string
	Center      string
	Distance    float64
	HasCenter   bool
	HasDistance bool
}

func mapSearchResult(resp api.SearchResponse) searchResult {
	result := searchResult{
		Total:      resp.Total,
		TotalPages: resp.TotalPages,
		ActualPage: resp.ActualPage,
	}

	for _, p := range resp.ElementList {
		item := searchResultItem{
			AdID:         p.PropertyCode,
			PropertyType: p.PropertyType,
			Operation:    p.Operation,
			Price:        p.Price,
			Currency:     p.PriceInfo.Price.CurrencySuffix,
			PricePerM2:   p.PriceByArea,
			SizeM2:       p.Size,
			Rooms:        p.Rooms,
			Bathrooms:    p.Bathrooms,
			Floor:        p.Floor,
			Exterior:     p.Exterior,
			WebLink:      p.URL,
		}

		if d := p.PriceInfo.Price.PriceDropInfo; d != nil {
			item.PriceDrop = &searchPriceDropResult{
				FormerPrice:    d.FormerPrice,
				DropValue:      d.PriceDropValue,
				DropPercentage: d.PriceDropPercentage,
			}
		}

		item.Location.Address = p.Address
		item.Location.Neighborhood = p.Neighborhood
		item.Location.District = p.District
		item.Location.Municipality = p.Municipality
		item.Location.GoogleMapsURL = fmt.Sprintf(
			"https://www.google.com/maps?q=%f,%f", p.Latitude, p.Longitude,
		)
		item.Location.Distance = p.Distance

		item.Features.Lift = p.HasLift
		item.Features.AirCon = p.Features.HasAirConditioning
		item.Features.Terrace = p.Features.HasTerrace
		item.Features.SwimmingPool = p.Features.HasSwimmingPool
		item.Features.Garden = p.Features.HasGarden
		item.Features.BoxRoom = p.Features.HasBoxRoom
		if p.ParkingSpace.HasParkingSpace {
			item.Features.Parking = &struct {
				IncludedInPrice bool `json:"includedInPrice"`
			}{IncludedInPrice: p.ParkingSpace.IsParkingSpaceIncludedInPrice}
		}

		item.Contact.Name = p.ContactInfo.CommercialName
		item.Contact.Phone = p.ContactInfo.Phone1.FormattedPhone
		item.Contact.UserType = p.ContactInfo.UserType

		item.Media.HasVideo = p.HasVideo
		item.Media.Has3DTour = p.Has3DTour
		item.Media.Has360 = p.Has360

		result.Items = append(result.Items, item)
	}

	return result
}

func buildSearchBodyParams(scope searchScope, req mcp.CallToolRequest) (url.Values, error) {
	hasLocation := scope.LocationID != ""
	hasCenterSearch := scope.HasCenter || scope.HasDistance

	if !hasLocation && !hasCenterSearch {
		return nil, fmt.Errorf("search_ads requires either location_id or center + distance")
	}
	if hasLocation && hasCenterSearch {
		return nil, fmt.Errorf("search_ads accepts either location_id or center + distance, not both")
	}
	if hasCenterSearch && !(scope.HasCenter && scope.HasDistance) {
		return nil, fmt.Errorf("center and distance must be provided together")
	}
	if scope.HasDistance && scope.Distance <= 0 {
		return nil, fmt.Errorf("distance must be greater than zero")
	}
	if scope.HasCenter {
		if err := validateCenter(scope.Center); err != nil {
			return nil, err
		}
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
		"operation":    {operation},
		"propertyType": {propertyType},
		"sort":         {"asc"},
	}
	if hasLocation {
		bodyParams.Set("locationId", scope.LocationID)
	}
	if scope.HasCenter {
		bodyParams.Set("center", scope.Center)
		bodyParams.Set("distance", strconv.FormatFloat(scope.Distance, 'f', 0, 64))
		bodyParams.Set("order", "distance")
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

	return bodyParams, nil
}

func validateCenter(center string) error {
	parts := strings.Split(center, ",")
	if len(parts) != 2 {
		return fmt.Errorf("center must be in 'lat,lng' format")
	}
	lat := strings.TrimSpace(parts[0])
	lng := strings.TrimSpace(parts[1])
	if lat == "" || lng == "" {
		return fmt.Errorf("center must be in 'lat,lng' format")
	}
	if _, err := strconv.ParseFloat(lat, 64); err != nil {
		return fmt.Errorf("center latitude must be a valid number")
	}
	if _, err := strconv.ParseFloat(lng, 64); err != nil {
		return fmt.Errorf("center longitude must be a valid number")
	}
	return nil
}

// RegisterSearchTools adds the search_ads tool to the MCP server.
func RegisterSearchTools(s *server.MCPServer, runtime *RuntimeProvider) {
	s.AddTool(
		mcp.NewTool("search_ads",
			mcp.WithDescription("Search Idealista property listings using either a location_id or a center + distance search, with optional price, size, and room filters"),
			mcp.WithString("location_id",
				mcp.Description("Location ID obtained from search_locations (e.g. '0-EU-ES-46-02-002-250'). Cannot be combined with center + distance."),
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
				mcp.Description("Sort order: asc (default) or desc (newest first)."),
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
			mcp.WithString("center",
				mcp.Description("Center point for distance search in 'lat,lng' format, e.g. '39.474545,-0.343552'. Must be provided together with distance and cannot be combined with location_id."),
			),
			mcp.WithNumber("distance",
				mcp.Description("Distance in metres for center-based search. Must be provided together with center and cannot be combined with location_id."),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			apiClient, err := runtimeAPIClient(runtime)
			if err != nil {
				return nil, err
			}

			args := req.GetArguments()
			scope := searchScope{
				LocationID:  req.GetString("location_id", ""),
				Center:      req.GetString("center", ""),
				Distance:    req.GetFloat("distance", 0),
				HasCenter:   args["center"] != nil,
				HasDistance: args["distance"] != nil,
			}

			bodyParams, err := buildSearchBodyParams(scope, req)
			if err != nil {
				return nil, err
			}

			respBody, err := apiClient.Post(ctx, api.SearchURL, bodyParams)
			if err != nil {
				return nil, fmt.Errorf("search ads: %w", err)
			}

			var resp api.SearchResponse
			if err := json.Unmarshal(respBody, &resp); err != nil {
				return nil, fmt.Errorf("parse search response: %w", err)
			}
			return structJSON(mapSearchResult(resp))
		},
	)
}
