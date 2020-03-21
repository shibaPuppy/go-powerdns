package mocks

import (
	"fmt"
	"net/http"

	"github.com/jarcoal/httpmock"
	"github.com/joeig/go-powerdns/v2/lib"
)

func (m *Mock) RegisterServersMockResponder() {
	httpmock.RegisterResponder("GET", m.generateTestAPIURL()+"/servers",
		func(req *http.Request) (*http.Response, error) {
			if res := m.verifyAPIKey(req); res != nil {
				return res, nil
			}

			serversMock := []lib.Server{
				{
					Type:       lib.String("Server"),
					ID:         lib.String(m.TestVHost),
					DaemonType: lib.String("authoritative"),
					Version:    lib.String("4.1.2"),
					URL:        lib.String("/api/v1/servers/" + m.TestVHost),
					ConfigURL:  lib.String("/api/v1/servers/" + m.TestVHost + "/config{/config_setting}"),
					ZonesURL:   lib.String("/api/v1/servers/" + m.TestVHost + "/zones{/zone}"),
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

			serverMock := lib.Server{
				Type:       lib.String("Server"),
				ID:         lib.String(m.TestVHost),
				DaemonType: lib.String("authoritative"),
				Version:    lib.String("4.1.2"),
				URL:        lib.String("/api/v1/servers/" + m.TestVHost),
				ConfigURL:  lib.String("/api/v1/servers/" + m.TestVHost + "/config{/config_setting}"),
				ZonesURL:   lib.String("/api/v1/servers/" + m.TestVHost + "/zones{/zone}"),
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

			if req.URL.Query().Get("domain") != lib.MakeDomainCanonical(testDomain) {
				return httpmock.NewStringResponse(http.StatusUnprocessableEntity, "Unprocessable Eneity"), nil
			}

			cacheFlushResultMock := lib.CacheFlushResult{
				Count:  lib.Uint32(1),
				Result: lib.String("foo"),
			}
			return httpmock.NewJsonResponse(http.StatusOK, cacheFlushResultMock)
		},
	)
}
