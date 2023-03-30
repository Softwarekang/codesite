package main

import (
	"fmt"
	"time"
)

func main() {
	yesterday := time.Now().AddDate(0, 0, -1) // 获取昨天的日期
	yesterdayStart := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, yesterday.Location())
	for i := 1; i < 25; i++ {
		rounded := yesterdayStart.Add(time.Hour * time.Duration(i))
		fmt.Println(rounded.Hour() + i)
		timestamp := rounded.Unix()
		fmt.Println(timestamp)
	}
}

// 1680058800