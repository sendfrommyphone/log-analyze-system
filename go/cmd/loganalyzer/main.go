package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"loganalyzer/pkg/analyzer" // 修改导入路径
)

func main() {
	inputFile := flag.String("input", "", "Input file path")
	flag.Parse()

	if *inputFile == "" {
		log.Fatal("Input file is required")
	}

	// 读取文件
	data, err := os.ReadFile(*inputFile)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// 分析文件
	result := analyzer.Analyze(data)

	// 转换为JSON输出
	jsonResult, err := json.Marshal(result)
	if err != nil {
		log.Fatalf("Failed to marshal result: %v", err)
	}

	fmt.Println(string(jsonResult))
}
