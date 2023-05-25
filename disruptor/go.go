package main

import (
	"fmt"
	"time"
)

func main() {
	var arr [][8]int64
	arr = make([][8]int64, 1024*1024)
	for i := 0; i < 1024*1024; i++ {
		for j := 0; j < 8; j++ {
			arr[i][j] = 1
		}
	}

	marked := time.Now().UnixMicro()
	for i := 0; i < 1024*1024; i++ {
		for j := 0; j < 8; j++ {
			arr[i][j] = 2
		}
	}
	fmt.Println("Loop times:", time.Now().UnixMicro()-marked, "ms")

	marked = time.Now().UnixMicro()
	for i := 0; i < 8; i++ {
		for j := 0; j < 1024*1024; j++ {
			arr[j][i] = 2
		}
	}
	fmt.Println("Loop times:", time.Now().UnixMicro()-marked, "ms")
}
