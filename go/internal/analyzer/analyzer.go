package analyzer

import (
	"bufio"
	"encoding/json"
	"os"
	"regexp"
	"sort"
	"strconv"

	"loganalyzer/internal/models"
)

type Analyzer struct {
	stats models.LogStats
	paths map[string]*models.EndpointStat
}

func NewAnalyzer() *Analyzer {
	return &Analyzer{
		stats: models.LogStats{
			StatusCodes:  make(map[string]int),
			TopEndpoints: make([]models.EndpointStat, 0),
		},
		paths: make(map[string]*models.EndpointStat),
	}
}

func (a *Analyzer) AnalyzeFile(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	totalTime := 0.0

	// 示例日志格式: [2024-03-22 15:04:05] "GET /api/v1/users" 200 0.123
	logPattern := regexp.MustCompile(`\[(.*?)\] "(.*?) (.*?)" (\d{3}) (\d+\.?\d*)`)

	for scanner.Scan() {
		line := scanner.Text()
		matches := logPattern.FindStringSubmatch(line)
		if len(matches) < 6 {
			continue
		}

		// 解析日志组件
		timestamp := matches[1]
		method := matches[2]
		path := matches[3]
		statusCode := matches[4]
		responseTime, _ := strconv.ParseFloat(matches[5], 64)

		// 更新统计信息
		a.stats.TotalRequests++
		a.stats.StatusCodes[statusCode]++
		totalTime += responseTime

		// 更新端点统计
		key := method + " " + path
		if _, exists := a.paths[key]; !exists {
			a.paths[key] = &models.EndpointStat{
				Path:    path,
				Method:  method,
				Count:   0,
				AvgTime: 0,
			}
		}
		stat := a.paths[key]
		stat.Count++
		stat.AvgTime = (stat.AvgTime*float64(stat.Count-1) + responseTime) / float64(stat.Count)

		// 更新时间范围
		if a.stats.StartTime == "" || timestamp < a.stats.StartTime {
			a.stats.StartTime = timestamp
		}
		if timestamp > a.stats.EndTime {
			a.stats.EndTime = timestamp
		}
	}

	// 计算平均响应时间
	if a.stats.TotalRequests > 0 {
		a.stats.AverageTime = totalTime / float64(a.stats.TotalRequests)
	}

	// 转换并排序端点统计
	for _, stat := range a.paths {
		a.stats.TopEndpoints = append(a.stats.TopEndpoints, *stat)
	}
	sort.Slice(a.stats.TopEndpoints, func(i, j int) bool {
		return a.stats.TopEndpoints[i].Count > a.stats.TopEndpoints[j].Count
	})

	// 只保留前 10 个最常访问的端点
	if len(a.stats.TopEndpoints) > 10 {
		a.stats.TopEndpoints = a.stats.TopEndpoints[:10]
	}

	return scanner.Err()
}

func (a *Analyzer) SaveStats(outputPath string) error {
	data, err := json.MarshalIndent(a.stats, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(outputPath, data, 0644)
}

func (a *Analyzer) GetStats() models.LogStats {
	return a.stats
}
