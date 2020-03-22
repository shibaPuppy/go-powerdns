package mocks

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/jarcoal/httpmock"
	"github.com/joeig/go-powerdns/v2/lib"
)

func validateChangeType(changeType lib.ChangeType) error {
	matched, err := regexp.MatchString(`^(REPLACE|DELETE)$`, string(changeType))
	if !matched || err != nil {
		return &lib.Error{}
	}

	return nil
}

func validateRRType(rrType lib.RRType) error {
	matched, err := regexp.MatchString(`^(A|AAAA|AFSDB|ALIAS|CAA|CERT|CDNSKEY|CDS|CNAME|DNSKEY|DNAME|DS|HINFO|KEY|LOC|MX|NAPTR|NS|NSEC|NSEC3|NSEC3PARAM|OPENPGPKEY|PTR|RP|RRSIG|SOA|SPF|SSHFP|SRV|TKEY|TSIG|TLSA|SMIMEA|TXT|URI)$`, string(rrType))
	if !matched || err != nil {
		return &lib.Error{}
	}

	return nil
}

func validateCNAMEContent(content string) error {
	if !strings.HasSuffix(content, ".") {
		return &lib.Error{}
	}

	return nil
}

// RegisterRecordMockResponder registers a record mock responder
func (m *Mock) RegisterRecordMockResponder(testDomain string) {
	httpmock.RegisterResponder("PATCH", m.generateTestAPIZoneURL(testDomain),
		func(req *http.Request) (*http.Response, error) {
			if res := m.verifyAPIKey(req); res != nil {
				return res, nil
			}

			if req.Body == nil {
				log.Print("Request body is nil")
				return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
			}

			var rrsets lib.RRsets
			if json.NewDecoder(req.Body).Decode(&rrsets) != nil {
				log.Print("Cannot decode request body")
				return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
			}

			for _, set := range *rrsets.Sets {
				if validateChangeType(*set.ChangeType) != nil {
					log.Print("Invalid change type", *set.ChangeType)
					return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
				}

				if validateRRType(*set.Type) != nil {
					log.Print("Invalid record type", *set.Type)
					return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
				}

				if *set.Type == lib.RRTypeCNAME || *set.Type == lib.RRTypeMX {
					for _, record := range *set.Records {
						if validateCNAMEContent(*record.Content) != nil {
							log.Print("CNAME content validation failed")
							return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
						}
					}
				}
			}

			zoneMock := m.generateTestZone(testDomain, lib.NativeZoneKind)

			return httpmock.NewJsonResponse(http.StatusOK, zoneMock)
		},
	)
}
