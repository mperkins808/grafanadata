# grafanadata

Query the Grafana API based off dashboards and panels and return the data in either prometheus or grafana format

## Install

```bash
go get github.com/mperkins808/grafanadata/go/pkg/grafanadata
```

## Example usage

```go
package main

import (
	"log"
	"time"

	"github.com/mperkins808/grafanadata/go/pkg/grafanadata"
)

func main() {

	u := "http://localhost:3000/d/bebca380-068d-463d-9c9c-1bb19cb8d2b3/new-dashboard?orgId=1&viewPanel=2"
	t := "glsa_5N21WQvXza0oWkbqQvjOhII8yJYxGS0G_fbb82943"
	client, err := grafanadata.NewGrafanaClient(u, t)

	if err != nil {
		log.Fatal(err)
	}

	start := time.Now().Add(time.Hour * 24 * -2)
	data, err := client.GetPanelDataFromURL(u, start)

	if err != nil {
		log.Fatal(err)
	}

	log.Default().Println(data)
}

```

### Convert response to a Prometheus Format

```go
func main() {

	...

	uid := "bebca380-068d-463d-9c9c-1bb19cb8d2b3"
	panelID := 7

	// Get last 24 hours
	data, err := client.GetPanelDataFromID(uid, panelID, time.Now().Add(time.Hour*24*-1))
	if err != nil {
		return err
	}
	prometheusData := grafanadata.ConvertResultToPrometheusFormat(data)
	log.Println(prometheusData)

}
```

### Iterate through an entire dashboard

```go
func main() {
	...

	uid := "bebca380-068d-463d-9c9c-1bb19cb8d2b3"

	resp, err := client.GetDashboardWithUID(uid)
	if err != nil {
		return
	}

	start := time.Now().Add(time.Hour * 24 * -2)
	for _, panel := range resp.Dashboard.Panels {
		data, _ := client.GetPanelDataFromID(uid, panel.ID, start)
		log.Println(data)
	}
}
```
