package tools

import (
	"net/url"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestBuildSearchBodyParamsLocationOnly(t *testing.T) {
	req := toolRequest(map[string]any{
		"operation":     "sale",
		"property_type": "homes",
	})

	params, err := buildSearchBodyParams(searchScope{LocationID: "0-EU-ES-46"}, req)
	if err != nil {
		t.Fatalf("buildSearchBodyParams() error = %v", err)
	}

	assertParamEquals(t, params, "locationId", "0-EU-ES-46")
	assertParamEquals(t, params, "operation", "sale")
	assertParamEquals(t, params, "propertyType", "homes")
	assertParamEquals(t, params, "numPage", "1")
	assertParamEquals(t, params, "maxItems", "20")
	assertParamEquals(t, params, "sort", "asc")
	assertParamEquals(t, params, "order", "")
}

func TestBuildSearchBodyParamsCenterOnly(t *testing.T) {
	req := toolRequest(map[string]any{
		"operation":     "sale",
		"property_type": "homes",
		"center":        "39.474545,-0.343552",
		"distance":      500.0,
	})

	params, err := buildSearchBodyParams(searchScope{
		Center:      "39.474545,-0.343552",
		Distance:    500,
		HasCenter:   true,
		HasDistance: true,
	}, req)
	if err != nil {
		t.Fatalf("buildSearchBodyParams() error = %v", err)
	}

	if got := params.Get("locationId"); got != "" {
		t.Fatalf("locationId = %q, want empty", got)
	}
	assertParamEquals(t, params, "center", "39.474545,-0.343552")
	assertParamEquals(t, params, "distance", "500")
	assertParamEquals(t, params, "sort", "asc")
	assertParamEquals(t, params, "order", "distance")
}

func TestBuildSearchBodyParamsRejectsMixedInput(t *testing.T) {
	req := toolRequest(map[string]any{
		"operation":     "rent",
		"property_type": "homes",
		"center":        "39.474545,-0.343552",
		"distance":      500.0,
	})

	_, err := buildSearchBodyParams(searchScope{
		LocationID:  "0-EU-ES-46",
		Center:      "39.474545,-0.343552",
		Distance:    500,
		HasCenter:   true,
		HasDistance: true,
	}, req)
	if err == nil {
		t.Fatal("expected error for mixed location and geo input")
	}
}

func TestBuildSearchBodyParamsRejectsPartialCenterSearch(t *testing.T) {
	req := toolRequest(map[string]any{
		"operation":     "sale",
		"property_type": "homes",
		"center":        "39.474545,-0.343552",
	})

	_, err := buildSearchBodyParams(searchScope{
		Center:    "39.474545,-0.343552",
		HasCenter: true,
	}, req)
	if err == nil {
		t.Fatal("expected error for partial center search input")
	}
}

func TestBuildSearchBodyParamsRejectsInvalidDistance(t *testing.T) {
	req := toolRequest(map[string]any{
		"operation":     "sale",
		"property_type": "homes",
		"center":        "39.474545,-0.343552",
		"distance":      0.0,
	})

	_, err := buildSearchBodyParams(searchScope{
		Center:      "39.474545,-0.343552",
		Distance:    0,
		HasCenter:   true,
		HasDistance: true,
	}, req)
	if err == nil {
		t.Fatal("expected error for non-positive distance")
	}
}

func TestBuildSearchBodyParamsRejectsInvalidCenter(t *testing.T) {
	req := toolRequest(map[string]any{
		"operation":     "sale",
		"property_type": "homes",
		"center":        "39.474545",
		"distance":      500.0,
	})

	_, err := buildSearchBodyParams(searchScope{
		Center:      "39.474545",
		Distance:    500,
		HasCenter:   true,
		HasDistance: true,
	}, req)
	if err == nil {
		t.Fatal("expected error for invalid center")
	}
}

func TestBuildSearchBodyParamsRejectsMissingScope(t *testing.T) {
	req := toolRequest(map[string]any{
		"operation":     "sale",
		"property_type": "homes",
	})

	_, err := buildSearchBodyParams(searchScope{}, req)
	if err == nil {
		t.Fatal("expected error for missing search scope")
	}
}

func toolRequest(args map[string]any) mcp.CallToolRequest {
	req := mcp.CallToolRequest{}
	req.Params.Name = "search_ads"
	req.Params.Arguments = args
	return req
}

func assertParamEquals(t *testing.T, params url.Values, key, want string) {
	t.Helper()
	if got := params.Get(key); got != want {
		t.Fatalf("%s = %q, want %q", key, got, want)
	}
}
