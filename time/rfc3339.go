package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now().UTC().Truncate(time.Hour)
	rfc3339 := now.Format(time.RFC3339)
	fmt.Println(rfc3339)

	parse, _ := time.Parse(time.RFC3339, rfc3339)
	parse = parse.Local()
	fmt.Println(parse)
	unix := time.Unix(1678418640, 0)
	fmt.Println(unix.UTC().Format(time.RFC3339))

	fmt.Println(unix.Local())
	fmt.Printf("now:%s\n", time.Now().UTC().Format(time.RFC3339))
	start, end := getCollectorRange(2)
	fmt.Printf("start:%s end:%s", start, end)
}

// http://localhost:9090/api/v1/query_range?query=count(node_cpu_seconds_total{mode=%22system%22})by%20(instance)&start=2023-03-10T06:00:00Z&end=2023-03-10T07:00:00Z&step=300s
// http://localhost:9090/api/v1/query_range?query=count(node_cpu_seconds_total{mode=%22system%22})by%20(instance)&start=2023-03-10T10:00:00Z&end=2023-03-11T11:00:00Z&step=5s

// 根据 scope 获取采集时间范围, scope 单位为 H , 最大为24。
// 案例：nowTime:2023-03-10T03:24:00Z  scope =2,将会返回 2023-03-19T03:00:00Z,2023-03-10904:55:00Z
func getCollectorRange(scope int) (string, string) {
	if scope <= 0 {
		scope = 1
	}

	if scope > 24 {
		scope = 24
	}

	collectorStartTime := time.Now().AddDate(0, 0, -1).Truncate(time.Hour)
	collectorEndTime := collectorStartTime.Add(time.Duration(scope) * time.Hour).Add(-time.Duration(5) * time.Minute)
	return collectorStartTime.UTC().Format(time.RFC3339), collectorEndTime.UTC().Format(time.RFC3339)
}
