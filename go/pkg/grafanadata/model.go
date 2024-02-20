package grafanadata

//////////////////////////////////////////////////
// The important parts of a Grafana dashboard json
//////////////////////////////////////////////////

type GrafanaDashboardResponse struct {
	Dashboard Dashboard `json:"dashboard"`
}

type Dashboard struct {
	ID     int     `json:"id"`
	Panels []Panel `json:"panels"`
}

type Panel struct {
	ID         int           `json:"id"`
	Datasource Datasource    `json:"datasource"`
	Targets    []interface{} `json:"targets"`
	Title      string        `json:"title"`
}

type Datasource struct {
	Type string `json:"type"`
	UID  string `json:"uid"`
}

type Target struct {
	AliasBy        string                 `json:"aliasBy"`
	Datasource     Datasource             `json:"datasource"`
	QueryType      string                 `json:"queryType"`
	RefID          string                 `json:"refId"`
	TimeSeriesList map[string]interface{} `json:"timeSeriesList"`
}

type GrafanaDataQueryRequest struct {
	Queries interface{} `json:"queries"`
	From    string      `json:"from"`
	To      string      `json:"to"`
}

type Query struct {
	AliasBy        string                 `json:"aliasBy"`
	Datasource     Datasource             `json:"datasource"`
	QueryType      string                 `json:"queryType"`
	RefID          string                 `json:"refId"`
	TimeSeriesList map[string]interface{} `json:"timeSeriesList"`
}

////////////////////////////////////////////////////////
// Building the request that needs to be sent to Grafana
////////////////////////////////////////////////////////

type Results struct {
	Results map[string]Result `json:"results"`
}

type Result struct {
	Status int     `json:"status"`
	Frames []Frame `json:"frames"`
}

type Frame struct {
	Schema Schema `json:"schema"`
	Data   Data   `json:"data"`
}

type Schema struct {
	RefId  string                 `json:"refId"`
	Meta   map[string]interface{} `json:"meta"`
	Fields []Field                `json:"fields"`
}

type Field struct {
	Name     string                 `json:"name"`
	Type     string                 `json:"type"`
	TypeInfo map[string]interface{} `json:"typeInfo"`
	Config   map[string]interface{} `json:"config"`
	Labels   map[string]string      `json:"labels,omitempty"`
}

type Data struct {
	Values [][]float64 `json:"values"`
}

/////////////////////////////////////////////////
// To Convert Grafana data into prometheus format
/////////////////////////////////////////////////

type PrometheusMetricResponse struct {
	Status string               `json:"status"`
	Data   PrometheusMetricData `json:"data"`
}

type PrometheusMetricData struct {
	ResultType string                       `json:"resultType"`
	Result     []PrometheusMetricDataResult `json:"result"`
}

type PrometheusMetricDataResult struct {
	Metric map[string]string `json:"metric"`
	Values [][]interface{}   `json:"values"`
}

type PrometheusValues struct {
	Timestamp float64 `json:"ts"`
	Value     float64 `json:"value"`
	Index     int     `json:"index"`
}

type PrometheusMetricsList struct {
	Status string   `json:"status"`
	Data   []string `json:"data"`
}

type PrometheusMetricLabels struct {
	Status string              `json:"status"`
	Data   []map[string]string `json:"data"`
}

/////////////////////////////////////////////////
// To Query Dashboards and panels
/////////////////////////////////////////////////

type DashboardSearch struct {
	UID   string `json:"uid"`
	Title string `json:"title"`
}

type PanelSearch struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}
