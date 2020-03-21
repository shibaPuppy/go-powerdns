package mocks

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"github.com/jarcoal/httpmock"
	"github.com/joeig/go-powerdns/v2/lib"
)

func validateZoneType(zoneType lib.ZoneType) error {
	if zoneType != "Zone" {
		return &lib.Error{}
	}

	return nil
}

func validateZoneKind(zoneKind lib.ZoneKind) error {
	matched, err := regexp.MatchString(`^(Native|Master|Slave)$`, string(zoneKind))
	if !matched || err != nil {
		return &lib.Error{}
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
			zonesMock := []lib.Zone{
				{
					ID:             lib.String(lib.MakeDomainCanonical(testDomain)),
					Name:           lib.String(lib.MakeDomainCanonical(testDomain)),
					URL:            lib.String("/api/v1/servers/" + m.TestVHost + "/zones/" + lib.MakeDomainCanonical(testDomain)),
					Kind:           lib.ZoneKindPtr(lib.NativeZoneKind),
					Serial:         lib.Uint32(1337),
					NotifiedSerial: lib.Uint32(1337),
				},
			}
			return httpmock.NewJsonResponse(http.StatusOK, zonesMock)
		},
	)
}

func (m *Mock) RegisterZoneMockResponder(testDomain string, zoneKind lib.ZoneKind) {
	httpmock.RegisterResponder("GET", m.generateTestAPIVHostURL()+"/zones/"+testDomain,
		func(req *http.Request) (*http.Response, error) {
			if res := m.verifyAPIKey(req); res != nil {
				return res, nil
			}

			if req.Body != nil {
				log.Print("Request body is not nil")
				return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
			}

			zoneMock := lib.Zone{
				ID:   lib.String(lib.MakeDomainCanonical(testDomain)),
				Name: lib.String(lib.MakeDomainCanonical(testDomain)),
				URL:  lib.String("/api/v1/servers/" + m.TestVHost + "/zones/" + lib.MakeDomainCanonical(testDomain)),
				Kind: lib.ZoneKindPtr(lib.NativeZoneKind),
				RRsets: []lib.RRset{
					{
						Name: lib.String(lib.MakeDomainCanonical(testDomain)),
						Type: lib.RRTypePtr(lib.RRTypeSOA),
						TTL:  lib.Uint32(3600),
						Records: []lib.Record{
							{
								Content: lib.String("a.misconfigured.powerdns.server. hostmaster." + lib.MakeDomainCanonical(testDomain) + " 1337 10800 3600 604800 3600"),
							},
						},
					},
				},
				Serial:         lib.Uint32(1337),
				NotifiedSerial: lib.Uint32(1337),
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

			var zone lib.Zone
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

			var zoneMock lib.Zone
			switch zoneKind {
			case lib.NativeZoneKind, lib.MasterZoneKind:
				zoneMock = lib.Zone{
					ID:   lib.String(lib.MakeDomainCanonical(testDomain)),
					Name: lib.String(lib.MakeDomainCanonical(testDomain)),
					Type: lib.ZoneTypePtr(lib.ZoneZoneType),
					URL:  lib.String("api/v1/servers/" + m.TestVHost + "/zones/" + lib.MakeDomainCanonical(testDomain)),
					Kind: lib.ZoneKindPtr(zoneKind),
					RRsets: []lib.RRset{
						{
							Name: lib.String(lib.MakeDomainCanonical(testDomain)),
							Type: lib.RRTypePtr(lib.RRTypeSOA),
							TTL:  lib.Uint32(3600),
							Records: []lib.Record{
								{
									Content:  lib.String("a.misconfigured.powerdns.server. hostmaster." + lib.MakeDomainCanonical(testDomain) + " 0 10800 3600 604800 3600"),
									Disabled: lib.Bool(false),
								},
							},
						},
						{
							Name: lib.String(lib.MakeDomainCanonical(testDomain)),
							Type: lib.RRTypePtr(lib.RRTypeNS),
							TTL:  lib.Uint32(3600),
							Records: []lib.Record{
								{
									Content:  lib.String("ns.example.tld."),
									Disabled: lib.Bool(false),
								},
							},
						},
					},
					Serial:      lib.Uint32(0),
					Masters:     []string{},
					DNSsec:      lib.Bool(true),
					Nsec3Param:  lib.String(""),
					Nsec3Narrow: lib.Bool(false),
					SOAEdit:     lib.String("foo"),
					SOAEditAPI:  lib.String("foo"),
					APIRectify:  lib.Bool(true),
					Account:     lib.String(""),
				}
			case lib.SlaveZoneKind:
				zoneMock = lib.Zone{
					ID:          lib.String(lib.MakeDomainCanonical(testDomain)),
					Name:        lib.String(lib.MakeDomainCanonical(testDomain)),
					Type:        lib.ZoneTypePtr(lib.ZoneZoneType),
					URL:         lib.String("api/v1/servers/" + m.TestVHost + "/zones/" + lib.MakeDomainCanonical(testDomain)),
					Kind:        lib.ZoneKindPtr(zoneKind),
					Serial:      lib.Uint32(0),
					Masters:     []string{"127.0.0.1"},
					DNSsec:      lib.Bool(true),
					Nsec3Param:  lib.String(""),
					Nsec3Narrow: lib.Bool(false),
					SOAEdit:     lib.String(""),
					SOAEditAPI:  lib.String("DEFAULT"),
					APIRectify:  lib.Bool(true),
					Account:     lib.String(""),
				}
			default:
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

			return httpmock.NewStringResponse(http.StatusOK, lib.MakeDomainCanonical(testDomain)+"	3600	SOA	a.misconfigured.powerdns.server. hostmaster."+lib.MakeDomainCanonical(testDomain)+" 1 10800 3600 604800 3600"), nil
		},
	)
}
