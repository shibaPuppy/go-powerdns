package types

import "testing"

func TestZoneTypePtr(t *testing.T) {
	source := ZoneZoneType
	if *ZoneTypePtr(source) != source {
		t.Error("Invalid return value")
	}
}

func TestZoneKindPtr(t *testing.T) {
	source := NativeZoneKind
	if *ZoneKindPtr(source) != source {
		t.Error("Invalid return value")
	}
}
