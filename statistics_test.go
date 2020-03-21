package powerdns

import (
	"testing"
)

func TestListStatistics(t *testing.T) {
	mock.RegisterStatisticsMockResponder()
	p := initialisePowerDNSTestClient(&mock)

	statistics, err := p.Statistics.List()
	if err != nil {
		t.Errorf("%s", err)
	}

	if len(statistics) == 0 {
		t.Error("Received amount of statistics is 0")
	}
}

func TestListStatisticsError(t *testing.T) {
	p := initialisePowerDNSTestClient(&mock)
	p.Port = "x"

	if _, err := p.Statistics.List(); err == nil {
		t.Error("error is nil")
	}
}

func TestGetStatistics(t *testing.T) {
	mock.RegisterStatisticsMockResponder()
	p := initialisePowerDNSTestClient(&mock)

	statistics, err := p.Statistics.Get("corrupt-packets")
	if err != nil {
		t.Errorf("%s", err)
	}

	if len(statistics) != 1 {
		t.Error("Received amount of statistics is not 1")
	}
}

func TestGetStatisticsError(t *testing.T) {
	p := initialisePowerDNSTestClient(&mock)
	p.Port = "x"

	if _, err := p.Statistics.Get("corrupt-packets"); err == nil {
		t.Error("error is nil")
	}
}
