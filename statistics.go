package powerdns

import (
	"fmt"
	"net/url"

	"github.com/joeig/go-powerdns/v2/lib"
)

// StatisticsService handles communication with the statistics related methods of the Client API
type StatisticsService service

// List retrieves a list of Statistics
func (s *StatisticsService) List() ([]lib.Statistic, error) {
	req, err := s.client.newRequest("GET", fmt.Sprintf("servers/%s/statistics", s.client.VHost), nil, nil)
	if err != nil {
		return nil, err
	}

	statistics := make([]lib.Statistic, 0)
	_, err = s.client.do(req, &statistics)

	return statistics, err
}

// Get retrieves certain Statistics
func (s *StatisticsService) Get(statisticName string) ([]lib.Statistic, error) {
	query := url.Values{}
	query.Add("statistic", statisticName)

	req, err := s.client.newRequest("GET", fmt.Sprintf("servers/%s/statistics", s.client.VHost), &query, nil)
	if err != nil {
		return nil, err
	}

	statistics := make([]lib.Statistic, 0)
	_, err = s.client.do(req, &statistics)

	return statistics, err
}
