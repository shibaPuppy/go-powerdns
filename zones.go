package powerdns

import (
	"fmt"
	"github.com/joeig/go-powerdns/v2/types"
	"io/ioutil"
)

// ZonesService handles communication with the zones related methods of the Client API
type ZonesService service

// List retrieves a list of Zones
func (z *ZonesService) List() ([]types.Zone, error) {
	req, err := z.client.newRequest("GET", fmt.Sprintf("servers/%s/zones", z.client.VHost), nil, nil)
	if err != nil {
		return nil, err
	}

	zones := make([]types.Zone, 0)
	_, err = z.client.do(req, &zones)
	return zones, err
}

// Get returns a certain Zone for a given domain
func (z *ZonesService) Get(domain string) (*types.Zone, error) {
	req, err := z.client.newRequest("GET", fmt.Sprintf("servers/%s/zones/%s", z.client.VHost, types.TrimDomain(domain)), nil, nil)
	if err != nil {
		return nil, err
	}

	zone := &types.Zone{}
	_, err = z.client.do(req, &zone)
	return zone, err
}

// AddNative creates a new native zone
func (z *ZonesService) AddNative(domain string, dnssec bool, nsec3Param string, nsec3Narrow bool, soaEdit, soaEditApi string, apiRectify bool, nameservers []string) (*types.Zone, error) {
	zone := types.Zone{
		Name:        types.String(domain),
		Kind:        types.ZoneKindPtr(types.NativeZoneKind),
		DNSsec:      types.Bool(dnssec),
		Nsec3Param:  types.String(nsec3Param),
		Nsec3Narrow: types.Bool(nsec3Narrow),
		SOAEdit:     types.String(soaEdit),
		SOAEditAPI:  types.String(soaEditApi),
		APIRectify:  types.Bool(apiRectify),
		Nameservers: nameservers,
	}
	return z.postZone(&zone)
}

// AddMaster creates a new master zone
func (z *ZonesService) AddMaster(domain string, dnssec bool, nsec3Param string, nsec3Narrow bool, soaEdit, soaEditApi string, apiRectify bool, nameservers []string) (*types.Zone, error) {
	zone := types.Zone{
		Name:        types.String(domain),
		Kind:        types.ZoneKindPtr(types.MasterZoneKind),
		DNSsec:      types.Bool(dnssec),
		Nsec3Param:  types.String(nsec3Param),
		Nsec3Narrow: types.Bool(nsec3Narrow),
		SOAEdit:     types.String(soaEdit),
		SOAEditAPI:  types.String(soaEditApi),
		APIRectify:  types.Bool(apiRectify),
		Nameservers: nameservers,
	}
	return z.postZone(&zone)
}

// AddSlave creates a new slave zone
func (z *ZonesService) AddSlave(domain string, masters []string) (*types.Zone, error) {
	zone := types.Zone{
		Name:    types.String(domain),
		Kind:    types.ZoneKindPtr(types.SlaveZoneKind),
		Masters: masters,
	}
	return z.postZone(&zone)
}

func (z *ZonesService) postZone(zone *types.Zone) (*types.Zone, error) {
	zone.Name = types.String(types.MakeDomainCanonical(*zone.Name))
	zone.Type = types.ZoneTypePtr(types.ZoneZoneType)

	req, err := z.client.newRequest("POST", fmt.Sprintf("servers/%s/zones", z.client.VHost), nil, zone)
	if err != nil {
		return nil, err
	}

	createdZone := new(types.Zone)
	_, err = z.client.do(req, &createdZone)
	return createdZone, err
}

// Change modifies an existing zone
func (z *ZonesService) Change(domain string, zone *types.Zone) error {
	zone.ID = nil
	zone.Name = nil
	zone.Type = nil
	zone.URL = nil

	req, err := z.client.newRequest("PUT", fmt.Sprintf("servers/%s/zones/%s", z.client.VHost, types.TrimDomain(domain)), nil, zone)
	if err != nil {
		return err
	}

	_, err = z.client.do(req, nil)
	return err
}

// Delete removes a certain Zone for a given domain
func (z *ZonesService) Delete(domain string) error {
	req, err := z.client.newRequest("DELETE", fmt.Sprintf("servers/%s/zones/%s", z.client.VHost, types.TrimDomain(domain)), nil, nil)
	if err != nil {
		return err
	}

	_, err = z.client.do(req, nil)
	return err
}

// Notify sends a DNS notify packet to all slaves
func (z *ZonesService) Notify(domain string) (*types.NotifyResult, error) {
	req, err := z.client.newRequest("PUT", fmt.Sprintf("servers/%s/zones/%s/notify", z.client.VHost, types.TrimDomain(domain)), nil, nil)
	if err != nil {
		return nil, err
	}

	notifyResult := &types.NotifyResult{}
	_, err = z.client.do(req, notifyResult)
	return notifyResult, err
}

// Export returns a BIND-like Zone file
func (z *ZonesService) Export(domain string) (types.Export, error) {
	req, err := z.client.newRequest("GET", fmt.Sprintf("servers/%s/zones/%s/export", z.client.VHost, types.TrimDomain(domain)), nil, nil)
	if err != nil {
		return "", err
	}

	resp, err := z.client.do(req, nil)
	if err != nil {
		return "", err
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	return types.Export(bodyBytes), nil
}
