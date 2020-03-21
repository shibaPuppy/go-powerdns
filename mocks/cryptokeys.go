package mocks

import (
	"fmt"
	"github.com/jarcoal/httpmock"
	"github.com/joeig/go-powerdns/v2/types"
	"net/http"
)

func (m *Mock) RegisterCryptokeysMockResponder(testDomain string) {
	httpmock.RegisterResponder("GET", m.generateTestAPIVHostURL()+"/zones/"+testDomain+"/cryptokeys",
		func(req *http.Request) (*http.Response, error) {
			if res := m.verifyAPIKey(req); res != nil {
				return res, nil
			}

			cryptokeysMock := []types.Cryptokey{
				{
					Type:      types.String("Cryptokey"),
					ID:        types.Uint64(11),
					KeyType:   types.String("zsk"),
					Active:    types.Bool(true),
					DNSkey:    types.String("256 3 8 thisIsTheKey"),
					Algorithm: types.String("ECDSAP256SHA256"),
					Bits:      types.Uint64(1024),
				},
				{
					Type:    types.String("Cryptokey"),
					ID:      types.Uint64(10),
					KeyType: types.String("lsk"),
					Active:  types.Bool(true),
					DNSkey:  types.String("257 3 8 thisIsTheKey"),
					DS: []string{
						"997 8 1 foo",
						"997 8 2 foo",
						"997 8 4 foo",
					},
					Algorithm: types.String("ECDSAP256SHA256"),
					Bits:      types.Uint64(2048),
				},
			}
			return httpmock.NewJsonResponse(http.StatusOK, cryptokeysMock)
		},
	)
}

func (m *Mock) RegisterCryptokeyMockResponder(testDomain string, id uint64) {
	httpmock.RegisterResponder("GET", m.generateTestAPIVHostURL()+"/zones/"+testDomain+"/cryptokeys/"+types.CryptokeyIDToString(id),
		func(req *http.Request) (*http.Response, error) {
			if res := m.verifyAPIKey(req); res != nil {
				return res, nil
			}

			cryptokeyMock := types.Cryptokey{
				Type:       types.String("Cryptokey"),
				ID:         types.Uint64(0),
				KeyType:    types.String("zsk"),
				Active:     types.Bool(true),
				DNSkey:     types.String("256 3 8 thisIsTheKey"),
				Privatekey: types.String("Private-key-format: v1.2\nAlgorithm: 8 (ECDSAP256SHA256)\nModulus: foo\nPublicExponent: foo\nPrivateExponent: foo\nPrime1: foo\nPrime2: foo\nExponent1: foo\nExponent2: foo\nCoefficient: foo\n"),
				Algorithm:  types.String("ECDSAP256SHA256"),
				Bits:       types.Uint64(1024),
			}
			return httpmock.NewJsonResponse(http.StatusOK, cryptokeyMock)
		},
	)

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("%s/zones/%s/cryptokeys/%s", m.generateTestAPIVHostURL(), testDomain, types.CryptokeyIDToString(id)),
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == m.TestAPIKey {
				return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
			}
			return httpmock.NewStringResponse(http.StatusUnauthorized, "Unauthorized"), nil
		},
	)
}
