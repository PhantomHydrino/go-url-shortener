package main

import (
	"log"
	"net/http"

	"github.com/PhantomHydrino/go-url-shortener/utils"
)

func main() {
	http.HandleFunc("/", utils.ShortenURLHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
