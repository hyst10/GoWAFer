package api_handler

// LogStats 用于表示图表数据的结构
type LogStats struct {
	Time     string `json:"time"`
	Count200 int    `json:"count200"`
	Count403 int    `json:"count403"`
}

type LogStatsResponse struct {
	TimeSegments []string `json:"timeSegments"`
	Count200s    []int    `json:"count200s"`
	Count403s    []int    `json:"count403s"`
}
