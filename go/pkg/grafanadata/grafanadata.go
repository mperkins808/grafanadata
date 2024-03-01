package grafanadata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// GrafanaClient interface defines the methods that our client will implement.
type GrafanaClient interface {
	SetToken(token string)
	NewRequest(method, endpoint string, body io.Reader) (*http.Request, error)
	Do(req *http.Request) (*http.Response, error)
	GetDashboardFromURL(urlstr string) (GrafanaDashboardResponse, error)
	GetDashboardWithUID(uid string) (GrafanaDashboardResponse, error)
	GetPanelDataFromURL(urlstr string, start time.Time) (Results, error)
	GetPanelDataFromID(uid string, panelID int, start time.Time) (Results, error)
	FetchDashboards() ([]DashboardSearch, error)
	FetchPanelsFromDashboard(dashboard GrafanaDashboardResponse) []PanelSearch
}

// needed for unit tests
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type grafanaClient struct {
	baseURL *url.URL
	token   string
	client  HTTPClient
}

// NewGrafanaClient creates a new Grafana Client with an API token and returns the GrafanaClient interface
func NewGrafanaClient(urlstr string, token string) (GrafanaClient, error) {
	parsed, err := url.Parse(urlstr)
	if err != nil {
		return nil, fmt.Errorf("failed to create GrafanaClient. %v", err)
	}

	return &grafanaClient{
		baseURL: parsed,
		token:   token,
		client:  &http.Client{},
	}, nil
}

// NewCustomGrafanaClient creates a new Grafana Client with your own custom http Client
func NewCustomGrafanaClient(c *http.Client, urlstr string, token string) (GrafanaClient, error) {
	parsed, err := url.Parse(urlstr)
	if err != nil {
		return nil, fmt.Errorf("failed to create Custom GrafanaClient. %v", err)
	}

	return &grafanaClient{
		baseURL: parsed,
		token:   token,
		client:  c,
	}, nil
}

// SetToken sets the API token for the client.
func (c *grafanaClient) SetToken(token string) {
	c.token = token
}

func (c *grafanaClient) getDashboard(uid string) (GrafanaDashboardResponse, error) {
	var grafanaDashboardResponse GrafanaDashboardResponse

	host := fmt.Sprintf("%v://%v", c.baseURL.Scheme, c.baseURL.Host)
	query := fmt.Sprintf("%v/api/dashboards/uid/%v", host, uid)

	req, err := c.NewRequest(http.MethodGet, query, nil)
	if err != nil {
		return grafanaDashboardResponse, fmt.Errorf("failed to get dashboard %v with error %v", uid, err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return grafanaDashboardResponse, fmt.Errorf("failed to get dashboard %v with error %v", uid, err)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return grafanaDashboardResponse, fmt.Errorf("could not read response body with error %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return grafanaDashboardResponse, fmt.Errorf("grafana returned status %v", resp.StatusCode)
	}

	err = json.Unmarshal(b, &grafanaDashboardResponse)
	if err != nil {
		return grafanaDashboardResponse, fmt.Errorf("could not unmarshal response %v", err)
	}

	return grafanaDashboardResponse, nil
}

// retrieves the data for a panel in a dashboard.
func (c *grafanaClient) getPanelData(panelID int, dashboard GrafanaDashboardResponse, start time.Time) (Results, error) {
	var result Results

	var targets []interface{}
	for i := range dashboard.Dashboard.Panels {
		p := dashboard.Dashboard.Panels[i]
		if p.ID != panelID {
			continue
		}
		targets = append(targets, p.Targets...)
	}

	endTime := time.Now().Unix() * int64(1000)
	startTime := start.Unix() * int64(1000)

	request := GrafanaDataQueryRequest{
		Queries: targets,
		From:    fmt.Sprint(startTime),
		To:      fmt.Sprint(endTime),
	}

	b, err := json.Marshal(&request)
	if err != nil {
		return result, fmt.Errorf("failed to build request object %v", err)
	}

	query := fmt.Sprintf("%v://%v/api/ds/query", c.baseURL.Scheme, c.baseURL.Host)
	req, err := c.NewRequest(http.MethodPost, query, bytes.NewBuffer(b))
	if err != nil {
		return result, fmt.Errorf("failed to build request %v", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return result, err
	}

	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf(string(b))
	}

	err = json.Unmarshal(b, &result)

	return result, err
}

// retrieves a dashboard object from the URL
func (c *grafanaClient) GetDashboardFromURL(urlstr string) (GrafanaDashboardResponse, error) {
	var grafanaDashboardResponse GrafanaDashboardResponse
	uid, _ := ExtractArgs(urlstr)

	if uid == "" {
		return grafanaDashboardResponse, fmt.Errorf("failed to extract uid from url")
	}

	grafanaDashboardResponse, err := c.getDashboard(uid)
	return grafanaDashboardResponse, err
}

// retrieves a dashboard object from a uid
func (c *grafanaClient) GetDashboardWithUID(uid string) (GrafanaDashboardResponse, error) {
	return c.getDashboard(uid)
}

// retrieves the panel data from a URL
func (c *grafanaClient) GetPanelDataFromURL(urlstr string, start time.Time) (Results, error) {
	uid, id := ExtractArgs(urlstr)

	var result Results
	dashboard, err := c.getDashboard(uid)
	if err != nil {
		return result, err
	}

	result, err = c.getPanelData(id, dashboard, start)
	return result, err
}

// retrieves the panel data from an id
func (c *grafanaClient) GetPanelDataFromID(uid string, panelID int, start time.Time) (Results, error) {
	var result Results

	dashboard, err := c.getDashboard(uid)
	if err != nil {
		return result, err
	}

	result, err = c.getPanelData(panelID, dashboard, start)
	return result, err
}

// retrieves the panel data from title
func (c *grafanaClient) GetPanelDataFromTitle(uid string, title string, start time.Time) (Results, error) {
	var result Results

	dashboard, err := c.getDashboard(uid)
	if err != nil {
		return result, err
	}

	for i := range dashboard.Dashboard.Panels {
		p := dashboard.Dashboard.Panels[i]
		if p.Title != title {
			continue
		}
		result, err = c.getPanelData(p.ID, dashboard, start)
		return result, err
	}
	return result, fmt.Errorf("failed to find panel %v", title)
}

// extracts the uid and panel id from a url
func ExtractArgs(urlStr string) (string, int) {
	parsedUrl, err := url.Parse(urlStr)
	if err != nil {
		return "", 0
	}

	segs := strings.Split(parsedUrl.Path, "/")
	var uid string
	if len(segs) >= 3 {
		uid = segs[2]
	} else {
		return "", 0
	}

	viewPanel := parsedUrl.Query().Get("viewPanel")
	if viewPanel == "" {
		return "", 0
	}

	id, err := strconv.ParseInt(viewPanel, 0, 0)
	if err != nil {
		return "", 0
	}

	return uid, int(id)
}
