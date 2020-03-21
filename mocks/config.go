package mocks

import (
	"net/http"

	"github.com/jarcoal/httpmock"
	"github.com/joeig/go-powerdns/v2/types"
)

func (m *Mock) RegisterConfigsMockResponder() {
	httpmock.RegisterResponder("GET", m.generateTestAPIVHostURL()+"/config",
		func(req *http.Request) (*http.Response, error) {
			if res := m.verifyAPIKey(req); res != nil {
				return res, nil
			}

			configMock := []types.ConfigSetting{
				{
					Name:  types.String("signing-threads"),
					Type:  types.String("ConfigSetting"),
					Value: types.String("3"),
				},
			}
			return httpmock.NewJsonResponse(http.StatusOK, configMock)
		},
	)
}
