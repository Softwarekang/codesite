package main

import (
	"fmt"
	"net"
	"os"
	"syscall"
	"time"
)

func main1() {
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

func main() {
	go func() {
		server()
	}()

	go func() {
		client8080()
	}()

	for {
		time.Sleep(2 * time.Second)
	}
}
func server() {
	// 创建一个 TCP 地址
	addr, err := net.ResolveTCPAddr("tcp", ":8080")
	if err != nil {
		fmt.Println("ResolveTCPAddr failed:", err)
		os.Exit(1)
	}

	// 创建一个 TCP 监听器
	l, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		fmt.Println("socket failed:", err)
		os.Exit(1)
	}

	// 绑定地址到监听器
	err = syscall.Bind(l, &syscall.SockaddrInet4{
		Port: addr.Port,
	})
	if err != nil {
		fmt.Println("bind failed:", err)
		os.Exit(1)
	}

	// 开始监听
	err = syscall.Listen(l, syscall.SOMAXCONN)
	if err != nil {
		fmt.Println("listen failed:", err)
		os.Exit(1)
	}

	fmt.Println("listening on", addr.String())

	// 接受连接
	for {
		fd, sa, err := syscall.Accept(l)
		if err != nil {
			fmt.Println("accept failed:", err)
			continue
		}

		fmt.Println("accepted connection from", sa)

		// 处理连接
		go handleConnection(fd)
	}
}

func handleConnection(fd int) {
	defer syscall.Close(fd)

	// 从连接中读取数据并打印
	buf := make([]byte, 1024)
	n, err := syscall.Read(fd, buf)
	if err != nil {
		fmt.Println("read failed:", err)
		return
	}

	fmt.Println("received data:", string(buf[:n]))
}

func client8080() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		return
	}

	defer conn.Close()

	message := "Hello, server. This is the client."

	conn.Write([]byte(message))
}
