package lib

// RRset structure with JSON API metadata
type RRset struct {
	Name       *string     `json:"name,omitempty"`
	Type       *RRType     `json:"type,omitempty"`
	TTL        *uint32     `json:"ttl,omitempty"`
	ChangeType *ChangeType `json:"changetype,omitempty"`
	Records    *[]Record   `json:"records"`
	Comments   *[]Comment  `json:"comments,omitempty"`
}

// Record structure with JSON API metadata
type Record struct {
	Content  *string `json:"content,omitempty"`
	Disabled *bool   `json:"disabled,omitempty"`
	SetPTR   *bool   `json:"set-ptr,omitempty"`
}

// RecordSlicePtr is a helper function that allocates a new record slice to store v and returns a pointer to it.
func RecordSlicePtr(v []Record) *[]Record {
	return &v
}

// Comment structure with JSON API metadata
type Comment struct {
	Content    *string `json:"content,omitempty"`
	Account    *string `json:"account,omitempty"`
	ModifiedAt *uint64 `json:"modified_at,omitempty"`
}

// RRsets structure with JSON API metadata
type RRsets struct {
	Sets *[]RRset `json:"rrsets,omitempty"`
}

// RRsetSlicePtr is a helper function that allocates a new RRset slice to store v and returns a pointer to it.
func RRsetSlicePtr(v []RRset) *[]RRset {
	return &v
}

// ChangeType represents a string-valued change type
type ChangeType string

// ChangeTypePtr is a helper function that allocates a new ChangeType value to store v and returns a pointer to it.
func ChangeTypePtr(v ChangeType) *ChangeType {
	return &v
}

const (
	// ChangeTypeReplace represents the REPLACE change type
	ChangeTypeReplace ChangeType = "REPLACE"
	// ChangeTypeDelete represents the DELETE change type
	ChangeTypeDelete ChangeType = "DELETE"
)

// RRType represents a string-valued resource record type
type RRType string

// RRTypePtr is a helper function that allocates a new RRType value to store v and returns a pointer to it.
func RRTypePtr(v RRType) *RRType {
	return &v
}

const (
	// RRTypeA represents the A resource record type
	RRTypeA RRType = "A"
	// RRTypeAAAA represents the AAAA resource record type
	RRTypeAAAA RRType = "AAAA"
	// RRTypeAFSDB represents the AFSDB resource record type
	RRTypeAFSDB RRType = "AFSDB"
	// RRTypeALIAS represents the ALIAS resource record type
	RRTypeALIAS RRType = "ALIAS"
	// RRTypeCAA represents the CAA resource record type
	RRTypeCAA RRType = "CAA"
	// RRTypeCERT represents the CERT resource record type
	RRTypeCERT RRType = "CERT"
	// RRTypeCDNSKEY represents the CDNSKEY resource record type
	RRTypeCDNSKEY RRType = "CDNSKEY"
	// RRTypeCDS represents the CDS resource record type
	RRTypeCDS RRType = "CDS"
	// RRTypeCNAME represents the CNAME resource record type
	RRTypeCNAME RRType = "CNAME"
	// RRTypeDNSKEY represents the DNSKEY resource record type
	RRTypeDNSKEY RRType = "DNSKEY"
	// RRTypeDNAME represents the DNAME resource record type
	RRTypeDNAME RRType = "DNAME"
	// RRTypeDS represents the DS resource record type
	RRTypeDS RRType = "DS"
	// RRTypeHINFO represents the HINFO resource record type
	RRTypeHINFO RRType = "HINFO"
	// RRTypeKEY represents the KEY resource record type
	RRTypeKEY RRType = "KEY"
	// RRTypeLOC represents the LOC resource record type
	RRTypeLOC RRType = "LOC"
	// RRTypeMX represents the MX resource record type
	RRTypeMX RRType = "MX"
	// RRTypeNAPTR represents the NAPTR resource record type
	RRTypeNAPTR RRType = "NAPTR"
	// RRTypeNS represents the NS resource record type
	RRTypeNS RRType = "NS"
	// RRTypeNSEC represents the NSEC resource record type
	RRTypeNSEC RRType = "NSEC"
	// RRTypeNSEC3 represents the NSEC3 resource record type
	RRTypeNSEC3 RRType = "NSEC3"
	// RRTypeNSEC3PARAM represents the NSEC3PARAM resource record type
	RRTypeNSEC3PARAM RRType = "NSEC3PARAM"
	// RRTypeOPENPGPKEY represents the OPENPGPKEY resource record type
	RRTypeOPENPGPKEY RRType = "OPENPGPKEY"
	// RRTypePTR represents the PTR resource record type
	RRTypePTR RRType = "PTR"
	// RRTypeRP represents the RP resource record type
	RRTypeRP RRType = "RP"
	// RRTypeRRSIG represents the RRSIG resource record type
	RRTypeRRSIG RRType = "RRSIG"
	// RRTypeSOA represents the SOA resource record type
	RRTypeSOA RRType = "SOA"
	// RRTypeSPF represents the SPF resource record type
	RRTypeSPF RRType = "SPF"
	// RRTypeSSHFP represents the SSHFP resource record type
	RRTypeSSHFP RRType = "SSHFP"
	// RRTypeSRV represents the SRV resource record type
	RRTypeSRV RRType = "SRV"
	// RRTypeTKEY represents the TKEY resource record type
	RRTypeTKEY RRType = "TKEY"
	// RRTypeTSIG represents the TSIG resource record type
	RRTypeTSIG RRType = "TSIG"
	// RRTypeTLSA represents the TLSA resource record type
	RRTypeTLSA RRType = "TLSA"
	// RRTypeSMIMEA represents the SMIMEA resource record type
	RRTypeSMIMEA RRType = "SMIMEA"
	// RRTypeTXT represents the TXT resource record type
	RRTypeTXT RRType = "TXT"
	// RRTypeURI represents the URI resource record type
	RRTypeURI RRType = "URI"
)

func canonicalResourceRecordValues(records []Record) {
	for i := range records {
		records[i].Content = StringPtr(MakeDomainCanonical(*records[i].Content))
	}
}

// FixRRset fixes a given RRset in order to comply with PowerDNS requirements
func FixRRset(rrset *RRset) {
	if *rrset.Type != RRTypeCNAME && *rrset.Type != RRTypeMX {
		return
	}

	canonicalResourceRecordValues(*rrset.Records)
}
