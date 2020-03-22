package lib

import "testing"

func TestBoolPtr(t *testing.T) {
	source := true

	if *BoolPtr(source) != source {
		t.Error("Invalid return value")
	}
}

func TestBoolValue(t *testing.T) {
	source := true

	if BoolValue(&source) != source {
		t.Error("Invalid return value")
	}

	if BoolValue(nil) != false {
		t.Error("Unexpected return value")
	}
}

func TestUint32Ptr(t *testing.T) {
	source := uint32(1337)

	if *Uint32Ptr(source) != source {
		t.Error("Invalid return value")
	}
}

func TestUint32Value(t *testing.T) {
	source := uint32(1337)

	if Uint32Value(&source) != source {
		t.Error("Invalid return value")
	}

	if Uint32Value(nil) != 0 {
		t.Error("Unexpected return value")
	}
}

func TestUint64Ptr(t *testing.T) {
	source := uint64(1337)

	if *Uint64Ptr(source) != source {
		t.Error("Invalid return value")
	}
}

func TestUint64Value(t *testing.T) {
	source := uint64(1337)

	if Uint64Value(&source) != source {
		t.Error("Invalid return value")
	}

	if Uint64Value(nil) != 0 {
		t.Error("Unexpected return value")
	}
}

func TestStringPtr(t *testing.T) {
	source := "foo"

	if *StringPtr(source) != source {
		t.Error("Invalid return value")
	}
}

func TestStringValue(t *testing.T) {
	source := "foo"

	if StringValue(&source) != source {
		t.Error("Invalid return value")
	}

	if StringValue(nil) != "" {
		t.Error("Unexpected return value")
	}
}

func TestStringSlicePtr(t *testing.T) {
	source := []string{"foo"}
	if (*StringSlicePtr(source))[0] != source[0] {
		t.Error("Invalid return value")
	}
}
