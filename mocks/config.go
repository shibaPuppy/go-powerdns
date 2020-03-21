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
					Name:  lib.String("signing-threads"),
					Type:  lib.String("ConfigSetting"),
					Value: lib.String("3"),
				},
			}
			return httpmock.NewJsonResponse(http.StatusOK, configMock)
		},
	)
}
