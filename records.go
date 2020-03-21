package powerdns

import (
	"fmt"

	"github.com/joeig/go-powerdns/v2/types"
)

// RecordsService handles communication with the records related methods of the Client API
type RecordsService service

// Add creates a new resource record
func (r *RecordsService) Add(domain string, name string, recordType types.RRType, ttl uint32, content []string) error {
	return r.Change(domain, name, recordType, ttl, content)
}

// Change replaces an existing resource record
func (r *RecordsService) Change(domain string, name string, recordType types.RRType, ttl uint32, content []string) error {
	rrset := new(types.RRset)
	rrset.Name = &name
	rrset.Type = &recordType
	rrset.TTL = &ttl
	rrset.ChangeType = types.ChangeTypePtr(types.ChangeTypeReplace)
	rrset.Records = make([]types.Record, 0)

	for _, c := range content {
		r := types.Record{Content: types.String(c), Disabled: types.Bool(false), SetPTR: types.Bool(false)}
		rrset.Records = append(rrset.Records, r)
	}

	return r.patchRRset(domain, *rrset)
}

// Delete removes an existing resource record
func (r *RecordsService) Delete(domain string, name string, recordType types.RRType) error {
	rrset := new(types.RRset)
	rrset.Name = &name
	rrset.Type = &recordType
	rrset.ChangeType = types.ChangeTypePtr(types.ChangeTypeDelete)

	return r.patchRRset(domain, *rrset)
}

func canonicalResourceRecordValues(records []types.Record) {
	for i := range records {
		records[i].Content = types.String(types.MakeDomainCanonical(*records[i].Content))
	}
}

func fixRRset(rrset *types.RRset) {
	if *rrset.Type != types.RRTypeCNAME && *rrset.Type != types.RRTypeMX {
		return
	}

	canonicalResourceRecordValues(rrset.Records)
}

func (r *RecordsService) patchRRset(domain string, rrset types.RRset) error {
	rrset.Name = types.String(types.MakeDomainCanonical(*rrset.Name))

	fixRRset(&rrset)

	payload := types.RRsets{}
	payload.Sets = append(payload.Sets, rrset)

	req, err := r.client.newRequest("PATCH", fmt.Sprintf("servers/%s/zones/%s", r.client.VHost, types.TrimDomain(domain)), nil, payload)
	if err != nil {
		return err
	}

	_, err = r.client.do(req, nil)

	return err
}
