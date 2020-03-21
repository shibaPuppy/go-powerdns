package mocks

import (
	"fmt"
	"github.com/jarcoal/httpmock"
	"github.com/joeig/go-powerdns/v2/types"
	"net/http"
)

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

func (m *Mock) RegisterDoMockResponder() {
	httpmock.RegisterResponder("GET", fmt.Sprintf("%s/servers/doesntExist", m.generateTestAPIURL()),
		func(req *http.Request) (*http.Response, error) {
			if res := m.verifyAPIKey(req); res != nil {
				return res, nil
			}
			return httpmock.NewStringResponse(http.StatusNotFound, "Not Found"), nil
		},
	)

	httpmock.RegisterResponder("GET", fmt.Sprintf("%s/server", m.generateTestAPIURL()),
		func(req *http.Request) (*http.Response, error) {
			mock := types.Error{
				Status:  "Not Found",
				Message: "Not Found",
			}
			return httpmock.NewJsonResponse(http.StatusNotImplemented, mock)
		},
	)
}
