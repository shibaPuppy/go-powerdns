package powerdns

import (
	"fmt"
	"strings"
	"testing"

	"github.com/joeig/go-powerdns/v2/lib"
)

func generateTestZone(autoAddZone bool) string {
	domain := fmt.Sprintf("test-%s.com", randomString(16))

	if mock.Disabled() && autoAddZone {
		p := initialisePowerDNSTestClient(&mock)

		zone, err := p.Zones.AddNative(domain, true, "", false, "", "", true, []string{"ns.foo.tld."})
		if err != nil {
			fmt.Printf("Error creating %s\n", domain)
			fmt.Printf("%v\n", err)
			fmt.Printf("%v\n", zone)
		} else {
			fmt.Printf("Created domain %s\n", domain)
		}
	}

	return domain
}

func TestListZones(t *testing.T) {
	mock.RegisterZonesMockResponder()
	p := initialisePowerDNSTestClient(&mock)

	zones, err := p.Zones.List()
	if err != nil {
		t.Errorf("%s", err)
	}

	if len(zones) == 0 {
		t.Error("Received amount of statistics is 0")
	}
}

func TestListZonesError(t *testing.T) {
	p := initialisePowerDNSTestClient(&mock)
	p.Port = "x"

	if _, err := p.Zones.List(); err == nil {
		t.Error("error is nil")
	}
}

func TestGetZone(t *testing.T) {
	testDomain := generateTestZone(true)
	mock.RegisterZoneMockResponders(testDomain, lib.NativeZoneKind)
	p := initialisePowerDNSTestClient(&mock)

	zone, err := p.Zones.Get(testDomain)
	if err != nil {
		t.Errorf("%s", err)
	}

	if *zone.ID != lib.MakeDomainCanonical(testDomain) {
		t.Error("Received no zone")
	}
}

func TestGetZonesError(t *testing.T) {
	testDomain := generateTestZone(false)
	p := initialisePowerDNSTestClient(&mock)
	p.Port = "x"

	if _, err := p.Zones.Get(testDomain); err == nil {
		t.Error("error is nil")
	}
}

func TestAddNativeZone(t *testing.T) {
	testDomain := generateTestZone(false)
	mock.RegisterZoneMockResponders(testDomain, lib.NativeZoneKind)
	p := initialisePowerDNSTestClient(&mock)

	zone, err := p.Zones.AddNative(testDomain, true, "", false, "foo", "foo", true, []string{"ns.foo.tld."})
	if err != nil {
		t.Errorf("%s", err)
	}

	if *zone.ID != lib.MakeDomainCanonical(testDomain) || *zone.Kind != lib.NativeZoneKind {
		t.Error("Zone wasn't created")
	}
}

func TestAddNativeZoneError(t *testing.T) {
	testDomain := generateTestZone(false)
	p := initialisePowerDNSTestClient(&mock)
	p.Port = "x"

	if _, err := p.Zones.AddNative(testDomain, true, "", false, "foo", "foo", true, []string{"ns.foo.tld."}); err == nil {
		t.Error("error is nil")
	}
}

func TestAddMasterZone(t *testing.T) {
	testDomain := generateTestZone(false)
	mock.RegisterZoneMockResponders(testDomain, lib.MasterZoneKind)
	p := initialisePowerDNSTestClient(&mock)

	zone, err := p.Zones.AddMaster(testDomain, true, "", false, "foo", "foo", true, []string{"ns.foo.tld."})
	if err != nil {
		t.Errorf("%s", err)
	}

	if *zone.ID != lib.MakeDomainCanonical(testDomain) || *zone.Kind != lib.MasterZoneKind {
		t.Error("Zone wasn't created")
	}
}

func TestAddMasterZoneError(t *testing.T) {
	testDomain := generateTestZone(false)
	p := initialisePowerDNSTestClient(&mock)
	p.Port = "x"

	if _, err := p.Zones.AddMaster(testDomain, true, "", false, "foo", "foo", true, []string{"ns.foo.tld."}); err == nil {
		t.Error("error is nil")
	}
}

func TestAddSlaveZone(t *testing.T) {
	testDomain := generateTestZone(false)
	mock.RegisterZoneMockResponders(testDomain, lib.SlaveZoneKind)
	p := initialisePowerDNSTestClient(&mock)

	zone, err := p.Zones.AddSlave(testDomain, []string{"127.0.0.1"})
	if err != nil {
		t.Errorf("%s", err)
	}

	if *zone.ID != lib.MakeDomainCanonical(testDomain) || *zone.Kind != lib.SlaveZoneKind {
		t.Error("Zone wasn't created")
	}
}

func TestAddSlaveZoneError(t *testing.T) {
	testDomain := generateTestZone(false)
	p := initialisePowerDNSTestClient(&mock)
	p.Port = "x"

	if _, err := p.Zones.AddSlave(testDomain, []string{"ns5.foo.tld."}); err == nil {
		t.Error("error is nil")
	}
}

func TestChangeZone(t *testing.T) {
	testDomain := generateTestZone(true)

	mock.RegisterZoneMockResponders(testDomain, lib.NativeZoneKind)

	p := initialisePowerDNSTestClient(&mock)

	t.Run("ChangeValidZone", func(t *testing.T) {
		if err := p.Zones.Change(testDomain, &lib.Zone{Nameservers: lib.StringSlicePtr([]string{"ns23.foo.tld."})}); err != nil {
			t.Errorf("%s", err)
		}
	})
	t.Run("ChangeInvalidZone", func(t *testing.T) {
		if err := p.Zones.Change("doesnt-exist", &lib.Zone{Nameservers: lib.StringSlicePtr([]string{"ns23.foo.tld."})}); err == nil {
			t.Errorf("Changing an invalid zone does not return an error")
		}
	})
}

func TestChangeZoneError(t *testing.T) {
	testDomain := generateTestZone(false)
	p := initialisePowerDNSTestClient(&mock)
	p.Port = "x"

	if err := p.Zones.Change(testDomain, &lib.Zone{Nameservers: lib.StringSlicePtr([]string{"ns23.foo.tld."})}); err == nil {
		t.Error("error is nil")
	}
}

func TestDeleteZone(t *testing.T) {
	testDomain := generateTestZone(true)
	mock.RegisterZoneMockResponders(testDomain, lib.NativeZoneKind)
	p := initialisePowerDNSTestClient(&mock)

	if err := p.Zones.Delete(testDomain); err != nil {
		t.Errorf("%s", err)
	}
}

func TestDeleteZoneError(t *testing.T) {
	testDomain := generateTestZone(false)
	p := initialisePowerDNSTestClient(&mock)
	p.Port = "x"

	if err := p.Zones.Delete(testDomain); err == nil {
		t.Error("error is nil")
	}
}

func TestNotify(t *testing.T) {
	testDomain := generateTestZone(true)
	mock.RegisterZoneMockResponders(testDomain, lib.MasterZoneKind)
	p := initialisePowerDNSTestClient(&mock)

	notifyResult, err := p.Zones.Notify(testDomain)
	if err != nil {
		t.Errorf("%s", err)
	}

	if *notifyResult.Result != "Notification queued" {
		t.Error("Notification was not queued successfully")
	}
}

func TestNotifyError(t *testing.T) {
	testDomain := generateTestZone(false)
	p := initialisePowerDNSTestClient(&mock)
	p.Port = "x"

	if _, err := p.Zones.Notify(testDomain); err == nil {
		t.Error("error is nil")
	}
}

func TestExport(t *testing.T) {
	testDomain := generateTestZone(true)
	mock.RegisterZoneMockResponders(testDomain, lib.NativeZoneKind)
	p := initialisePowerDNSTestClient(&mock)

	export, err := p.Zones.Export(testDomain)
	if err != nil {
		t.Errorf("%s", err)
	}

	if !strings.HasPrefix(string(export), testDomain) {
		t.Errorf("Export payload wrong: \"%s\"", export)
	}
}

func TestExportError(t *testing.T) {
	testDomain := generateTestZone(false)
	p := initialisePowerDNSTestClient(&mock)
	p.Hostname = "invalid"

	if _, err := p.Zones.Export(testDomain); err == nil {
		t.Error("error is nil")
	}

	p.Port = "x"

	if _, err := p.Zones.Export(testDomain); err == nil {
		t.Error("error is nil")
	}
}
