package mocks

import (
	"net/http"

	"github.com/jarcoal/httpmock"
)

func (m *Mock) RegisterStatisticsMockResponder() {
	httpmock.RegisterResponder("GET", m.generateTestAPIVHostURL()+"/statistics",
		func(req *http.Request) (*http.Response, error) {
			if res := m.verifyAPIKey(req); res != nil {
				return res, nil
			}

			statisticsMock := "[{\"name\": \"corrupt-packets\", \"type\": \"StatisticItem\", \"value\": \"0\"}, {\"name\": \"response-by-rcode\", \"type\": \"MapStatisticItem\", \"value\": [{\"name\": \"foo1\", \"value\": \"bar1\"}, {\"name\": \"foo2\", \"value\": \"bar2\"}]}, {\"name\": \"logmessages\", \"size\": \"10000\", \"type\": \"RingStatisticItem\", \"value\": [{\"name\": \"gmysql Connection successful. Connected to database 'powerdns' on 'mariadb'.\", \"value\": \"235\"}]}]"

			statisticQueryString := req.URL.Query().Get("statistic")
			if statisticQueryString != "" {
				if statisticQueryString == "corrupt-packets" {
					statisticsMock = "[{\"name\": \"corrupt-packets\", \"type\": \"StatisticItem\", \"value\": \"0\"}]"
				} else {
					return httpmock.NewStringResponse(http.StatusUnprocessableEntity, "Unprocessable Entity"), nil
				}
			}

			return httpmock.NewStringResponse(http.StatusOK, statisticsMock), nil
		},
	)
}
