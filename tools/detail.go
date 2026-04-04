package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"

	"github.com/galimru/idealista-mcp/client"
	"github.com/galimru/idealista-mcp/internal/api"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var inmuebleRe = regexp.MustCompile(`/inmueble/(\d+)`)

type adDetailImageResult struct {
	URL  string `json:"url"`
	Room string `json:"room"`
}

type adDetailFeatureGroupResult struct {
	Title    string   `json:"title"`
	Features []string `json:"features"`
}

type adDetailResult struct {
	AdID         int64   `json:"adId"`
	Price        float64 `json:"price"`
	Currency     string  `json:"currency"`
	Operation    string  `json:"operation"`
	PropertyType string  `json:"propertyType"`
	State        string  `json:"state"`

	Location struct {
		Title        string  `json:"title"`
		Latitude     float64 `json:"latitude"`
		Longitude    float64 `json:"longitude"`
		Neighborhood string  `json:"neighborhood,omitempty"`
		District     string  `json:"district,omitempty"`
		City         string  `json:"city,omitempty"`
		Region       string  `json:"region,omitempty"`
	} `json:"location"`

	Contact struct {
		Name              string `json:"name"`
		Phone             string `json:"phone,omitempty"`
		ExternalReference string `json:"externalReference,omitempty"`
		UserType          string `json:"userType"`
		Address           string `json:"address,omitempty"`
	} `json:"contact"`

	Characteristics struct {
		Rooms           int     `json:"rooms"`
		Bathrooms       int     `json:"bathrooms"`
		ConstructedArea float64 `json:"constructedAreaM2"`
		UsableArea      float64 `json:"usableAreaM2,omitempty"`
		Floor           string  `json:"floor,omitempty"`
		Layout          string  `json:"layout,omitempty"`
		Exterior        bool    `json:"exterior"`
		Lift            bool    `json:"lift"`
		Garden          bool    `json:"garden"`
		SwimmingPool    bool    `json:"swimmingPool"`
		IsDuplex        bool    `json:"isDuplex"`
		IsPenthouse     bool    `json:"isPenthouse"`
		IsStudio        bool    `json:"isStudio"`
		Furnishings     string  `json:"furnishings,omitempty"`
		Condition       string  `json:"condition,omitempty"`
		CommunityCosts  float64 `json:"communityCostsPerMonth,omitempty"`
		EnergyRating    string  `json:"energyRating,omitempty"`
	} `json:"characteristics"`

	Features    []adDetailFeatureGroupResult `json:"features"`
	Description string                       `json:"description,omitempty"`
	WebLink     string                       `json:"webLink"`
	Images      []adDetailImageResult        `json:"images"`
	VideoCount  int                          `json:"videoCount,omitempty"`
	TourURL     string                       `json:"tourUrl,omitempty"`

	AllowsCounterOffer bool   `json:"allowsCounterOffer"`
	AllowsRemoteVisit  bool   `json:"allowsRemoteVisit"`
	LastUpdated        string `json:"lastUpdated,omitempty"`
}

func mapAdDetail(r api.AdDetail) adDetailResult {
	d := adDetailResult{
		AdID:               r.AdID,
		Price:              r.Price,
		Currency:           r.PriceInfo.CurrencySuffix,
		Operation:          r.Operation,
		PropertyType:       r.ExtendedPropertyType,
		State:              r.State,
		WebLink:            r.DetailWebLink,
		AllowsCounterOffer: r.AllowsCounterOffers,
		AllowsRemoteVisit:  r.AllowsRemoteVisit,
		LastUpdated:        r.ModificationDate.Text,
	}

	d.Location.Title = r.Ubication.Title
	d.Location.Latitude = r.Ubication.Latitude
	d.Location.Longitude = r.Ubication.Longitude
	d.Location.Neighborhood = r.Ubication.AdministrativeAreaLevel4
	d.Location.District = r.Ubication.AdministrativeAreaLevel3
	d.Location.City = r.Ubication.AdministrativeAreaLevel2
	d.Location.Region = r.Ubication.AdministrativeAreaLevel1

	d.Contact.Name = r.ContactInfo.CommercialName
	d.Contact.Phone = r.ContactInfo.Phone1.FormattedPhoneWithPrefix
	d.Contact.ExternalReference = r.ContactInfo.ExternalReference
	d.Contact.UserType = r.ContactInfo.UserType
	if !r.Ubication.HasHiddenAddress && r.ContactInfo.Address.StreetName != "" {
		a := r.ContactInfo.Address
		if a.StreetNumber > 0 {
			d.Contact.Address = fmt.Sprintf("%s %d, %s %s", a.StreetName, a.StreetNumber, a.PostalCode, a.LocationName)
		} else {
			d.Contact.Address = fmt.Sprintf("%s, %s %s", a.StreetName, a.PostalCode, a.LocationName)
		}
	}

	d.Characteristics.Rooms = r.MoreCharacteristics.RoomNumber
	d.Characteristics.Bathrooms = r.MoreCharacteristics.BathNumber
	d.Characteristics.ConstructedArea = r.MoreCharacteristics.ConstructedArea
	d.Characteristics.UsableArea = r.MoreCharacteristics.UsableArea
	d.Characteristics.Floor = r.MoreCharacteristics.Floor
	d.Characteristics.Layout = r.TranslatedTexts.LayoutDescription
	d.Characteristics.Exterior = r.MoreCharacteristics.Exterior
	d.Characteristics.Lift = r.MoreCharacteristics.Lift
	d.Characteristics.Garden = r.MoreCharacteristics.Garden
	d.Characteristics.SwimmingPool = r.MoreCharacteristics.SwimmingPool
	d.Characteristics.IsDuplex = r.MoreCharacteristics.IsDuplex
	d.Characteristics.IsPenthouse = r.MoreCharacteristics.IsPenthouse
	d.Characteristics.IsStudio = r.MoreCharacteristics.IsStudio
	d.Characteristics.Furnishings = r.MoreCharacteristics.HousingFurnitures
	d.Characteristics.Condition = r.MoreCharacteristics.Status
	d.Characteristics.CommunityCosts = r.MoreCharacteristics.CommunityCosts
	d.Characteristics.EnergyRating = r.MoreCharacteristics.EnergyCertificationType

	// Description: prefer English, fallback to default language
	for _, c := range r.Comments {
		if c.Language == "en" {
			d.Description = c.PropertyComment
			break
		}
	}
	if d.Description == "" {
		for _, c := range r.Comments {
			if c.DefaultLanguage {
				d.Description = c.PropertyComment
				break
			}
		}
	}

	// Features from translatedTexts
	for _, g := range r.TranslatedTexts.CharacteristicsDescriptions {
		group := adDetailFeatureGroupResult{Title: g.Title}
		for _, f := range g.DetailFeatures {
			group.Features = append(group.Features, f.Phrase)
		}
		d.Features = append(d.Features, group)
	}

	// Images
	for _, img := range r.Multimedia.Images {
		d.Images = append(d.Images, adDetailImageResult{URL: img.URL, Room: img.Tag})
	}

	d.VideoCount = len(r.Multimedia.Videos)

	if len(r.Multimedia.Virtual3DTours) > 0 {
		d.TourURL = r.Multimedia.Virtual3DTours[0].URL
	}

	return d
}

func fetchAdDetail(ctx context.Context, apiClient client.APIClient, adID string) (*mcp.CallToolResult, error) {
	rawURL := fmt.Sprintf("%s/%s", api.DetailURL, adID)
	queryParams := url.Values{"language": {"en"}}

	body, err := apiClient.Get(ctx, rawURL, queryParams)
	if err != nil {
		return nil, fmt.Errorf("get ad details: %w", err)
	}

	var resp api.AdDetail
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse ad detail response: %w", err)
	}

	return structJSON(mapAdDetail(resp))
}

func RegisterDetailTools(s *server.MCPServer, runtime *RuntimeProvider) {
	s.AddTool(
		mcp.NewTool("get_ad_details",
			mcp.WithDescription("Get full details for a property listing by its ad ID. Use get_ad_details_by_url when you have a shareable property URL instead."),
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
			return fetchAdDetail(ctx, apiClient, adID)
		},
	)

	s.AddTool(
		mcp.NewTool("get_ad_details_by_url",
			mcp.WithDescription("Get full details for a property listing from its shareable URL, e.g. https://www.idealista.com/inmueble/111033504/"),
			mcp.WithString("url",
				mcp.Required(),
				mcp.Description("The shareable property URL from idealista.com"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			apiClient, err := runtimeAPIClient(runtime)
			if err != nil {
				return nil, err
			}

			propertyURL, err := req.RequireString("url")
			if err != nil {
				return nil, err
			}

			matches := inmuebleRe.FindStringSubmatch(propertyURL)
			if len(matches) < 2 {
				return nil, fmt.Errorf("could not extract ad ID from URL: %s", propertyURL)
			}

			return fetchAdDetail(ctx, apiClient, matches[1])
		},
	)
}
