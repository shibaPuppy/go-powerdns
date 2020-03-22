package mocks

import (
	"fmt"
	"net/http"

	"github.com/jarcoal/httpmock"
	"github.com/joeig/go-powerdns/v2/lib"
)

// RegisterServersMockResponders registers server mock responders
func (m *Mock) RegisterServersMockResponders() {
	httpmock.RegisterResponder("GET", m.generateTestAPIURL()+"/servers",
		func(req *http.Request) (*http.Response, error) {
			if res := m.verifyAPIKey(req); res != nil {
				return res, nil
			}

			serversMock := []lib.Server{
				{
					Type:       lib.StringPtr("Server"),
					ID:         lib.StringPtr(m.TestVHost),
					DaemonType: lib.StringPtr("authoritative"),
					Version:    lib.StringPtr("4.1.2"),
					URL:        lib.StringPtr("/api/v1/servers/" + m.TestVHost),
					ConfigURL:  lib.StringPtr("/api/v1/servers/" + m.TestVHost + "/config{/config_setting}"),
					ZonesURL:   lib.StringPtr("/api/v1/servers/" + m.TestVHost + "/zones{/zone}"),
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
				Type:       lib.StringPtr("Server"),
				ID:         lib.StringPtr(m.TestVHost),
				DaemonType: lib.StringPtr("authoritative"),
				Version:    lib.StringPtr("4.1.2"),
				URL:        lib.StringPtr("/api/v1/servers/" + m.TestVHost),
				ConfigURL:  lib.StringPtr("/api/v1/servers/" + m.TestVHost + "/config{/config_setting}"),
				ZonesURL:   lib.StringPtr("/api/v1/servers/" + m.TestVHost + "/zones{/zone}"),
			}
			return httpmock.NewJsonResponse(http.StatusOK, serverMock)
		},
	)
}

// RegisterCacheFlushMockResponder registers a cache flush mock responder
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
				Count:  lib.Uint32Ptr(1),
				Result: lib.StringPtr("foo"),
			}
			return httpmock.NewJsonResponse(http.StatusOK, cacheFlushResultMock)
		},
	)
}
