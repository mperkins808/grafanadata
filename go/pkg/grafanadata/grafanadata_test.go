package grafanadata

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"
)

type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func createMockClient(t *testing.T, file string, expectedCode int) *MockHTTPClient {
	path := fmt.Sprintf("./test/%v", file)
	return &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			file, err := os.Open(path)
			if err != nil {
				t.Fatal(err)
			}
			return &http.Response{
				StatusCode: expectedCode,
				Body:       io.NopCloser(file),
			}, nil
		},
	}
}

func createMockGrafanaClient(t *testing.T, mockClient *MockHTTPClient) *grafanaClient {
	return &grafanaClient{
		baseURL: &url.URL{
			Scheme: "http",
			Host:   "example.com",
		},
		token:  "test_token",
		client: mockClient,
	}

}

func TestCreateGrafanaClient(t *testing.T) {

	c, err := NewGrafanaClient("foo", "bar")
	if err != nil {
		t.Fatalf("creating new grafana client error %v", err)
	}

	_, err = c.GetDashboardFromURL("foo")
	if err == nil {
		t.Fatalf("should have failed with bad url")
	}
}

func TestGetDashboard(t *testing.T) {

	client := createMockClient(t, "dashboard.json", http.StatusOK)

	g := createMockGrafanaClient(t, client)

	dashboard, err := g.getDashboard("foo")
	if err != nil {
		t.Fatal(err)
	}

	panels := len(dashboard.Dashboard.Panels)
	if panels != 2 {
		t.Fatalf("wanted 2 panels. got %v", panels)
	}

	client = createMockClient(t, "dashboard.json", http.StatusNotFound)

	g = createMockGrafanaClient(t, client)
	_, err = g.getDashboard("foo")
	if err == nil {
		t.Fatal("wanted error but was nil")
	}
}

func TestGetPanelData(t *testing.T) {
	// loading in the dashboard
	client := createMockClient(t, "dashboard.json", http.StatusOK)

	g := createMockGrafanaClient(t, client)

	dashboard, err := g.getDashboard("foo")
	if err != nil {
		t.Fatal(err)
	}

	// loading in the panel
	client = createMockClient(t, "data.json", http.StatusOK)

	g = createMockGrafanaClient(t, client)

	data, err := g.getPanelData(0, dashboard, time.Now())
	if err != nil {
		t.Fatal(err)
	}

	lenRes := len(data.Results)
	if lenRes != 2 {
		t.Fatalf("wanted len 2 but was %v", lenRes)
	}

}

func TestExtractArgs(t *testing.T) {
	u := "https://grafana.com/d/foobar/fizz?viewPanel=4"
	uid, id := ExtractArgs(u)
	if uid != "foobar" {
		t.Fatalf("wanted foobar. was %v", uid)
	}

	if id != 4 {
		t.Fatalf("wanted 4. was %v", id)
	}

}
