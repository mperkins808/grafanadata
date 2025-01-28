package grafanadata

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// returns all the dashboards for a grafana instance
func (c *grafanaClient) FetchDashboards() ([]DashboardSearch, error) {
	host := strings.TrimSuffix(c.baseURL.String(), "/")

	q := fmt.Sprintf("%v/api/search?type=dash-db", host)
	req, err := c.NewRequest(http.MethodGet, q, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("grafana response %v. %v", resp.StatusCode, string(body))
	}

	var search []DashboardSearch
	err = json.Unmarshal(body, &search)
	return search, err
}

// returns all the panels for a dashboard
func (c *grafanaClient) FetchPanelsFromDashboard(dashboard GrafanaDashboardResponse) []PanelSearch {
	var search []PanelSearch
	for i := range dashboard.Dashboard.Panels {
		p := dashboard.Dashboard.Panels[i]
		search = append(search, PanelSearch{
			ID:    p.ID,
			Title: p.Title,
		})
	}
	return search
}
