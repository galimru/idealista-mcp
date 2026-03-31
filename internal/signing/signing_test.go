package signing

import (
	"net/url"
	"os"
	"testing"
)

func TestSortedEncode(t *testing.T) {
	tests := []struct {
		name   string
		params url.Values
		want   string
	}{
		{
			name:   "empty",
			params: url.Values{},
			want:   "",
		},
		{
			name:   "single param",
			params: url.Values{"operation": {"sale"}},
			want:   "operation=sale",
		},
		{
			name: "sorted params",
			params: url.Values{
				"propertyType": {"homes"},
				"operation":    {"sale"},
				"locationId":   {"0-EU-ES-46"},
			},
			want: "locationId=0-EU-ES-46&operation=sale&propertyType=homes",
		},
		{
			name:   "space encoded as %20",
			params: url.Values{"prefix": {"New York"}},
			want:   "prefix=New%20York",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sortedEncode(tt.params)
			if got != tt.want {
				t.Errorf("sortedEncode() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSign(t *testing.T) {
	os.Setenv("IDEALISTA_SIGNING_SECRET", "testSecret")
	s, err := New()
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	seed, sig, err := s.Sign("GET", url.Values{"prefix": {"Valencia"}}, nil)
	if err != nil {
		t.Fatalf("Sign() error: %v", err)
	}
	if seed == "" {
		t.Error("seed should not be empty")
	}
	if sig == "" {
		t.Error("signature should not be empty")
	}
	if len(sig) != 64 {
		t.Errorf("HMAC-SHA256 hex should be 64 chars, got %d", len(sig))
	}
}

func TestNewMissingSecret(t *testing.T) {
	os.Unsetenv("IDEALISTA_SIGNING_SECRET")
	_, err := New()
	if err == nil {
		t.Error("expected error when IDEALISTA_SIGNING_SECRET is missing")
	}
}
