package analyzer

import (
	"bufio"
	"bytes"
	"strings"
)

type AnalysisResult struct {
	TotalLines   int                    `json:"total_lines"`
	FileSize     int64                  `json:"file_size"`
	LinePatterns map[string]int         `json:"line_patterns"`
	ErrorCount   int                    `json:"error_count"`
	WarningCount int                    `json:"warning_count"`
	Statistics   map[string]interface{} `json:"statistics"`
}

func Analyze(data []byte) AnalysisResult {
	result := AnalysisResult{
		LinePatterns: make(map[string]int),
		Statistics:   make(map[string]interface{}),
	}

	// 计算文件大小
	result.FileSize = int64(len(data))

	// 使用 bufio.Scanner 逐行读取
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		result.TotalLines++

		// 计算行模式统计
		pattern := getLinePattern(line)
		result.LinePatterns[pattern]++

		// 统计错误和警告
		lower := strings.ToLower(line)
		if strings.Contains(lower, "error") {
			result.ErrorCount++
		}
		if strings.Contains(lower, "warning") {
			result.WarningCount++
		}
	}

	// 添加基本统计信息
	result.Statistics["avg_line_length"] = calculateAvgLineLength(data)

	return result
}

func getLinePattern(line string) string {
	// 简单的行模式识别逻辑
	switch {
	case strings.HasPrefix(line, "[ERROR]"):
		return "error_log"
	case strings.HasPrefix(line, "[WARN]"):
		return "warning_log"
	case strings.HasPrefix(line, "[INFO]"):
		return "info_log"
	case strings.HasPrefix(line, "[DEBUG]"):
		return "debug_log"
	default:
		return "other"
	}
}

func calculateAvgLineLength(data []byte) float64 {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	var totalLength int
	var lineCount int

	for scanner.Scan() {
		totalLength += len(scanner.Text())
		lineCount++
	}

	if lineCount == 0 {
		return 0
	}
	return float64(totalLength) / float64(lineCount)
}
