package lib

import (
	"testing"
)

func TestConvertCryptokeyIDToString(t *testing.T) {
	if CryptokeyIDToString(1337) != "1337" {
		t.Error("Cryptokey ID to string conversion failed")
	}
}
