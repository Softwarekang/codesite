// server.go

package main

import (
	"fmt"
	"net"
	"os"
	"syscall"
)

const (
	PORT    = "8888"
	BACKLOG = 5
	MAXBUF  = 1024
)

var (
	readyCh = make(chan struct{})
)

func main() {
	go server()
	<-readyCh
	client()
}
func server() {
	// 创建 Socket
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		fmt.Printf("socket creation error: %s\n", err)
		os.Exit(1)
	}
	defer syscall.Close(fd)

	// 绑定 Socket
	laddr, err := net.ResolveTCPAddr("tcp", ":"+PORT)
	if err != nil {
		fmt.Printf("resolve address error: %s\n", err)
		os.Exit(1)
	}
	err = syscall.Bind(fd, &syscall.SockaddrInet4{Port: laddr.Port})
	if err != nil {
		fmt.Printf("bind error: %s\n", err)
		os.Exit(1)
	}

	// 监听 Socket
	err = syscall.Listen(fd, BACKLOG)
	if err != nil {
		fmt.Printf("listen error: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("server listening on", laddr.String())

	close(readyCh)
	// 接受客户端连接
	cfd, _, err := syscall.Accept(fd)
	if err != nil {
		fmt.Printf("accept error: %s\n", err)
		return
	}

	// 接收客户端数据
	buf := make([]byte, MAXBUF)
	n, err := syscall.Read(cfd, buf)
	if err != nil {
		fmt.Printf("read error: %s\n", err)
		syscall.Close(cfd)
		return
	}

	fmt.Printf("server received message from client: %s\n", string(buf[:n]))

	// 发送响应数据
	resp := "Hello, client!"
	_, err = syscall.Write(cfd, []byte(resp))
	if err != nil {
		fmt.Printf("write error: %s\n", err)
	}

	syscall.Close(cfd)
}

func client() {
	// 创建一个IPv4的TCP套接字
	socket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		fmt.Println("Socket creation failed:", err)
		return
	}

	// 设置服务器地址
	serverAddr := &syscall.SockaddrInet4{
		Port: 8888,
		Addr: [4]byte{127, 0, 0, 1},
	}

	// 连接到服务器
	err = syscall.Connect(socket, serverAddr)
	if err != nil {
		fmt.Println("Connection failed:", err)
		return
	}

	// 发送数据到服务器
	data := []byte("Hello, server!")
	_, err = syscall.Write(socket, data)
	if err != nil {
		fmt.Println("Write failed:", err)
		return
	}

	// 接收来自服务器的数据
	buffer := make([]byte, 1024)
	_, err = syscall.Read(socket, buffer)
	if err != nil {
		fmt.Println("Read failed:", err)
		return
	}

	// 输出接收到的数据
	fmt.Println("Received from server:", string(buffer))
}
