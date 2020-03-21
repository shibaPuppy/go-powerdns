package mocks

import (
	"encoding/json"
	"github.com/jarcoal/httpmock"
	"github.com/joeig/go-powerdns/v2/types"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func validateChangeType(changeType types.ChangeType) error {
	matched, err := regexp.MatchString(`^(REPLACE|DELETE)$`, string(changeType))
	if matched == false || err != nil {
		return &types.Error{}
	}
	return nil
}

func validateRRType(rrType types.RRType) error {
	matched, err := regexp.MatchString(`^(A|AAAA|AFSDB|ALIAS|CAA|CERT|CDNSKEY|CDS|CNAME|DNSKEY|DNAME|DS|HINFO|KEY|LOC|MX|NAPTR|NS|NSEC|NSEC3|NSEC3PARAM|OPENPGPKEY|PTR|RP|RRSIG|SOA|SPF|SSHFP|SRV|TKEY|TSIG|TLSA|SMIMEA|TXT|URI)$`, string(rrType))
	if matched == false || err != nil {
		return &types.Error{}
	}
	return nil
}

func validateCNAMEContent(content string) error {
	if !strings.HasSuffix(content, ".") {
		return &types.Error{}
	}
	return nil
}

func (m *Mock) RegisterRecordMockResponder(testDomain string) {
	httpmock.RegisterResponder("PATCH", m.generateTestAPIVHostURL()+"/zones/"+testDomain,
		func(req *http.Request) (*http.Response, error) {
			if res := m.verifyAPIKey(req); res != nil {
				return res, nil
			}

			if req.Body == nil {
				log.Print("Request body is nil")
				return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
			}

			var rrsets types.RRsets
			if json.NewDecoder(req.Body).Decode(&rrsets) != nil {
				log.Print("Cannot decode request body")
				return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
			}

			for _, set := range rrsets.Sets {
				if validateChangeType(*set.ChangeType) != nil {
					log.Print("Invalid change type", *set.ChangeType)
					return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
				}

				if validateRRType(*set.Type) != nil {
					log.Print("Invalid record type", *set.Type)
					return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
				}

				if *set.Type == types.RRTypeCNAME || *set.Type == types.RRTypeMX {
					for _, record := range set.Records {
						if validateCNAMEContent(*record.Content) != nil {
							log.Print("CNAME content validation failed")
							return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
						}
					}
				}
			}

			zoneMock := types.Zone{
				Name: types.String(types.MakeDomainCanonical(testDomain)),
				URL:  types.String("/api/v1/servers/" + m.TestVHost + "/zones/" + types.MakeDomainCanonical(testDomain)),
			}
			return httpmock.NewJsonResponse(http.StatusOK, zoneMock)
		},
	)
}
