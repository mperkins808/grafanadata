package grafanadata

// Convert a Grafana data response into prometheus format
func ConvertResultToPrometheusFormat(results Results) PrometheusMetricResponse {
	promResponse := PrometheusMetricResponse{
		Status: "success",
		Data: PrometheusMetricData{
			ResultType: "matrix",
		},
	}

	for _, result := range results.Results {
		for _, frame := range result.Frames {

			var promResult PrometheusMetricDataResult

			metricLabels := make(map[string]string)
			for _, field := range frame.Schema.Fields {
				for labelKey, labelValue := range field.Labels {
					metricLabels[labelKey] = labelValue
				}
			}

			promResult.Metric = metricLabels
			if len(frame.Data.Values) >= 2 {
				timestamps := frame.Data.Values[0]
				values := frame.Data.Values[1]

				for index, timestamp := range timestamps {
					if index < len(values) {
						value := values[index]
						promResult.Values = append(promResult.Values, []interface{}{timestamp / 1000, value})
					}
				}
			}

			if len(promResult.Values) > 0 {
				promResponse.Data.Result = append(promResponse.Data.Result, promResult)
			}
		}
	}

	return promResponse
}
