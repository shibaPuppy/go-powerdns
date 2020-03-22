package lib

import (
	"fmt"
	"testing"
)

func TestRecordSlicePtr(t *testing.T) {
	source := []Record{
		{
			Content: StringPtr("foo"),
		},
	}

	if (*RecordSlicePtr(source))[0].Content != source[0].Content {
		t.Error("Invalid return value")
	}
}

func TestRRsetSlicePtr(t *testing.T) {
	source := []RRset{
		{
			Name: StringPtr("foo"),
		},
	}

	if (*RRsetSlicePtr(source))[0].Name != source[0].Name {
		t.Error("Invalid return value")
	}
}

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

func TestCanonicalResourceRecordValues(t *testing.T) {
	testCases := []struct {
		records     []Record
		wantContent []string
	}{
		{[]Record{{Content: StringPtr("foo.tld")}}, []string{"foo.tld."}},
		{[]Record{{Content: StringPtr("foo.tld.")}}, []string{"foo.tld."}},
		{[]Record{{Content: StringPtr("foo.tld")}, {Content: StringPtr("foo.tld.")}}, []string{"foo.tld.", "foo.tld."}},
	}

	for i, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("TestCase%d", i), func(t *testing.T) {
			canonicalResourceRecordValues(tc.records)

			for j := range tc.records {
				isContent := *tc.records[j].Content
				wantContent := tc.wantContent[j]

				if isContent != wantContent {
					t.Errorf("Comparison failed: %s != %s", isContent, wantContent)
				}
			}
		})
	}
}

func TestFixRRset(t *testing.T) {
	testCases := []struct {
		rrset                     RRset
		wantFixedCanonicalRecords bool
	}{
		{RRset{Type: RRTypePtr(RRTypeMX), Records: RecordSlicePtr([]Record{{Content: StringPtr("foo.tld")}})}, true},
		{RRset{Type: RRTypePtr(RRTypeCNAME), Records: RecordSlicePtr([]Record{{Content: StringPtr("foo.tld")}})}, true},
		{RRset{Type: RRTypePtr(RRTypeA), Records: RecordSlicePtr([]Record{{Content: StringPtr("foo.tld")}})}, false},
	}

	for i, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("TestCase%d", i), func(t *testing.T) {
			FixRRset(&tc.rrset)

			if tc.wantFixedCanonicalRecords {
				for _, r := range *tc.rrset.Records {
					isContent := *r.Content
					wantContent := MakeDomainCanonical(*r.Content)

					if isContent != wantContent {
						t.Errorf("Comparison failed: %s != %s", isContent, wantContent)
					}
				}
			} else {
				for _, r := range *tc.rrset.Records {
					isContent := *r.Content
					wrongContent := MakeDomainCanonical(*r.Content)

					if isContent == wrongContent {
						t.Errorf("Comparison failed: %s == %s", isContent, wrongContent)
					}
				}
			}
		})
	}
}
