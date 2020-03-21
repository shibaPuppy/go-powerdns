package lib

// Zone structure with JSON API metadata
type Zone struct {
	ID               *string   `json:"id,omitempty"`
	Name             *string   `json:"name,omitempty"`
	Type             *ZoneType `json:"type,omitempty"`
	URL              *string   `json:"url,omitempty"`
	Kind             *ZoneKind `json:"kind,omitempty"`
	RRsets           []RRset   `json:"rrsets,omitempty"`
	Serial           *uint32   `json:"serial,omitempty"`
	NotifiedSerial   *uint32   `json:"notified_serial,omitempty"`
	Masters          []string  `json:"masters,omitempty"`
	DNSsec           *bool     `json:"dnssec,omitempty"`
	Nsec3Param       *string   `json:"nsec3param,omitempty"`
	Nsec3Narrow      *bool     `json:"nsec3narrow,omitempty"`
	Presigned        *bool     `json:"presigned,omitempty"`
	SOAEdit          *string   `json:"soa_edit,omitempty"`
	SOAEditAPI       *string   `json:"soa_edit_api,omitempty"`
	APIRectify       *bool     `json:"api_rectify,omitempty"`
	Zone             *string   `json:"zone,omitempty"`
	Account          *string   `json:"account,omitempty"`
	Nameservers      []string  `json:"nameservers,omitempty"`
	MasterTSIGKeyIDs []string  `json:"master_tsig_key_ids,omitempty"`
	SlaveTSIGKeyIDs  []string  `json:"slave_tsig_key_ids,omitempty"`
}

// NotifyResult structure with JSON API metadata
type NotifyResult struct {
	Result *string `json:"result,omitempty"`
}

// Export string type
type Export string

// ZoneType string type
type ZoneType string

// ZoneZoneType sets the zone's type to zone
const ZoneZoneType ZoneType = "Zone"

// ZoneTypePtr is a helper function that allocates a new ZoneType value to store v and returns a pointer to it.
func ZoneTypePtr(v ZoneType) *ZoneType {
	return &v
}

// ZoneKind string type
type ZoneKind string

// ZoneKindPtr is a helper function that allocates a new ZoneKind value to store v and returns a pointer to it.
func ZoneKindPtr(v ZoneKind) *ZoneKind {
	return &v
}

const (
	// NativeZoneKind sets the zone's kind to native
	NativeZoneKind ZoneKind = "Native"
	// MasterZoneKind sets the zone's kind to master
	MasterZoneKind ZoneKind = "Master"
	// SlaveZoneKind sets the zone's kind to slave
	SlaveZoneKind ZoneKind = "Slave"
)
