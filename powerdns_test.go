package powerdns

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/joeig/go-powerdns/v2/lib"
	"github.com/joeig/go-powerdns/v2/mocks"
)

var mock mocks.Mock

func TestMain(m *testing.M) {
	mock.TestBaseURL = "http://localhost:8080"
	mock.TestVHost = "localhost"
	mock.TestAPIKey = "apipw"

	mock.Activate()
	defer mock.DeactivateAndReset()

	os.Exit(m.Run())
}

func initialisePowerDNSTestClient(m *mocks.Mock) *Client {
	return NewClient(m.TestBaseURL, m.TestVHost, map[string]string{"X-API-Key": m.TestAPIKey}, nil)
}

func TestNewClient(t *testing.T) {
	t.Run("TestValidURL", func(t *testing.T) {
		tmpl := &Client{"http", "localhost", "8080", "localhost", map[string]string{"X-API-Key": "apipw"}, http.DefaultClient, service{}, nil, nil, nil, nil, nil, nil}
		p := NewClient("http://localhost:8080", "localhost", map[string]string{"X-API-Key": "apipw"}, http.DefaultClient)
		if p.Hostname != tmpl.Hostname {
			t.Error("NewClient returns invalid Client object")
		}
	})

	t.Run("TestInvalidURL", func(t *testing.T) {
		originalLogFatalf := logFatalf
		defer func() {
			logFatalf = originalLogFatalf
		}()
		var errors []string
		logFatalf = func(format string, args ...interface{}) {
			if len(args) > 0 {
				errors = append(errors, fmt.Sprintf(format, args))
			} else {
				errors = append(errors, format)
			}
		}

		_ = NewClient("http://1.2:foo", "localhost", map[string]string{"X-API-Key": "apipw"}, http.DefaultClient)

		if len(errors) < 1 {
			t.Error("NewClient does not exit with fatal error")
		}
	})
}

func TestNewRequest(t *testing.T) {
	p := initialisePowerDNSTestClient(&mock)

	t.Run("TestValidRequest", func(t *testing.T) {
		if _, err := p.newRequest("GET", "servers", nil, nil); err != nil {
			t.Error("error is not nil")
		}
	})
}

func TestDo(t *testing.T) {
	mock.RegisterDoMockResponder()

	p := initialisePowerDNSTestClient(&mock)

	t.Run("TestStringErrorResponse", func(t *testing.T) {
		req, _ := p.newRequest("GET", "servers/doesntExist", nil, nil)
		if resp, err := p.do(req, nil); err == nil {
			defer func() {
				_ = resp.Body.Close()
			}()

			t.Error("err is nil")
		}
	})
	t.Run("Test401Handling", func(t *testing.T) {
		p.Headers = nil
		req, _ := p.newRequest("GET", "servers", nil, nil)
		if resp, err := p.do(req, nil); err == nil {
			defer func() {
				_ = resp.Body.Close()
			}()

			t.Error("401 response does not result into an error")
		}
	})
	t.Run("Test404Handling", func(t *testing.T) {
		req, _ := p.newRequest("GET", "servers/doesntExist", nil, nil)
		if resp, err := p.do(req, nil); err == nil {
			defer func() {
				_ = resp.Body.Close()
			}()

			t.Error("404 response does not result into an error")
		}
	})
	t.Run("TestJSONResponseHandling", func(t *testing.T) {
		req, _ := p.newRequest("GET", "server", nil, &lib.Server{})
		if resp, err := p.do(req, nil); err.(*lib.Error).Message != "Not Found" {
			defer func() {
				_ = resp.Body.Close()
			}()

			t.Error("501 JSON response does not result into Error structure")
		}
	})
}

func TestParseBaseURL(t *testing.T) {
	testCases := []struct {
		baseURL      string
		wantScheme   string
		wantHostname string
		wantPort     string
		wantError    bool
	}{
		{"https://example.com", "https", "example.com", "443", false},
		{"http://example.com", "http", "example.com", "80", false},
		{"https://example.com:8080", "https", "example.com", "8080", false},
		{"http://example.com:8080", "http", "example.com", "8080", false},
		{"http%%%foo", "http", "", "", true},
	}

	for i, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("TestCase%d", i), func(t *testing.T) {
			scheme, hostname, port, err := parseBaseURL(tc.baseURL)

			if err != nil && tc.wantError == true {
				return
			}
			if err != nil && tc.wantError == false {
				t.Error("Error was returned unexpectedly")
			}
			if err == nil && tc.wantError == true {
				t.Error("No error was returned")
			}
			if scheme != tc.wantScheme {
				t.Errorf("Scheme parsing failed: %s != %s", scheme, tc.wantScheme)
			}
			if hostname != tc.wantHostname {
				t.Errorf("Hostname parsing failed: %s != %s", hostname, tc.wantHostname)
			}
			if port != tc.wantPort {
				t.Errorf("Port parsing failed: %s != %s", port, tc.wantPort)
			}
		})
	}
}

func TestParseVHost(t *testing.T) {
	testCases := []struct {
		vHost     string
		wantVHost string
	}{
		{"example.com", "example.com"},
		{"", "localhost"},
	}

	for i, tc := range testCases {
		tc := tc

		t.Run(fmt.Sprintf("TestCase%d", i), func(t *testing.T) {
			if parseVHost(tc.vHost) != tc.wantVHost {
				t.Error("parseVHost returned an invalid value")
			}
		})
	}
}

func TestGenerateAPIURL(t *testing.T) {
	tmpl := "https://localhost:8080/api/v1/foo?a=b"
	query := url.Values{}
	query.Add("a", "b")
	g := generateAPIURL("https", "localhost", "8080", "foo", &query)

	if tmpl != g.String() {
		t.Errorf("Template does not match generated API URL: %s", g.String())
	}
}
