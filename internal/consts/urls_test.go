package consts

import (
	"testing"
)

func TestRegionalURL(t *testing.T) {
	tests := []struct {
		name     string
		url      URL
		region   string
		expected URL
		wantErr  bool
	}{
		{
			name:     "Empty region",
			url:      "https://example.com",
			region:   "",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "China region",
			url:      "https://example.com",
			region:   "china",
			expected: "https://example.cn",
			wantErr:  false,
		},
		{
			name:     "Non-China region",
			url:      "https://example.com",
			region:   "us",
			expected: "https://example.com",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RegionalURL(tt.url, tt.region)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegionalURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("RegionalURL() = %v, expected %v", got, tt.expected)
			}
		})
	}
}
