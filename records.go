package powerdns

import (
	"fmt"

	"github.com/joeig/go-powerdns/v2/lib"
)

// RecordsService handles communication with the records related methods of the Client API
type RecordsService service

// Add creates a new resource record
func (r *RecordsService) Add(domain string, name string, recordType lib.RRType, ttl uint32, content []string) error {
	return r.Change(domain, name, recordType, ttl, content)
}

// Change replaces an existing resource record
func (r *RecordsService) Change(domain string, name string, recordType lib.RRType, ttl uint32, content []string) error {
	rrset := new(lib.RRset)
	rrset.Name = &name
	rrset.Type = &recordType
	rrset.TTL = &ttl
	rrset.ChangeType = lib.ChangeTypePtr(lib.ChangeTypeReplace)
	rrset.Records = make([]lib.Record, 0)

	for _, c := range content {
		r := lib.Record{Content: lib.String(c), Disabled: lib.Bool(false), SetPTR: lib.Bool(false)}
		rrset.Records = append(rrset.Records, r)
	}

	return r.patchRRset(domain, *rrset)
}

// Delete removes an existing resource record
func (r *RecordsService) Delete(domain string, name string, recordType lib.RRType) error {
	rrset := new(lib.RRset)
	rrset.Name = &name
	rrset.Type = &recordType
	rrset.ChangeType = lib.ChangeTypePtr(lib.ChangeTypeDelete)

	return r.patchRRset(domain, *rrset)
}

func canonicalResourceRecordValues(records []lib.Record) {
	for i := range records {
		records[i].Content = lib.String(lib.MakeDomainCanonical(*records[i].Content))
	}
}

func fixRRset(rrset *lib.RRset) {
	if *rrset.Type != lib.RRTypeCNAME && *rrset.Type != lib.RRTypeMX {
		return
	}

	canonicalResourceRecordValues(rrset.Records)
}

func (r *RecordsService) patchRRset(domain string, rrset lib.RRset) error {
	rrset.Name = lib.String(lib.MakeDomainCanonical(*rrset.Name))

	fixRRset(&rrset)

	payload := lib.RRsets{}
	payload.Sets = append(payload.Sets, rrset)

	req, err := r.client.newRequest("PATCH", fmt.Sprintf("servers/%s/zones/%s", r.client.VHost, lib.TrimDomain(domain)), nil, payload)
	if err != nil {
		return err
	}

	_, err = r.client.do(req, nil)

	return err
}
