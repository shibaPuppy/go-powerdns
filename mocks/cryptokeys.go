package mocks

import (
	"fmt"
	"net/http"

	"github.com/jarcoal/httpmock"
	"github.com/joeig/go-powerdns/v2/lib"
)

func (m *Mock) RegisterCryptokeysMockResponder(testDomain string) {
	httpmock.RegisterResponder("GET", m.generateTestAPIVHostURL()+"/zones/"+testDomain+"/cryptokeys",
		func(req *http.Request) (*http.Response, error) {
			if res := m.verifyAPIKey(req); res != nil {
				return res, nil
			}

			cryptokeysMock := []lib.Cryptokey{
				{
					Type:      lib.String("Cryptokey"),
					ID:        lib.Uint64(11),
					KeyType:   lib.String("zsk"),
					Active:    lib.Bool(true),
					DNSkey:    lib.String("256 3 8 thisIsTheKey"),
					Algorithm: lib.String("ECDSAP256SHA256"),
					Bits:      lib.Uint64(1024),
				},
				{
					Type:    lib.String("Cryptokey"),
					ID:      lib.Uint64(10),
					KeyType: lib.String("lsk"),
					Active:  lib.Bool(true),
					DNSkey:  lib.String("257 3 8 thisIsTheKey"),
					DS: []string{
						"997 8 1 foo",
						"997 8 2 foo",
						"997 8 4 foo",
					},
					Algorithm: lib.String("ECDSAP256SHA256"),
					Bits:      lib.Uint64(2048),
				},
			}
			return httpmock.NewJsonResponse(http.StatusOK, cryptokeysMock)
		},
	)
}

func (m *Mock) RegisterCryptokeyMockResponder(testDomain string, id uint64) {
	httpmock.RegisterResponder("GET", m.generateTestAPIVHostURL()+"/zones/"+testDomain+"/cryptokeys/"+lib.CryptokeyIDToString(id),
		func(req *http.Request) (*http.Response, error) {
			if res := m.verifyAPIKey(req); res != nil {
				return res, nil
			}

			cryptokeyMock := lib.Cryptokey{
				Type:       lib.String("Cryptokey"),
				ID:         lib.Uint64(0),
				KeyType:    lib.String("zsk"),
				Active:     lib.Bool(true),
				DNSkey:     lib.String("256 3 8 thisIsTheKey"),
				Privatekey: lib.String("Private-key-format: v1.2\nAlgorithm: 8 (ECDSAP256SHA256)\nModulus: foo\nPublicExponent: foo\nPrivateExponent: foo\nPrime1: foo\nPrime2: foo\nExponent1: foo\nExponent2: foo\nCoefficient: foo\n"),
				Algorithm:  lib.String("ECDSAP256SHA256"),
				Bits:       lib.Uint64(1024),
			}
			return httpmock.NewJsonResponse(http.StatusOK, cryptokeyMock)
		},
	)

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("%s/zones/%s/cryptokeys/%s", m.generateTestAPIVHostURL(), testDomain, lib.CryptokeyIDToString(id)),
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == m.TestAPIKey {
				return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
			}
			return httpmock.NewStringResponse(http.StatusUnauthorized, "Unauthorized"), nil
		},
	)
}
