package types

import (
	"fmt"
	"strings"
)

func TrimDomain(domain string) string {
	return strings.TrimSuffix(domain, ".")
}

func MakeDomainCanonical(domain string) string {
	return fmt.Sprintf("%s.", TrimDomain(domain))
}
