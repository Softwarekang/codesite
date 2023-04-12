package main

import (
	"fmt"
)

func main() {
	src := make([]byte, 10)
	dst := make([]byte, 1)
	n := copy(dst, src)
	fmt.Println(n)

	dst = make([]byte, 5)
	dst2 := dst[:0]

	dst3 := dst2[:cap(dst2)]
	fmt.Println(len(dst2))
	fmt.Println(len(dst3))
}
