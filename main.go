package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", ShortenURLHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
