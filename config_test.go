package powerdns

import (
	"testing"
)

func TestListConfig(t *testing.T) {
	mock.RegisterConfigsMockResponder()

	p := initialisePowerDNSTestClient(&mock)
	config, err := p.Config.List()
	if err != nil {
		t.Errorf("%s", err)
	}
	if len(config) == 0 {
		t.Error("Received amount of config settings is 0")
	}
}

func TestListConfigError(t *testing.T) {
	p := initialisePowerDNSTestClient(&mock)
	p.Port = "x"
	if _, err := p.Config.List(); err == nil {
		t.Error("error is nil")
	}
}
