package lib

// Statistic structure with JSON API metadata
type Statistic struct {
	Name *string `json:"name,omitempty"`
	Type *string `json:"type,omitempty"`

	// Contrary to the authoritative API specification, the "size" field has actually been implemented as string instead of integer.
	Size *string `json:"size,omitempty"`

	// The "value" field contains either a string or a list of objects, depending on the "type".
	Value interface{} `json:"value,omitempty"`
}
