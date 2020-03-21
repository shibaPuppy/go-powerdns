package lib

import (
	"fmt"
	"testing"
)

func TestTrimDomain(t *testing.T) {
	testCases := []struct {
		domain     string
		wantDomain string
	}{
		{"example.com.", "example.com"},
		{"example.com", "example.com"},
	}

	for i, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("TestCase%d", i), func(t *testing.T) {
			if TrimDomain(tc.domain) != tc.wantDomain {
				t.Error("TrimDomain returned an invalid value")
			}
		})
	}
}

func TestMakeDomainCanonical(t *testing.T) {
	testCases := []struct {
		domain     string
		wantDomain string
	}{
		{"example.com.", "example.com."},
		{"example.com", "example.com."},
	}

	for i, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("TestCase%d", i), func(t *testing.T) {
			if MakeDomainCanonical(tc.domain) != tc.wantDomain {
				t.Error("MakeDomainCanonical returned an invalid value")
			}
		})
	}
}
