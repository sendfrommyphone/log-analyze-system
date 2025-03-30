package models

type LogStats struct {
	TotalRequests int            `json:"total_requests"`
	StatusCodes   map[string]int `json:"status_codes"` // 修改为 string 类型的 key
	AverageTime   float64        `json:"average_time"`
	TopEndpoints  []EndpointStat `json:"top_endpoints"`
	StartTime     string         `json:"start_time"` // 添加时间范围
	EndTime       string         `json:"end_time"`
}

type EndpointStat struct {
	Path    string  `json:"path"`
	Count   int     `json:"count"`
	AvgTime float64 `json:"avg_time"`
	Method  string  `json:"method"` // 添加 HTTP 方法
}
