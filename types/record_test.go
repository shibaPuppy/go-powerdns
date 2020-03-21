package types

import "testing"

func TestChangeTypePtr(t *testing.T) {
	source := ChangeTypeReplace
	if *ChangeTypePtr(source) != source {
		t.Error("Invalid return value")
	}
}

func TestRRTypePtr(t *testing.T) {
	source := RRTypeDNSKEY
	if *RRTypePtr(source) != source {
		t.Error("Invalid return value")
	}
}
