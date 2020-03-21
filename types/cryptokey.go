package types

import "strconv"

// Cryptokey structure with JSON API metadata
type Cryptokey struct {
	Type       *string  `json:"type,omitempty"`
	ID         *uint64  `json:"id,omitempty"`
	KeyType    *string  `json:"keytype,omitempty"`
	Active     *bool    `json:"active,omitempty"`
	DNSkey     *string  `json:"dnskey,omitempty"`
	DS         []string `json:"ds,omitempty"`
	Privatekey *string  `json:"privatekey,omitempty"`
	Algorithm  *string  `json:"algorithm,omitempty"`
	Bits       *uint64  `json:"bits,omitempty"`
}

func CryptokeyIDToString(id uint64) string {
	return strconv.FormatUint(id, 10)
}
