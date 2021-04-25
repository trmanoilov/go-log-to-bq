package structs

type AccessLogEntry struct {
	Timestamp string `bigquery:"timestamp" json:"timestamp"`
	App       string `bigquery:"app" json:"app"`
	Request   string `bigquery:"request" json:"request"`
	Code      int    `bigquery:"code" json:"code"`
}

type UsageLogEntry struct {
	Timestamp string  `bigquery:"timestamp" json:"timestamp"`
	App       string  `bigquery:"app" json:"app"`
	CPU       float64 `bigquery:"cpu" json:"cpu"`
	Memory    float64 `bigquery:"memory" json:"memory"`
}

type SlackRequestBody struct {
	Text string `json:"text"`
}
