# grafanadata

Query the Grafana API based off dashboards and panels and return the data in either prometheus or grafana format

## Example usage

```go

package main

import (
"log"
"time"
)

func  main() {
    u  :=  "http://localhost:3000/d/bebca380-068d-463d-9c9c-1bb19cb8d2b3/new-dashboard?orgId=1&viewPanel=2"
    t  :=  "glsa_5N21WQvXza0oWkbqQvjOhII8yJYxGS0G_fbb82943"
    client, err  :=  NewGrafanaClient(u, t)

    if err !=  nil {
    log.Fatal(err)
    }

    start  := time.Now().Add(time.Hour *  24  *  -2)
    data, err  := client.GetPanelDataFromURL(u, start)

    if err !=  nil {
    log.Fatal(err)
    }

    logger  := log.Default().Println(data)
}
```
