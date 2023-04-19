package main

import (
	"fmt"
	"sync"
)

func main() {
	m := sync.Map{}
	m.Store(1, 2)
	m.Store(1, 3)
	m.Store(12, 2)
	m.Range(func(key, value any) bool {
		fmt.Printf("k:%v, v:%v", key, value)
		return true
	})
}
