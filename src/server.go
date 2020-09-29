package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello, world!\n")
}

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func main() {
	// Hello world, the web server

	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/headers", headers)


	fmt.Println("Server is running...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
