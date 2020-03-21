package mocks

import (
	"encoding/json"
	"github.com/jarcoal/httpmock"
	"github.com/joeig/go-powerdns/v2/types"
	"log"
	"net/http"
	"regexp"
)

func validateZoneType(zoneType types.ZoneType) error {
	if zoneType != "Zone" {
		return &types.Error{}
	}
	return nil
}

func validateZoneKind(zoneKind types.ZoneKind) error {
	matched, err := regexp.MatchString(`^(Native|Master|Slave)$`, string(zoneKind))
	if matched == false || err != nil {
		return &types.Error{}
	}
	return nil
}

func (m *Mock) RegisterZonesMockResponder() {
	httpmock.RegisterResponder("GET", m.generateTestAPIVHostURL()+"/zones",
		func(req *http.Request) (*http.Response, error) {
			if res := m.verifyAPIKey(req); res != nil {
				return res, nil
			}

			if req.Body != nil {
				log.Print("Request body is not nil")
				return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
			}

			testDomain := "example.com"
			zonesMock := []types.Zone{
				{
					ID:             types.String(types.MakeDomainCanonical(testDomain)),
					Name:           types.String(types.MakeDomainCanonical(testDomain)),
					URL:            types.String("/api/v1/servers/" + m.TestVHost + "/zones/" + types.MakeDomainCanonical(testDomain)),
					Kind:           types.ZoneKindPtr(types.NativeZoneKind),
					Serial:         types.Uint32(1337),
					NotifiedSerial: types.Uint32(1337),
				},
			}
			return httpmock.NewJsonResponse(http.StatusOK, zonesMock)
		},
	)
}

func (m *Mock) RegisterZoneMockResponder(testDomain string, zoneKind types.ZoneKind) {
	httpmock.RegisterResponder("GET", m.generateTestAPIVHostURL()+"/zones/"+testDomain,
		func(req *http.Request) (*http.Response, error) {
			if res := m.verifyAPIKey(req); res != nil {
				return res, nil
			}

			if req.Body != nil {
				log.Print("Request body is not nil")
				return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
			}

			zoneMock := types.Zone{
				ID:   types.String(types.MakeDomainCanonical(testDomain)),
				Name: types.String(types.MakeDomainCanonical(testDomain)),
				URL:  types.String("/api/v1/servers/" + m.TestVHost + "/zones/" + types.MakeDomainCanonical(testDomain)),
				Kind: types.ZoneKindPtr(types.NativeZoneKind),
				RRsets: []types.RRset{
					{
						Name: types.String(types.MakeDomainCanonical(testDomain)),
						Type: types.RRTypePtr(types.RRTypeSOA),
						TTL:  types.Uint32(3600),
						Records: []types.Record{
							{
								Content: types.String("a.misconfigured.powerdns.server. hostmaster." + types.MakeDomainCanonical(testDomain) + " 1337 10800 3600 604800 3600"),
							},
						},
					},
				},
				Serial:         types.Uint32(1337),
				NotifiedSerial: types.Uint32(1337),
			}
			return httpmock.NewJsonResponse(http.StatusOK, zoneMock)
		},
	)

	httpmock.RegisterResponder("POST", m.generateTestAPIVHostURL()+"/zones",
		func(req *http.Request) (*http.Response, error) {
			if res := m.verifyAPIKey(req); res != nil {
				return res, nil
			}

			if req.Body == nil {
				log.Print("Request body is nil")
				return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
			}

			var zone types.Zone
			if json.NewDecoder(req.Body).Decode(&zone) != nil {
				log.Print("Cannot decode request body")
				return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
			}

			if validateZoneType(*zone.Type) != nil {
				log.Print("Invalid zone type", *zone.Type)
				return httpmock.NewStringResponse(http.StatusUnprocessableEntity, "Unprocessable Entity"), nil
			}

			if validateZoneKind(*zone.Kind) != nil {
				log.Print("Invalid zone kind", *zone.Kind)
				return httpmock.NewStringResponse(http.StatusUnprocessableEntity, "Unprocessable Entity"), nil
			}

			var zoneMock types.Zone
			if zoneKind == types.NativeZoneKind || zoneKind == types.MasterZoneKind {
				zoneMock = types.Zone{
					ID:   types.String(types.MakeDomainCanonical(testDomain)),
					Name: types.String(types.MakeDomainCanonical(testDomain)),
					Type: types.ZoneTypePtr(types.ZoneZoneType),
					URL:  types.String("api/v1/servers/" + m.TestVHost + "/zones/" + types.MakeDomainCanonical(testDomain)),
					Kind: types.ZoneKindPtr(zoneKind),
					RRsets: []types.RRset{
						{
							Name: types.String(types.MakeDomainCanonical(testDomain)),
							Type: types.RRTypePtr(types.RRTypeSOA),
							TTL:  types.Uint32(3600),
							Records: []types.Record{
								{
									Content:  types.String("a.misconfigured.powerdns.server. hostmaster." + types.MakeDomainCanonical(testDomain) + " 0 10800 3600 604800 3600"),
									Disabled: types.Bool(false),
								},
							},
						},
						{
							Name: types.String(types.MakeDomainCanonical(testDomain)),
							Type: types.RRTypePtr(types.RRTypeNS),
							TTL:  types.Uint32(3600),
							Records: []types.Record{
								{
									Content:  types.String("ns.example.tld."),
									Disabled: types.Bool(false),
								},
							},
						},
					},
					Serial:      types.Uint32(0),
					Masters:     []string{},
					DNSsec:      types.Bool(true),
					Nsec3Param:  types.String(""),
					Nsec3Narrow: types.Bool(false),
					SOAEdit:     types.String("foo"),
					SOAEditAPI:  types.String("foo"),
					APIRectify:  types.Bool(true),
					Account:     types.String(""),
				}
			} else if zoneKind == types.SlaveZoneKind {
				zoneMock = types.Zone{
					ID:          types.String(types.MakeDomainCanonical(testDomain)),
					Name:        types.String(types.MakeDomainCanonical(testDomain)),
					Type:        types.ZoneTypePtr(types.ZoneZoneType),
					URL:         types.String("api/v1/servers/" + m.TestVHost + "/zones/" + types.MakeDomainCanonical(testDomain)),
					Kind:        types.ZoneKindPtr(zoneKind),
					Serial:      types.Uint32(0),
					Masters:     []string{"127.0.0.1"},
					DNSsec:      types.Bool(true),
					Nsec3Param:  types.String(""),
					Nsec3Narrow: types.Bool(false),
					SOAEdit:     types.String(""),
					SOAEditAPI:  types.String("DEFAULT"),
					APIRectify:  types.Bool(true),
					Account:     types.String(""),
				}
			} else {
				return httpmock.NewStringResponse(http.StatusUnprocessableEntity, "Unprocessable Entity"), nil
			}

			return httpmock.NewJsonResponse(http.StatusCreated, zoneMock)
		},
	)

	httpmock.RegisterResponder("PUT", m.generateTestAPIVHostURL()+"/zones/"+testDomain,
		func(req *http.Request) (*http.Response, error) {
			if res := m.verifyAPIKey(req); res != nil {
				return res, nil
			}

			if req.Body == nil {
				log.Print("Request body is nil")
				return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
			}

			return httpmock.NewBytesResponse(http.StatusNoContent, []byte{}), nil
		},
	)

	httpmock.RegisterResponder("DELETE", m.generateTestAPIVHostURL()+"/zones/"+testDomain,
		func(req *http.Request) (*http.Response, error) {
			if res := m.verifyAPIKey(req); res != nil {
				return res, nil
			}

			if req.Body != nil {
				log.Print("Request body is not nil")
				return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
			}

			return httpmock.NewBytesResponse(http.StatusNoContent, []byte{}), nil
		},
	)

	httpmock.RegisterResponder("PUT", m.generateTestAPIVHostURL()+"/zones/"+testDomain+"/notify",
		func(req *http.Request) (*http.Response, error) {
			if res := m.verifyAPIKey(req); res != nil {
				return res, nil
			}

			if req.Body != nil {
				log.Print("Request body is not nil")
				return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
			}

			return httpmock.NewStringResponse(http.StatusOK, "{\"result\":\"Notification queued\"}"), nil
		},
	)

	httpmock.RegisterResponder("GET", m.generateTestAPIVHostURL()+"/zones/"+testDomain+"/export",
		func(req *http.Request) (*http.Response, error) {
			if res := m.verifyAPIKey(req); res != nil {
				return res, nil
			}

			if req.Body != nil {
				log.Print("Request body is not nil")
				return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
			}

			return httpmock.NewStringResponse(http.StatusOK, types.MakeDomainCanonical(testDomain)+"	3600	SOA	a.misconfigured.powerdns.server. hostmaster."+types.MakeDomainCanonical(testDomain)+" 1 10800 3600 604800 3600"), nil
		},
	)
}
