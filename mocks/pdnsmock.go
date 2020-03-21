package mocks

import "github.com/jarcoal/httpmock"

type Mock struct {
	TestBaseURL string
	TestVHost   string
	TestAPIKey  string
}

func (m *Mock) Activate() {
	httpmock.Activate()
}

func (m *Mock) DeactivateAndReset() {
	httpmock.DeactivateAndReset()
}

func (m *Mock) Disabled() bool {
	return httpmock.Disabled()
}
