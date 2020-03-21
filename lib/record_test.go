package lib

import (
	"fmt"
	"testing"
)

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
		{[]Record{{Content: String("foo.tld")}}, []string{"foo.tld."}},
		{[]Record{{Content: String("foo.tld.")}}, []string{"foo.tld."}},
		{[]Record{{Content: String("foo.tld")}, {Content: String("foo.tld.")}}, []string{"foo.tld.", "foo.tld."}},
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
		{RRset{Type: RRTypePtr(RRTypeMX), Records: []Record{{Content: String("foo.tld")}}}, true},
		{RRset{Type: RRTypePtr(RRTypeCNAME), Records: []Record{{Content: String("foo.tld")}}}, true},
		{RRset{Type: RRTypePtr(RRTypeA), Records: []Record{{Content: String("foo.tld")}}}, false},
	}

	for i, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("TestCase%d", i), func(t *testing.T) {
			FixRRset(&tc.rrset)

			if tc.wantFixedCanonicalRecords {
				for j := range tc.rrset.Records {
					isContent := *tc.rrset.Records[j].Content
					wantContent := MakeDomainCanonical(*tc.rrset.Records[j].Content)

					if isContent != wantContent {
						t.Errorf("Comparison failed: %s != %s", isContent, wantContent)
					}
				}
			} else {
				for j := range tc.rrset.Records {
					isContent := *tc.rrset.Records[j].Content
					wrongContent := MakeDomainCanonical(*tc.rrset.Records[j].Content)

					if isContent == wrongContent {
						t.Errorf("Comparison failed: %s == %s", isContent, wrongContent)
					}
				}
			}
		})
	}
}
