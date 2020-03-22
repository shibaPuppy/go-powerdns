package mocks

import (
	"net/http"

	"github.com/jarcoal/httpmock"
	"github.com/joeig/go-powerdns/v2/lib"
)

// RegisterConfigsMockResponder registers a config mock route
func (m *Mock) RegisterConfigsMockResponder() {
	httpmock.RegisterResponder("GET", m.generateTestAPIVHostURL()+"/config",
		func(req *http.Request) (*http.Response, error) {
			if res := m.verifyAPIKey(req); res != nil {
				return res, nil
			}

			configMock := []lib.ConfigSetting{
				{
					Name:  lib.StringPtr("signing-threads"),
					Type:  lib.StringPtr("ConfigSetting"),
					Value: lib.StringPtr("3"),
				},
			}
			return httpmock.NewJsonResponse(http.StatusOK, configMock)
		},
	)
}
