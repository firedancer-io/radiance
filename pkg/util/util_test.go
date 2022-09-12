package util

import (
	"testing"
)

func Test_isValidHostname(t *testing.T) {
	tests := []struct {
		hostname string
		want     bool
	}{
		{
			// Starts with an number.
			hostname: "1solana.com",
			want:     true,
		},
		{
			// Ends with an number.
			hostname: "solana.com1",
			want:     true,
		},
		{
			// Starts with an underscore.
			hostname: "_solana.com",
			want:     false,
		},
		{
			// Ends with an underscore.
			hostname: "solana.com_",
			want:     false,
		},
		{
			// No TLD.
			hostname: "solana",
			want:     false,
		},
		{
			// With TLD.
			hostname: "solana.com",
			want:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.hostname, func(t *testing.T) {
			if got := IsValidHostname(tt.hostname); got != tt.want {
				t.Errorf("isValidHostname(%q) = %v, want %v", tt.hostname, got, tt.want)
			}
		})
	}
}
