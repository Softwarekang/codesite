package main

import (
	"fmt"
	"time"
)

func main() {
	var CST = time.FixedZone("CST", 3600*8)

	date := time.Date(2023, time.February, 40, 0, 0, 0, 0, CST)
	fmt.Printf(date.String())
}

// 1680058800
