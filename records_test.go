package powerdns

import (
	"fmt"
	"github.com/joeig/go-powerdns/v2/types"
	"math/rand"
	"testing"
	"time"
)

func generateTestRecord(client *Client, domain string, autoAddRecord bool) string {
	rand.Seed(time.Now().UnixNano())
	name := fmt.Sprintf("test-%d.%s", rand.Int(), domain)

	if mock.Disabled() && autoAddRecord {
		if err := client.Records.Add(domain, name, types.RRTypeTXT, 300, []string{"\"Testing...\""}); err != nil {
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
	if err := p.Records.Add(testDomain, testRecordNameTXT, types.RRTypeTXT, 300, []string{"\"bar\""}); err != nil {
		t.Errorf("%s", err)
	}
	testRecordNameCNAME := generateTestRecord(p, testDomain, false)
	if err := p.Records.Add(testDomain, testRecordNameCNAME, types.RRTypeCNAME, 300, []string{"foo.tld"}); err != nil {
		t.Errorf("%s", err)
	}
}

func TestAddRecordError(t *testing.T) {
	p := initialisePowerDNSTestClient(&mock)
	p.Port = "x"
	testDomain := generateTestZone(false)
	testRecordName := generateTestRecord(p, testDomain, false)
	if err := p.Records.Add(testDomain, testRecordName, types.RRTypeTXT, 300, []string{"\"bar\""}); err == nil {
		t.Error("error is nil")
	}
}

func TestChangeRecord(t *testing.T) {
	testDomain := generateTestZone(true)

	p := initialisePowerDNSTestClient(&mock)
	testRecordName := generateTestRecord(p, testDomain, true)
	mock.RegisterRecordMockResponder(testDomain)
	if err := p.Records.Change(testDomain, testRecordName, types.RRTypeTXT, 300, []string{"\"bar\""}); err != nil {
		t.Errorf("%s", err)
	}
}

func TestChangeRecordError(t *testing.T) {
	p := initialisePowerDNSTestClient(&mock)
	p.Port = "x"
	testDomain := generateTestZone(false)
	testRecordName := generateTestRecord(p, testDomain, false)
	if err := p.Records.Change(testDomain, testRecordName, types.RRTypeTXT, 300, []string{"\"bar\""}); err == nil {
		t.Error("error is nil")
	}
}

func TestDeleteRecord(t *testing.T) {
	testDomain := generateTestZone(true)

	p := initialisePowerDNSTestClient(&mock)
	testRecordName := generateTestRecord(p, testDomain, true)
	mock.RegisterRecordMockResponder(testDomain)
	if err := p.Records.Delete(testDomain, testRecordName, types.RRTypeTXT); err != nil {
		t.Errorf("%s", err)
	}
}

func TestDeleteRecordError(t *testing.T) {
	p := initialisePowerDNSTestClient(&mock)
	p.Port = "x"
	testDomain := generateTestZone(false)
	testRecordName := generateTestRecord(p, testDomain, false)
	if err := p.Records.Delete(testDomain, testRecordName, types.RRTypeTXT); err == nil {
		t.Error("error is nil")
	}
}

func TestCanonicalResourceRecordValues(t *testing.T) {
	testCases := []struct {
		records     []types.Record
		wantContent []string
	}{
		{[]types.Record{{Content: types.String("foo.tld")}}, []string{"foo.tld."}},
		{[]types.Record{{Content: types.String("foo.tld.")}}, []string{"foo.tld."}},
		{[]types.Record{{Content: types.String("foo.tld")}, {Content: types.String("foo.tld.")}}, []string{"foo.tld.", "foo.tld."}},
	}

	for i, tc := range testCases {
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
		rrset                     types.RRset
		wantFixedCanonicalRecords bool
	}{
		{types.RRset{Type: types.RRTypePtr(types.RRTypeMX), Records: []types.Record{{Content: types.String("foo.tld")}}}, true},
		{types.RRset{Type: types.RRTypePtr(types.RRTypeCNAME), Records: []types.Record{{Content: types.String("foo.tld")}}}, true},
		{types.RRset{Type: types.RRTypePtr(types.RRTypeA), Records: []types.Record{{Content: types.String("foo.tld")}}}, false},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TestCase%d", i), func(t *testing.T) {
			fixRRset(&tc.rrset)

			if tc.wantFixedCanonicalRecords {
				for j := range tc.rrset.Records {
					isContent := *tc.rrset.Records[j].Content
					wantContent := types.MakeDomainCanonical(*tc.rrset.Records[j].Content)
					if isContent != wantContent {
						t.Errorf("Comparison failed: %s != %s", isContent, wantContent)
					}
				}
			} else {
				for j := range tc.rrset.Records {
					isContent := *tc.rrset.Records[j].Content
					wrongContent := types.MakeDomainCanonical(*tc.rrset.Records[j].Content)
					if isContent == wrongContent {
						t.Errorf("Comparison failed: %s == %s", isContent, wrongContent)
					}
				}
			}
		})
	}
}
