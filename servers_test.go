package powerdns

import (
	"testing"
)

func TestListServers(t *testing.T) {
	mock.RegisterServersMockResponder()

	p := initialisePowerDNSTestClient(&mock)

	servers, err := p.Servers.List()
	if err != nil {
		t.Errorf("%s", err)
	}

	if len(servers) == 0 {
		t.Error("Received amount of servers is 0")
	}
}

func TestListServersError(t *testing.T) {
	p := initialisePowerDNSTestClient(&mock)
	p.Port = "x"

	if _, err := p.Servers.List(); err == nil {
		t.Error("error is nil")
	}
}

func TestGetServer(t *testing.T) {
	mock.RegisterServersMockResponder()
	p := initialisePowerDNSTestClient(&mock)

	server, err := p.Servers.Get(mock.TestVHost)
	if err != nil {
		t.Errorf("%s", err)
	}

	if *server.ID != mock.TestVHost {
		t.Error("Received no server")
	}
}

func TestGetServerError(t *testing.T) {
	p := initialisePowerDNSTestClient(&mock)
	p.Port = "x"

	if _, err := p.Servers.Get(mock.TestVHost); err == nil {
		t.Error("error is nil")
	}
}

func TestCacheFlush(t *testing.T) {
	testDomain := generateTestZone(true)
	mock.RegisterCacheFlushMockResponder(testDomain)
	p := initialisePowerDNSTestClient(&mock)

	cacheFlushResult, err := p.Servers.CacheFlush(mock.TestVHost, testDomain)
	if err != nil {
		t.Errorf("%s", err)
	}

	if *cacheFlushResult.Count != 1 {
		t.Error("Received cache flush result is invalid")
	}
}

func TestCacheFlushResultError(t *testing.T) {
	testDomain := generateTestZone(false)
	p := initialisePowerDNSTestClient(&mock)
	p.Port = "x"

	if _, err := p.Servers.CacheFlush(mock.TestVHost, testDomain); err == nil {
		t.Error("error is nil")
	}
}
