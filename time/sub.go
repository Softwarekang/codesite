package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()
	befor2H := now.AddDate(0, 0, -1).Add(time.Minute)

	fmt.Println(now.Sub(befor2H).Hours())

	now = time.Now()                                                                                                                                    // 获取当前时间
	yesterday := now.AddDate(0, 0, -1)                                                                                                                  // 获取昨天的日期
	yesterdayStart := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, yesterday.Location())                                 // 获取昨天的起始时间点
	yesterdayEnd := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), yesterday.Location()) // 获取昨天的结束时间点

	fmt.Println("昨天的起始时间点：", yesterdayStart)
	fmt.Println("昨天的结束时间点：", yesterdayEnd)

}
