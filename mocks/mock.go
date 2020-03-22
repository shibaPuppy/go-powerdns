package mocks

import (
	"net/http"

	"github.com/jarcoal/httpmock"
)

// Mock is a structure, which contains a basic mock configuration
type Mock struct {
	TestBaseURL string
	TestVHost   string
	TestAPIKey  string
}

// Activate enables the mock backend
func (m *Mock) Activate() {
	httpmock.Activate()
}

// DeactivateAndReset stops and clears the mock backend
func (m *Mock) DeactivateAndReset() {
	httpmock.DeactivateAndReset()
}

// Disabled returns true if the mock backend is disabled
func (m *Mock) Disabled() bool {
	return httpmock.Disabled()
}

func (m *Mock) verifyAPIKey(req *http.Request) *http.Response {
	if req.Header.Get("X-Api-Key") != m.TestAPIKey {
		return httpmock.NewStringResponse(http.StatusUnauthorized, "Unauthorized")
	}

	return nil
}
