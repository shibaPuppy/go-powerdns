package powerdns

import (
	"testing"
)

func TestListCryptokeys(t *testing.T) {
	testDomain := generateTestZone(true)
	mock.RegisterCryptokeysMockResponder(testDomain)

	p := initialisePowerDNSTestClient(&mock)

	cryptokeys, err := p.Cryptokeys.List(testDomain)
	if err != nil {
		t.Errorf("%s", err)
	}

	if len(cryptokeys) == 0 {
		t.Error("Received amount of statistics is 0")
	}
}

func TestListCryptokeysError(t *testing.T) {
	testDomain := generateTestZone(false)
	p := initialisePowerDNSTestClient(&mock)
	p.Port = "x"

	if _, err := p.Cryptokeys.List(testDomain); err == nil {
		t.Error("error is nil")
	}
}

func TestGetCryptokey(t *testing.T) {
	testDomain := generateTestZone(true)
	p := initialisePowerDNSTestClient(&mock)

	mock.RegisterCryptokeysMockResponder(testDomain)

	cryptokeys, err := p.Cryptokeys.List(testDomain)
	if err != nil {
		t.Errorf("%s", err)
	}

	id := cryptokeys[0].ID

	mock.RegisterCryptokeyMockResponders(testDomain, *id)

	cryptokey, err := p.Cryptokeys.Get(testDomain, *id)
	if err != nil {
		t.Errorf("%s", err)
	}

	if *cryptokey.Algorithm != "ECDSAP256SHA256" {
		t.Error("Received cryptokey algorithm is wrong")
	}
}

func TestGetCryptokeyError(t *testing.T) {
	testDomain := generateTestZone(false)
	p := initialisePowerDNSTestClient(&mock)
	p.Port = "x"

	if _, err := p.Cryptokeys.Get(testDomain, uint64(0)); err == nil {
		t.Error("error is nil")
	}
}

func TestDeleteCryptokey(t *testing.T) {
	testDomain := generateTestZone(true)
	p := initialisePowerDNSTestClient(&mock)

	mock.RegisterCryptokeysMockResponder(testDomain)

	cryptokeys, err := p.Cryptokeys.List(testDomain)
	if err != nil {
		t.Errorf("%s", err)
	}

	id := cryptokeys[0].ID
	mock.RegisterCryptokeyMockResponders(testDomain, *id)

	if p.Cryptokeys.Delete(testDomain, *id) != nil {
		t.Errorf("%s", err)
	}
}

func TestDeleteCryptokeyError(t *testing.T) {
	testDomain := generateTestZone(false)
	p := initialisePowerDNSTestClient(&mock)
	p.Port = "x"

	if err := p.Cryptokeys.Delete(testDomain, uint64(0)); err == nil {
		t.Error("error is nil")
	}
}
