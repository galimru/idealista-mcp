package api

// TokenResponse is the response from the OAuth token endpoint.
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

// Location is a single result from the locations search endpoint.
type Location struct {
	Name                string `json:"name"`
	LocationID          string `json:"locationId"`
	Divisible           bool   `json:"divisible"`
	Type                string `json:"type"`
	SuggestedLocationID int    `json:"suggestedLocationId"`
	SubTypeText         string `json:"subTypeText"`
	Total               int    `json:"total"`
}

// LocationsResponse is the full response from the locations search endpoint.
type LocationsResponse struct {
	Locations []Location `json:"locations"`
	Total     int        `json:"total"`
}

// PriceDropInfo holds price drop details for a property.
type PriceDropInfo struct {
	FormerPrice        float64 `json:"formerPrice"`
	PriceDropValue     int     `json:"priceDropValue"`
	PriceDropPercentage int    `json:"priceDropPercentage"`
}

// PriceDetail holds the price amount and currency.
type PriceDetail struct {
	Amount         float64        `json:"amount"`
	CurrencySuffix string         `json:"currencySuffix"`
	PriceDropInfo  *PriceDropInfo `json:"priceDropInfo,omitempty"`
}

// PriceInfo wraps the price detail.
type PriceInfo struct {
	Price PriceDetail `json:"price"`
}

// Features holds boolean amenity flags for a property.
type Features struct {
	HasSwimmingPool   bool `json:"hasSwimmingPool"`
	HasTerrace        bool `json:"hasTerrace"`
	HasAirConditioning bool `json:"hasAirConditioning"`
	HasBoxRoom        bool `json:"hasBoxRoom"`
	HasGarden         bool `json:"hasGarden"`
}

// Property is a single property listing returned by the search endpoint.
type Property struct {
	PropertyCode    string    `json:"propertyCode"`
	Thumbnail       string    `json:"thumbnail"`
	NumPhotos       int       `json:"numPhotos"`
	Floor           string    `json:"floor"`
	Price           float64   `json:"price"`
	PriceInfo       PriceInfo `json:"priceInfo"`
	PropertyType    string    `json:"propertyType"`
	Operation       string    `json:"operation"`
	Size            float64   `json:"size"`
	Exterior        bool      `json:"exterior"`
	Rooms           int       `json:"rooms"`
	Bathrooms       int       `json:"bathrooms"`
	Address         string    `json:"address"`
	Province        string    `json:"province"`
	Municipality    string    `json:"municipality"`
	District        string    `json:"district"`
	Neighborhood    string    `json:"neighborhood"`
	LocationID      string    `json:"locationId"`
	Latitude        float64   `json:"latitude"`
	Longitude       float64   `json:"longitude"`
	URL             string    `json:"url"`
	Description     string    `json:"description"`
	Status          string    `json:"status"`
	NewDevelopment  bool      `json:"newDevelopment"`
	HasLift         bool      `json:"hasLift"`
	PriceByArea     float64   `json:"priceByArea"`
	Features        Features  `json:"features"`
}

// SearchResponse is the full response from the search endpoint.
type SearchResponse struct {
	ElementList  []Property `json:"elementList"`
	Total        int        `json:"total"`
	TotalPages   int        `json:"totalPages"`
	ActualPage   int        `json:"actualPage"`
	ItemsPerPage int        `json:"itemsPerPage"`
}
