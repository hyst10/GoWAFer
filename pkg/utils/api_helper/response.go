package api_helper

// Response 通用响应结构
type Response struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

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
