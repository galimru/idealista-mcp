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
	FormerPrice         float64 `json:"formerPrice"`
	PriceDropValue      int     `json:"priceDropValue"`
	PriceDropPercentage int     `json:"priceDropPercentage"`
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
	HasSwimmingPool    bool `json:"hasSwimmingPool"`
	HasTerrace         bool `json:"hasTerrace"`
	HasAirConditioning bool `json:"hasAirConditioning"`
	HasBoxRoom         bool `json:"hasBoxRoom"`
	HasGarden          bool `json:"hasGarden"`
}

// Property is a single property listing returned by the search endpoint.
type Property struct {
	PropertyCode string    `json:"propertyCode"`
	Floor        string    `json:"floor"`
	Price        float64   `json:"price"`
	PriceInfo    PriceInfo `json:"priceInfo"`
	PropertyType string    `json:"propertyType"`
	Operation    string    `json:"operation"`
	Size         float64   `json:"size"`
	Exterior     bool      `json:"exterior"`
	Rooms        int       `json:"rooms"`
	Bathrooms    int       `json:"bathrooms"`
	Address      string    `json:"address"`
	Municipality string    `json:"municipality"`
	District     string    `json:"district"`
	Neighborhood string    `json:"neighborhood"`
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	Distance     string    `json:"distance"`
	URL          string    `json:"url"`
	HasLift      bool      `json:"hasLift"`
	PriceByArea  float64   `json:"priceByArea"`
	Features     Features  `json:"features"`
	HasVideo     bool      `json:"hasVideo"`
	Has3DTour    bool      `json:"has3DTour"`
	Has360       bool      `json:"has360"`
	ContactInfo  struct {
		CommercialName string `json:"commercialName"`
		Phone1         struct {
			FormattedPhone string `json:"formattedPhone"`
		} `json:"phone1"`
		UserType string `json:"userType"`
	} `json:"contactInfo"`
	ParkingSpace struct {
		HasParkingSpace               bool `json:"hasParkingSpace"`
		IsParkingSpaceIncludedInPrice bool `json:"isParkingSpaceIncludedInPrice"`
	} `json:"parkingSpace"`
}

// SearchResponse is the full response from the search endpoint.
type SearchResponse struct {
	ElementList  []Property `json:"elementList"`
	Total        int        `json:"total"`
	TotalPages   int        `json:"totalPages"`
	ActualPage   int        `json:"actualPage"`
	ItemsPerPage int        `json:"itemsPerPage"`
}

// AdDetail is the response from the ad detail endpoint.
type AdDetail struct {
	AdID      int64   `json:"adid"`
	Price     float64 `json:"price"`
	PriceInfo struct {
		CurrencySuffix string `json:"currencySuffix"`
	} `json:"priceInfo"`
	Operation            string `json:"operation"`
	PropertyType         string `json:"propertyType"`
	ExtendedPropertyType string `json:"extendedPropertyType"`
	State                string `json:"state"`
	Multimedia           struct {
		Images []struct {
			URL string `json:"url"`
			Tag string `json:"tag"`
		} `json:"images"`
		Videos []struct {
			URL string `json:"url"`
		} `json:"videos"`
		Virtual3DTours []struct {
			URL string `json:"url"`
		} `json:"virtual3DTours"`
	} `json:"multimedia"`
	Ubication struct {
		Title                    string  `json:"title"`
		Latitude                 float64 `json:"latitude"`
		Longitude                float64 `json:"longitude"`
		AdministrativeAreaLevel1 string  `json:"administrativeAreaLevel1"`
		AdministrativeAreaLevel2 string  `json:"administrativeAreaLevel2"`
		AdministrativeAreaLevel3 string  `json:"administrativeAreaLevel3"`
		AdministrativeAreaLevel4 string  `json:"administrativeAreaLevel4"`
		HasHiddenAddress         bool    `json:"hasHiddenAddress"`
	} `json:"ubication"`
	ContactInfo struct {
		CommercialName string `json:"commercialName"`
		Phone1         struct {
			FormattedPhoneWithPrefix string `json:"formattedPhoneWithPrefix"`
		} `json:"phone1"`
		ExternalReference string `json:"externalReference"`
		UserType          string `json:"userType"`
		Address           struct {
			StreetName   string `json:"streetName"`
			StreetNumber int    `json:"streetNumber"`
			LocationName string `json:"locationName"`
			PostalCode   string `json:"postalCode"`
		} `json:"address"`
	} `json:"contactInfo"`
	MoreCharacteristics struct {
		CommunityCosts          float64 `json:"communityCosts"`
		RoomNumber              int     `json:"roomNumber"`
		BathNumber              int     `json:"bathNumber"`
		Exterior                bool    `json:"exterior"`
		Floor                   string  `json:"floor"`
		ConstructedArea         float64 `json:"constructedArea"`
		UsableArea              float64 `json:"usableArea"`
		Lift                    bool    `json:"lift"`
		Garden                  bool    `json:"garden"`
		SwimmingPool            bool    `json:"swimmingPool"`
		IsDuplex                bool    `json:"isDuplex"`
		IsPenthouse             bool    `json:"isPenthouse"`
		IsStudio                bool    `json:"isStudio"`
		HousingFurnitures       string  `json:"housingFurnitures"`
		EnergyCertificationType string  `json:"energyCertificationType"`
		Status                  string  `json:"status"`
	} `json:"moreCharacteristics"`
	TranslatedTexts struct {
		LayoutDescription           string `json:"layoutDescription"`
		CharacteristicsDescriptions []struct {
			Title          string `json:"title"`
			DetailFeatures []struct {
				Phrase string `json:"phrase"`
			} `json:"detailFeatures"`
		} `json:"characteristicsDescriptions"`
	} `json:"translatedTexts"`
	Comments []struct {
		PropertyComment string `json:"propertyComment"`
		Language        string `json:"language"`
		DefaultLanguage bool   `json:"defaultLanguage"`
	} `json:"comments"`
	DetailWebLink       string `json:"detailWebLink"`
	AllowsCounterOffers bool   `json:"allowsCounterOffers"`
	AllowsRemoteVisit   bool   `json:"allowsRemoteVisit"`
	ModificationDate    struct {
		Text string `json:"text"`
	} `json:"modificationDate"`
}

// AdStatValue is a single stat with numeric value and display text.
type AdStatValue struct {
	Value int    `json:"value"`
	Text  string `json:"text"`
}

// AdStats is the response from the ad stats endpoint.
type AdStats struct {
	Views        AdStatValue `json:"views"`
	ContactMails AdStatValue `json:"contactMails"`
	SentToFriend AdStatValue `json:"sentToFriend"`
	Favorites    AdStatValue `json:"favorites"`
}
