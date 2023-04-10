package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "10.122.2.71:8000")
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		return
	}

	defer conn.Close()

	message := "Hello, server. This is the client."

	go func() {
		for {
			n, err := conn.Write([]byte(message))
			if err != nil {
				fmt.Println(err)
			}

			fmt.Printf("client write len:%v\n", n)
		}

	}()

	go func() {
		for {
			buf := make([]byte, 1000)
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Printf("client read len:%v\n", n)
		}

	}()

	for {
		time.Sleep(2 * time.Second)
	}
}
