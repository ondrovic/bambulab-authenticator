package utils

import (
	"testing"
)

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected bool
	}{
		{
			name:     "Empty string",
			s:        "",
			expected: true,
		},
		{
			name:     "Non-Empty string",
			s:        "hello",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsEmpty(tt.s)
			if got != tt.expected {
				t.Errorf("IsEmpty() = %v, expected %v", got, tt.expected)
			}
		})
	}
}
