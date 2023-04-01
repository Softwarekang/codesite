package main

import "fmt"

func main() {
	src := make([]byte, 10)
	dst := make([]byte, 1)
	n := copy(dst, src)
	fmt.Println(n)
}
