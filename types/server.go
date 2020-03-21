package types

// Server structure with JSON API metadata
type Server struct {
	Type       *string `json:"type,omitempty"`
	ID         *string `json:"id,omitempty"`
	DaemonType *string `json:"daemon_type,omitempty"`
	Version    *string `json:"version,omitempty"`
	URL        *string `json:"url,omitempty"`
	ConfigURL  *string `json:"config_url,omitempty"`
	ZonesURL   *string `json:"zones_url,omitempty"`
}

// CacheFlushResult structure with JSON API metadata
type CacheFlushResult struct {
	Count  *uint32 `json:"count,omitempty"`
	Result *string `json:"result,omitempty"`
}
