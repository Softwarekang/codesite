package main

import (
	"fmt"
	"net/http"
)

func main() {
	server := &http.Server{
		Addr: ":8443",
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, HTTP/2")
	})

	err := server.ListenAndServeTLS("./crt/server.crt", "./crt/server.key")
	if err != nil {
		fmt.Println(err)
	}
}
