package mocks

import (
	"fmt"
	"github.com/jarcoal/httpmock"
	"github.com/joeig/go-powerdns/v2/types"
	"net/http"
)

func (m *Mock) RegisterServersMockResponder() {
	httpmock.RegisterResponder("GET", m.generateTestAPIURL()+"/servers",
		func(req *http.Request) (*http.Response, error) {
			if res := m.verifyAPIKey(req); res != nil {
				return res, nil
			}

			serversMock := []types.Server{
				{
					Type:       types.String("Server"),
					ID:         types.String(m.TestVHost),
					DaemonType: types.String("authoritative"),
					Version:    types.String("4.1.2"),
					URL:        types.String("/api/v1/servers/" + m.TestVHost),
					ConfigURL:  types.String("/api/v1/servers/" + m.TestVHost + "/config{/config_setting}"),
					ZonesURL:   types.String("/api/v1/servers/" + m.TestVHost + "/zones{/zone}"),
				},
			}
			return httpmock.NewJsonResponse(http.StatusOK, serversMock)
		},
	)

	httpmock.RegisterResponder("GET", m.generateTestAPIVHostURL(),
		func(req *http.Request) (*http.Response, error) {
			if res := m.verifyAPIKey(req); res != nil {
				return res, nil
			}

			serverMock := types.Server{
				Type:       types.String("Server"),
				ID:         types.String(m.TestVHost),
				DaemonType: types.String("authoritative"),
				Version:    types.String("4.1.2"),
				URL:        types.String("/api/v1/servers/" + m.TestVHost),
				ConfigURL:  types.String("/api/v1/servers/" + m.TestVHost + "/config{/config_setting}"),
				ZonesURL:   types.String("/api/v1/servers/" + m.TestVHost + "/zones{/zone}"),
			}
			return httpmock.NewJsonResponse(http.StatusOK, serverMock)
		},
	)
}

func (m *Mock) RegisterCacheFlushMockResponder(testDomain string) {
	httpmock.RegisterResponder("PUT", fmt.Sprintf("%s/cache/flush", m.generateTestAPIVHostURL()),
		func(req *http.Request) (*http.Response, error) {
			if res := m.verifyAPIKey(req); res != nil {
				return res, nil
			}

			if req.URL.Query().Get("domain") != types.MakeDomainCanonical(testDomain) {
				return httpmock.NewStringResponse(http.StatusUnprocessableEntity, "Unprocessable Eneity"), nil
			}

			cacheFlushResultMock := types.CacheFlushResult{
				Count:  types.Uint32(1),
				Result: types.String("foo"),
			}
			return httpmock.NewJsonResponse(http.StatusOK, cacheFlushResultMock)
		},
	)
}
