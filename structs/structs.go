package structs

type LogEntry struct {
	Timestamp string `bigquery:"timestamp" json:"timestamp"`
	App       string `bigquery:"app" json:"app"`
	Request   string `bigquery:"request" json:"request"`
	Code      int    `bigquery:"code" json:"code"`
}
