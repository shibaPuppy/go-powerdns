package powerdns

import (
	"fmt"
	"io/ioutil"

	"github.com/joeig/go-powerdns/v2/lib"
)

// ZonesService handles communication with the zones related methods of the Client API
type ZonesService service

// List retrieves a list of Zones
func (z *ZonesService) List() ([]lib.Zone, error) {
	req, err := z.client.newRequest("GET", fmt.Sprintf("servers/%s/zones", z.client.VHost), nil, nil)
	if err != nil {
		return nil, err
	}

	zones := make([]lib.Zone, 0)
	_, err = z.client.do(req, &zones)

	return zones, err
}

// Get returns a certain Zone for a given domain
func (z *ZonesService) Get(domain string) (*lib.Zone, error) {
	req, err := z.client.newRequest("GET", fmt.Sprintf("servers/%s/zones/%s", z.client.VHost, lib.TrimDomain(domain)), nil, nil)
	if err != nil {
		return nil, err
	}

	zone := &lib.Zone{}
	_, err = z.client.do(req, &zone)

	return zone, err
}

// AddNative creates a new native zone
func (z *ZonesService) AddNative(domain string, dnssec bool, nsec3Param string, nsec3Narrow bool, soaEdit, soaEditAPI string, apiRectify bool, nameservers []string) (*lib.Zone, error) {
	zone := lib.Zone{
		Name:        lib.String(domain),
		Kind:        lib.ZoneKindPtr(lib.NativeZoneKind),
		DNSsec:      lib.Bool(dnssec),
		Nsec3Param:  lib.String(nsec3Param),
		Nsec3Narrow: lib.Bool(nsec3Narrow),
		SOAEdit:     lib.String(soaEdit),
		SOAEditAPI:  lib.String(soaEditAPI),
		APIRectify:  lib.Bool(apiRectify),
		Nameservers: nameservers,
	}

	return z.postZone(&zone)
}

// AddMaster creates a new master zone
func (z *ZonesService) AddMaster(domain string, dnssec bool, nsec3Param string, nsec3Narrow bool, soaEdit, soaEditAPI string, apiRectify bool, nameservers []string) (*lib.Zone, error) {
	zone := lib.Zone{
		Name:        lib.String(domain),
		Kind:        lib.ZoneKindPtr(lib.MasterZoneKind),
		DNSsec:      lib.Bool(dnssec),
		Nsec3Param:  lib.String(nsec3Param),
		Nsec3Narrow: lib.Bool(nsec3Narrow),
		SOAEdit:     lib.String(soaEdit),
		SOAEditAPI:  lib.String(soaEditAPI),
		APIRectify:  lib.Bool(apiRectify),
		Nameservers: nameservers,
	}

	return z.postZone(&zone)
}

// AddSlave creates a new slave zone
func (z *ZonesService) AddSlave(domain string, masters []string) (*lib.Zone, error) {
	zone := lib.Zone{
		Name:    lib.String(domain),
		Kind:    lib.ZoneKindPtr(lib.SlaveZoneKind),
		Masters: masters,
	}

	return z.postZone(&zone)
}

func (z *ZonesService) postZone(zone *lib.Zone) (*lib.Zone, error) {
	zone.Name = lib.String(lib.MakeDomainCanonical(*zone.Name))
	zone.Type = lib.ZoneTypePtr(lib.ZoneZoneType)

	req, err := z.client.newRequest("POST", fmt.Sprintf("servers/%s/zones", z.client.VHost), nil, zone)
	if err != nil {
		return nil, err
	}

	createdZone := new(lib.Zone)
	_, err = z.client.do(req, &createdZone)

	return createdZone, err
}

// Change modifies an existing zone
func (z *ZonesService) Change(domain string, zone *lib.Zone) error {
	zone.ID = nil
	zone.Name = nil
	zone.Type = nil
	zone.URL = nil

	req, err := z.client.newRequest("PUT", fmt.Sprintf("servers/%s/zones/%s", z.client.VHost, lib.TrimDomain(domain)), nil, zone)
	if err != nil {
		return err
	}

	_, err = z.client.do(req, nil)

	return err
}

// Delete removes a certain Zone for a given domain
func (z *ZonesService) Delete(domain string) error {
	req, err := z.client.newRequest("DELETE", fmt.Sprintf("servers/%s/zones/%s", z.client.VHost, lib.TrimDomain(domain)), nil, nil)
	if err != nil {
		return err
	}

	_, err = z.client.do(req, nil)

	return err
}

// Notify sends a DNS notify packet to all slaves
func (z *ZonesService) Notify(domain string) (*lib.NotifyResult, error) {
	req, err := z.client.newRequest("PUT", fmt.Sprintf("servers/%s/zones/%s/notify", z.client.VHost, lib.TrimDomain(domain)), nil, nil)
	if err != nil {
		return nil, err
	}

	notifyResult := &lib.NotifyResult{}
	_, err = z.client.do(req, notifyResult)

	return notifyResult, err
}

// Export returns a BIND-like Zone file
func (z *ZonesService) Export(domain string) (lib.Export, error) {
	req, err := z.client.newRequest("GET", fmt.Sprintf("servers/%s/zones/%s/export", z.client.VHost, lib.TrimDomain(domain)), nil, nil)
	if err != nil {
		return "", err
	}

	resp, err := z.client.do(req, nil)
	if err != nil {
		return "", err
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	return lib.Export(bodyBytes), nil
}
