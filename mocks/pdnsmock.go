package mocks

import (
	"fmt"
	"net/http"

	"github.com/jarcoal/httpmock"
)

type Mock struct {
	TestBaseURL string
	TestVHost   string
	TestAPIKey  string
}

func (m *Mock) Activate() {
	httpmock.Activate()
}

func (m *Mock) DeactivateAndReset() {
	httpmock.DeactivateAndReset()
}

func (m *Mock) Disabled() bool {
	return httpmock.Disabled()
}

func (m *Mock) generateTestAPIURL() string {
	return fmt.Sprintf("%s/api/v1", m.TestBaseURL)
}

func (m *Mock) generateTestAPIVHostURL() string {
	return fmt.Sprintf("%s/servers/%s", m.generateTestAPIURL(), m.TestVHost)
}

func (m *Mock) verifyAPIKey(req *http.Request) *http.Response {
	if req.Header.Get("X-Api-Key") != m.TestAPIKey {
		return httpmock.NewStringResponse(http.StatusUnauthorized, "Unauthorized")
	}

	return nil
}
