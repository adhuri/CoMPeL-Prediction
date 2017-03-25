package fetcher

type DataPoint struct {
	AgentIp     string
	ContainerId string
	Timestamp   int64
	Value       float32
	MetricType  string
}
