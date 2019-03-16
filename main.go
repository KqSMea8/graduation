package main

import (
	"net/http"
)

func main() {
	initHttpHandler()
}

func initHttpHandler() {
	http.HandleFunc("/post", post)
	http.HandleFunc("/update", delete_)
}

