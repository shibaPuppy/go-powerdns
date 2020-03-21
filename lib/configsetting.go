package lib

// ConfigSetting structure with JSON API metadata
type ConfigSetting struct {
	Name  *string `json:"name,omitempty"`
	Type  *string `json:"type,omitempty"`
	Value *string `json:"value,omitempty"`
}
