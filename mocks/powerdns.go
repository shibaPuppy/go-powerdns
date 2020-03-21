package mocks

import (
	"fmt"
	"net/http"

	"github.com/jarcoal/httpmock"
	"github.com/joeig/go-powerdns/v2/lib"
)

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
			mock := lib.Error{
				Status:  "Not Found",
				Message: "Not Found",
			}
			return httpmock.NewJsonResponse(http.StatusNotImplemented, mock)
		},
	)
}
