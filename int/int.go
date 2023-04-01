package main

import (
	"fmt"
	"math"
)

func main() {
	maxV := max(math.MaxInt, math.MinInt8)
	fmt.Println(maxV)
}

func max(a, b int) int {
	if a < b {
		return b
	}

	return a
}
