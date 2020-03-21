package powerdns

import (
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/joeig/go-powerdns/v2/lib"
)

func randomString(length int) string {
	bytes := make([]byte, length)

	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}

	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	for i, b := range bytes {
		character := b % byte(len(characters))
		bytes[i] = character
	}

	return characters
}

func generateTestRecord(client *Client, domain string, autoAddRecord bool) string {
	name := fmt.Sprintf("test-%s.%s", randomString(16), domain)

	if mock.Disabled() && autoAddRecord {
		if err := client.Records.Add(domain, name, lib.RRTypeTXT, 300, []string{"\"Testing...\""}); err != nil {
			fmt.Printf("Error creating record: %s\n", name)
			fmt.Printf("%s\n", err)
		} else {
			fmt.Printf("Created record %s\n", name)
		}
	}

	return name
}

func TestAddRecord(t *testing.T) {
	testDomain := generateTestZone(true)
	p := initialisePowerDNSTestClient(&mock)

	mock.RegisterRecordMockResponder(testDomain)

	testRecordNameTXT := generateTestRecord(p, testDomain, false)
	if err := p.Records.Add(testDomain, testRecordNameTXT, lib.RRTypeTXT, 300, []string{"\"bar\""}); err != nil {
		t.Errorf("%s", err)
	}

	testRecordNameCNAME := generateTestRecord(p, testDomain, false)
	if err := p.Records.Add(testDomain, testRecordNameCNAME, lib.RRTypeCNAME, 300, []string{"foo.tld"}); err != nil {
		t.Errorf("%s", err)
	}
}

func TestAddRecordError(t *testing.T) {
	p := initialisePowerDNSTestClient(&mock)
	p.Port = "x"
	testDomain := generateTestZone(false)

	testRecordName := generateTestRecord(p, testDomain, false)
	if err := p.Records.Add(testDomain, testRecordName, lib.RRTypeTXT, 300, []string{"\"bar\""}); err == nil {
		t.Error("error is nil")
	}
}

func TestChangeRecord(t *testing.T) {
	testDomain := generateTestZone(true)

	p := initialisePowerDNSTestClient(&mock)
	mock.RegisterRecordMockResponder(testDomain)

	testRecordName := generateTestRecord(p, testDomain, true)
	if err := p.Records.Change(testDomain, testRecordName, lib.RRTypeTXT, 300, []string{"\"bar\""}); err != nil {
		t.Errorf("%s", err)
	}
}

func TestChangeRecordError(t *testing.T) {
	p := initialisePowerDNSTestClient(&mock)
	p.Port = "x"
	testDomain := generateTestZone(false)

	testRecordName := generateTestRecord(p, testDomain, false)
	if err := p.Records.Change(testDomain, testRecordName, lib.RRTypeTXT, 300, []string{"\"bar\""}); err == nil {
		t.Error("error is nil")
	}
}

func TestDeleteRecord(t *testing.T) {
	testDomain := generateTestZone(true)
	p := initialisePowerDNSTestClient(&mock)
	mock.RegisterRecordMockResponder(testDomain)

	testRecordName := generateTestRecord(p, testDomain, true)
	if err := p.Records.Delete(testDomain, testRecordName, lib.RRTypeTXT); err != nil {
		t.Errorf("%s", err)
	}
}

func TestDeleteRecordError(t *testing.T) {
	p := initialisePowerDNSTestClient(&mock)
	p.Port = "x"
	testDomain := generateTestZone(false)

	testRecordName := generateTestRecord(p, testDomain, false)
	if err := p.Records.Delete(testDomain, testRecordName, lib.RRTypeTXT); err == nil {
		t.Error("error is nil")
	}
}
