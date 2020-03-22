package mocks

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/jarcoal/httpmock"
	"github.com/joeig/go-powerdns/v2/lib"
)

func (m *Mock) generateTestAPIZonesURL() string {
	return fmt.Sprintf("%s/zones", m.generateTestAPIVHostURL())
}

func (m *Mock) generateTestAPIZoneURL(testDomain string) string {
	return fmt.Sprintf("%s/%s", m.generateTestAPIZonesURL(), testDomain)
}

func (m *Mock) generateTestAPIZoneNotifyURL(testDomain string) string {
	return fmt.Sprintf("%s/notify", m.generateTestAPIZoneURL(testDomain))
}

func (m *Mock) generateTestAPIZoneExportURL(testDomain string) string {
	return fmt.Sprintf("%s/export", m.generateTestAPIZoneURL(testDomain))
}

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

// RegisterZonesMockResponder registers a zones mock responder
func (m *Mock) RegisterZonesMockResponder() {
	httpmock.RegisterResponder("GET", m.generateTestAPIZonesURL(),
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
					ID:             lib.StringPtr(lib.MakeDomainCanonical(testDomain)),
					Name:           lib.StringPtr(lib.MakeDomainCanonical(testDomain)),
					URL:            lib.StringPtr(fmt.Sprintf("/api/v1/servers/%s/zones/%s", m.TestVHost, lib.MakeDomainCanonical(testDomain))),
					Kind:           lib.ZoneKindPtr(lib.NativeZoneKind),
					Serial:         lib.Uint32Ptr(1337),
					NotifiedSerial: lib.Uint32Ptr(1337),
				},
			}
			return httpmock.NewJsonResponse(http.StatusOK, zonesMock)
		},
	)
}

// RegisterZoneMockResponders registers zone mock responders
func (m *Mock) RegisterZoneMockResponders(testDomain string, zoneKind lib.ZoneKind) {
	httpmock.RegisterResponder("GET", m.generateTestAPIZoneURL(testDomain),
		func(req *http.Request) (*http.Response, error) {
			if res := m.verifyAPIKey(req); res != nil {
				return res, nil
			}

			if req.Body != nil {
				log.Print("Request body is not nil")
				return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
			}

			zoneMock := lib.Zone{
				ID:   lib.StringPtr(lib.MakeDomainCanonical(testDomain)),
				Name: lib.StringPtr(lib.MakeDomainCanonical(testDomain)),
				URL:  lib.StringPtr(fmt.Sprintf("/api/v1/servers/%s/zones/%s", m.TestVHost, lib.MakeDomainCanonical(testDomain))),
				Kind: lib.ZoneKindPtr(lib.NativeZoneKind),
				RRsets: lib.RRsetSlicePtr([]lib.RRset{
					{
						Name: lib.StringPtr(lib.MakeDomainCanonical(testDomain)),
						Type: lib.RRTypePtr(lib.RRTypeSOA),
						TTL:  lib.Uint32Ptr(3600),
						Records: lib.RecordSlicePtr([]lib.Record{
							{
								Content: lib.StringPtr("a.misconfigured.powerdns.server. hostmaster." + lib.MakeDomainCanonical(testDomain) + " 1337 10800 3600 604800 3600"),
							},
						}),
					},
				}),
				Serial:         lib.Uint32Ptr(1337),
				NotifiedSerial: lib.Uint32Ptr(1337),
			}
			return httpmock.NewJsonResponse(http.StatusOK, zoneMock)
		},
	)

	httpmock.RegisterResponder("POST", m.generateTestAPIZonesURL(),
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
					ID:   lib.StringPtr(lib.MakeDomainCanonical(testDomain)),
					Name: lib.StringPtr(lib.MakeDomainCanonical(testDomain)),
					Type: lib.ZoneTypePtr(lib.ZoneZoneType),
					URL:  lib.StringPtr(fmt.Sprintf("/api/v1/servers/%s/zones/%s", m.TestVHost, lib.MakeDomainCanonical(testDomain))),
					Kind: lib.ZoneKindPtr(zoneKind),
					RRsets: lib.RRsetSlicePtr([]lib.RRset{
						{
							Name: lib.StringPtr(lib.MakeDomainCanonical(testDomain)),
							Type: lib.RRTypePtr(lib.RRTypeSOA),
							TTL:  lib.Uint32Ptr(3600),
							Records: lib.RecordSlicePtr([]lib.Record{
								{
									Content:  lib.StringPtr("a.misconfigured.powerdns.server. hostmaster." + lib.MakeDomainCanonical(testDomain) + " 0 10800 3600 604800 3600"),
									Disabled: lib.BoolPtr(false),
								},
							}),
						},
						{
							Name: lib.StringPtr(lib.MakeDomainCanonical(testDomain)),
							Type: lib.RRTypePtr(lib.RRTypeNS),
							TTL:  lib.Uint32Ptr(3600),
							Records: lib.RecordSlicePtr([]lib.Record{
								{
									Content:  lib.StringPtr("ns.example.tld."),
									Disabled: lib.BoolPtr(false),
								},
							}),
						},
					}),
					Serial:      lib.Uint32Ptr(0),
					Masters:     lib.StringSlicePtr([]string{}),
					DNSsec:      lib.BoolPtr(true),
					Nsec3Param:  lib.StringPtr(""),
					Nsec3Narrow: lib.BoolPtr(false),
					SOAEdit:     lib.StringPtr("foo"),
					SOAEditAPI:  lib.StringPtr("foo"),
					APIRectify:  lib.BoolPtr(true),
					Account:     lib.StringPtr(""),
				}
			case lib.SlaveZoneKind:
				zoneMock = lib.Zone{
					ID:          lib.StringPtr(lib.MakeDomainCanonical(testDomain)),
					Name:        lib.StringPtr(lib.MakeDomainCanonical(testDomain)),
					Type:        lib.ZoneTypePtr(lib.ZoneZoneType),
					URL:         lib.StringPtr(fmt.Sprintf("/api/v1/servers/%s/zones/%s", m.TestVHost, lib.MakeDomainCanonical(testDomain))),
					Kind:        lib.ZoneKindPtr(zoneKind),
					Serial:      lib.Uint32Ptr(0),
					Masters:     lib.StringSlicePtr([]string{"127.0.0.1"}),
					DNSsec:      lib.BoolPtr(true),
					Nsec3Param:  lib.StringPtr(""),
					Nsec3Narrow: lib.BoolPtr(false),
					SOAEdit:     lib.StringPtr(""),
					SOAEditAPI:  lib.StringPtr("DEFAULT"),
					APIRectify:  lib.BoolPtr(true),
					Account:     lib.StringPtr(""),
				}
			default:
				return httpmock.NewStringResponse(http.StatusUnprocessableEntity, "Unprocessable Entity"), nil
			}

			return httpmock.NewJsonResponse(http.StatusCreated, zoneMock)
		},
	)

	httpmock.RegisterResponder("PUT", m.generateTestAPIZoneURL(testDomain),
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

	httpmock.RegisterResponder("DELETE", m.generateTestAPIZoneURL(testDomain),
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

	httpmock.RegisterResponder("PUT", m.generateTestAPIZoneNotifyURL(testDomain),
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

	httpmock.RegisterResponder("GET", m.generateTestAPIZoneExportURL(testDomain),
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
