package powerdns

import (
	"fmt"
	"github.com/joeig/go-powerdns/v2/types"
	"net/url"
)

// ServersService handles communication with the servers related methods of the Client API
type ServersService service

// List retrieves a list of Servers
func (s *ServersService) List() ([]types.Server, error) {
	req, err := s.client.newRequest("GET", "servers", nil, nil)
	if err != nil {
		return nil, err
	}

	servers := make([]types.Server, 0)
	_, err = s.client.do(req, &servers)
	return servers, err
}

// Get returns a certain Server
func (s *ServersService) Get(vHost string) (*types.Server, error) {
	req, err := s.client.newRequest("GET", fmt.Sprintf("servers/%s", vHost), nil, nil)
	if err != nil {
		return nil, err
	}

	server := &types.Server{}
	_, err = s.client.do(req, &server)
	return server, err
}

// CacheFlush flushes a cache-entry by name
func (s *ServersService) CacheFlush(vHost string, domain string) (*types.CacheFlushResult, error) {
	query := url.Values{}
	query.Add("domain", types.MakeDomainCanonical(domain))
	req, err := s.client.newRequest("PUT", fmt.Sprintf("servers/%s/cache/flush", vHost), &query, nil)
	if err != nil {
		return nil, err
	}

	cacheFlushResult := &types.CacheFlushResult{}
	_, err = s.client.do(req, &cacheFlushResult)
	return cacheFlushResult, err
}
