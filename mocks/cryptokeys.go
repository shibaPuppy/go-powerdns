package mocks

import (
	"fmt"
	"net/http"

	"github.com/jarcoal/httpmock"
	"github.com/joeig/go-powerdns/v2/lib"
)

func (m *Mock) generateTestAPICryptokeysURL(testDomain string) string {
	return fmt.Sprintf("%s/cryptokeys", m.generateTestAPIZoneURL(testDomain))
}

func (m *Mock) generateTestAPICryptokeyURL(testDomain string, id uint64) string {
	return fmt.Sprintf("%s/%s", m.generateTestAPICryptokeysURL(testDomain), lib.CryptokeyIDToString(id))
}

// RegisterCryptokeysMockResponder registers a cryptokeys mock route
func (m *Mock) RegisterCryptokeysMockResponder(testDomain string) {
	httpmock.RegisterResponder("GET", m.generateTestAPICryptokeysURL(testDomain),
		func(req *http.Request) (*http.Response, error) {
			if res := m.verifyAPIKey(req); res != nil {
				return res, nil
			}

			cryptokeysMock := []lib.Cryptokey{
				{
					Type:      lib.StringPtr("Cryptokey"),
					ID:        lib.Uint64Ptr(11),
					KeyType:   lib.StringPtr("zsk"),
					Active:    lib.BoolPtr(true),
					DNSkey:    lib.StringPtr("256 3 8 thisIsTheKey"),
					Algorithm: lib.StringPtr("ECDSAP256SHA256"),
					Bits:      lib.Uint64Ptr(1024),
				},
				{
					Type:    lib.StringPtr("Cryptokey"),
					ID:      lib.Uint64Ptr(10),
					KeyType: lib.StringPtr("lsk"),
					Active:  lib.BoolPtr(true),
					DNSkey:  lib.StringPtr("257 3 8 thisIsTheKey"),
					DS: lib.StringSlicePtr([]string{
						"997 8 1 foo",
						"997 8 2 foo",
						"997 8 4 foo",
					}),
					Algorithm: lib.StringPtr("ECDSAP256SHA256"),
					Bits:      lib.Uint64Ptr(2048),
				},
			}
			return httpmock.NewJsonResponse(http.StatusOK, cryptokeysMock)
		},
	)
}

// RegisterCryptokeyMockResponders registers cryptokey mock routes
func (m *Mock) RegisterCryptokeyMockResponders(testDomain string, id uint64) {
	httpmock.RegisterResponder("GET", m.generateTestAPICryptokeyURL(testDomain, id),
		func(req *http.Request) (*http.Response, error) {
			if res := m.verifyAPIKey(req); res != nil {
				return res, nil
			}

			cryptokeyMock := lib.Cryptokey{
				Type:       lib.StringPtr("Cryptokey"),
				ID:         lib.Uint64Ptr(0),
				KeyType:    lib.StringPtr("zsk"),
				Active:     lib.BoolPtr(true),
				DNSkey:     lib.StringPtr("256 3 8 thisIsTheKey"),
				Privatekey: lib.StringPtr("Private-key-format: v1.2\nAlgorithm: 8 (ECDSAP256SHA256)\nModulus: foo\nPublicExponent: foo\nPrivateExponent: foo\nPrime1: foo\nPrime2: foo\nExponent1: foo\nExponent2: foo\nCoefficient: foo\n"),
				Algorithm:  lib.StringPtr("ECDSAP256SHA256"),
				Bits:       lib.Uint64Ptr(1024),
			}
			return httpmock.NewJsonResponse(http.StatusOK, cryptokeyMock)
		},
	)

	httpmock.RegisterResponder("DELETE", m.generateTestAPICryptokeyURL(testDomain, id),
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == m.TestAPIKey {
				return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
			}
			return httpmock.NewStringResponse(http.StatusUnauthorized, "Unauthorized"), nil
		},
	)
}
