package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	inputPath             = "tmp_restore/exported.csv"
	outputPath            = "tmp_restore/out_shifted.lp"
	originalEndTimeString = "2025-05-30T00:00:00Z"
)

func escape(value string) string {
	return strings.NewReplacer(" ", "\\ ", ",", "\\,", "=", "\\=").Replace(value)
}

func isNumeric(s string) bool {
	_, err := fmt.Sscanf(s, "%f", new(float64))
	return err == nil
}

func isHeaderRow(row []string) bool {
	required := map[string]bool{
		"_time":  false,
		"_value": false,
		"_field": false,
	}
	for _, col := range row {
		if _, ok := required[col]; ok {
			required[col] = true
		}
	}
	for _, found := range required {
		if !found {
			return false
		}
	}
	return true
}

func main() {
	// 시간 기준 계산
	originalEndTime, err := time.Parse(time.RFC3339, originalEndTimeString)
	if err != nil {
		panic(err)
	}
	timeDiff := time.Since(originalEndTime)

	// 파일 열기
	inputFile, err := os.Open(inputPath)
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	var headers []string
	recordCount := 0
	skippedCount := 0

	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		reader := csv.NewReader(strings.NewReader(line))
		reader.FieldsPerRecord = -1
		row, err := reader.Read()
		if err != nil {
			skippedCount++
			continue
		}

		if isHeaderRow(row) {
			headers = row
			continue
		}

		if len(headers) == 0 || len(row) != len(headers) {
			skippedCount++
			continue
		}

		// Map 형태로 변환
		record := make(map[string]string)
		for i, key := range headers {
			record[key] = row[i]
		}

		measurement := record["_measurement"]
		field := record["_field"]
		value := record["_value"]
		timeStr := record["_time"]

		if measurement == "" || field == "" || timeStr == "" {
			skippedCount++
			continue
		}

		timestamp, err := time.Parse(time.RFC3339, timeStr)
		if err != nil {
			skippedCount++
			continue
		}
		shifted := timestamp.Add(timeDiff)
		timestampNs := shifted.UnixNano()

		// 태그 구성
		var tags []string
		for _, key := range headers {
			if strings.HasPrefix(key, "_") || key == "result" || key == "table" {
				continue
			}
			val := record[key]
			if val != "" && val != "None" {
				tags = append(tags, fmt.Sprintf("%s=%s", escape(key), escape(val)))
			}
		}
		tagStr := strings.Join(tags, ",")

		// 필드 값 처리
		if !isNumeric(value) {
			value = fmt.Sprintf(`"%s"`, value)
		}

		var lineOut string
		if tagStr != "" {
			lineOut = fmt.Sprintf("%s,%s %s=%s %d", measurement, tagStr, field, value, timestampNs)
		} else {
			lineOut = fmt.Sprintf("%s %s=%s %d", measurement, field, value, timestampNs)
		}

		_, err = writer.WriteString(lineOut + "\n")
		if err != nil {
			panic(err)
		}

		recordCount++
		if recordCount%100000 == 0 {
			fmt.Printf("✅ %d rows processed...\n", recordCount)
		}
	}

	fmt.Printf("🎯 변환 완료: %d행 처리됨, %d행 건너뜀\n", recordCount, skippedCount)
}
