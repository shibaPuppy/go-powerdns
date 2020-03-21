package powerdns

import (
	"fmt"

	"github.com/joeig/go-powerdns/v2/types"
)

// CryptokeysService handles communication with the cryptokeys related methods of the Client API
type CryptokeysService service

// List retrieves a list of Cryptokeys that belong to a Zone
func (c *CryptokeysService) List(domain string) ([]types.Cryptokey, error) {
	req, err := c.client.newRequest("GET", fmt.Sprintf("servers/%s/zones/%s/cryptokeys", c.client.VHost, types.TrimDomain(domain)), nil, nil)
	if err != nil {
		return nil, err
	}

	cryptokeys := make([]types.Cryptokey, 0)
	_, err = c.client.do(req, &cryptokeys)

	return cryptokeys, err
}

// Get returns a certain Cryptokey instance of a given Zone
func (c *CryptokeysService) Get(domain string, id uint64) (*types.Cryptokey, error) {
	req, err := c.client.newRequest("GET", fmt.Sprintf("servers/%s/zones/%s/cryptokeys/%s", c.client.VHost, types.TrimDomain(domain), types.CryptokeyIDToString(id)), nil, nil)
	if err != nil {
		return nil, err
	}

	cryptokey := new(types.Cryptokey)
	_, err = c.client.do(req, &cryptokey)

	return cryptokey, err
}

// Delete removes a given Cryptokey
func (c *CryptokeysService) Delete(domain string, id uint64) error {
	req, err := c.client.newRequest("DELETE", fmt.Sprintf("servers/%s/zones/%s/cryptokeys/%s", c.client.VHost, types.TrimDomain(domain), types.CryptokeyIDToString(id)), nil, nil)
	if err != nil {
		return err
	}

	_, err = c.client.do(req, nil)

	return err
}
