package lib

import (
	"fmt"
	"strings"
)

// TrimDomain removes a trailing dot from a domain
func TrimDomain(domain string) string {
	return strings.TrimSuffix(domain, ".")
}

// MakeDomainCanonical adds a trailing dot to a domain
func MakeDomainCanonical(domain string) string {
	return fmt.Sprintf("%s.", TrimDomain(domain))
}
