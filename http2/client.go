package main

import (
	"crypto/tls"
	"fmt"
	"golang.org/x/net/http2"
	"net/http"
)

func main() {
	// 创建一个HTTP/2客户端传输配置
	config := &http2.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	// 创建一个HTTP/2客户端
	client := &http.Client{
		Transport: config,
	}

	// 发送一个HTTP/2 GET请求
	resp, err := client.Get("https://localhost:8443/")
	if err != nil {
		fmt.Println(err)
	}

	// 读取响应并关闭响应体
	body := make([]byte, resp.ContentLength)
	resp.Body.Read(body)
	resp.Body.Close()

	fmt.Println(string(body))
}
