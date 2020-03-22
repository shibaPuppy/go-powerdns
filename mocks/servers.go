package mocks

import (
	"fmt"
	"net/http"

	"github.com/jarcoal/httpmock"
	"github.com/joeig/go-powerdns/v2/lib"
)

func (m *Mock) generateTestAPIServersURL() string {
	return fmt.Sprintf("%s/servers", m.generateTestAPIURL())
}

func (m *Mock) generateTestAPIVHostURL() string {
	return fmt.Sprintf("%s/%s", m.generateTestAPIServersURL(), m.TestVHost)
}

func (m *Mock) generateTestAPICustomVHostURL(vHost string) string {
	return fmt.Sprintf("%s/%s", m.generateTestAPIServersURL(), vHost)
}

func (m *Mock) generateTestAPIVHostCacheFlushURL() string {
	return fmt.Sprintf("%s/%s/cache/flush", m.generateTestAPIServersURL(), m.TestVHost)
}

func (m *Mock) generateTestServer() *lib.Server {
	return &lib.Server{
		Type:       lib.StringPtr("Server"),
		ID:         lib.StringPtr(m.TestVHost),
		DaemonType: lib.StringPtr("authoritative"),
		Version:    lib.StringPtr("4.1.2"),
		URL:        lib.StringPtr(fmt.Sprintf("/api/v1/servers/%s", m.TestVHost)),
		ConfigURL:  lib.StringPtr(fmt.Sprintf("/api/v1/servers/%s/config{/config_setting}", m.TestVHost)),
		ZonesURL:   lib.StringPtr(fmt.Sprintf("/api/v1/servers/%s/zones{/zone}", m.TestVHost)),
	}
}

// RegisterServersMockResponders registers server mock responders
func (m *Mock) RegisterServersMockResponders() {
	httpmock.RegisterResponder("GET", m.generateTestAPIServersURL(),
		func(req *http.Request) (*http.Response, error) {
			if res := m.verifyAPIKey(req); res != nil {
				return res, nil
			}

			return httpmock.NewJsonResponse(http.StatusOK, []lib.Server{*m.generateTestServer()})
		},
	)

	httpmock.RegisterResponder("GET", m.generateTestAPIVHostURL(),
		func(req *http.Request) (*http.Response, error) {
			if res := m.verifyAPIKey(req); res != nil {
				return res, nil
			}

			return httpmock.NewJsonResponse(http.StatusOK, *m.generateTestServer())
		},
	)
}

// RegisterCacheFlushMockResponder registers a cache flush mock responder
func (m *Mock) RegisterCacheFlushMockResponder(testDomain string) {
	httpmock.RegisterResponder("PUT", m.generateTestAPIVHostCacheFlushURL(),
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
