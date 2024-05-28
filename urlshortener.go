package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/teris-io/shortid"
)

func shortenurl(originalurl string) (string, error) {
	// //generator string
	// mapChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// var shortUrl strings.Builder

	// for n > 0 {
	// 	shortUrl.WriteByte(mapChars[n%62])
	// 	n = n / 62
	// }

	// return reverseString(shortUrl.String())

	sid, err := shortid.New(1, shortid.DefaultABC, 2342)

	if err != nil {
		return "", err
	}

	id, err := sid.Generate()
	if err != nil {
		return "", err
	}

	err = dbinsert(originalurl, id)
	if err != nil {
		return "", err
	}
	return id, err
}

func dbinsert(originalurl, generatedid string) error {

	// record := [][]string{
	// 	{"generatedid", "originalurl"},
	// }

	file, err := os.OpenFile("dummydb.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	w := csv.NewWriter(file)

	if err := w.Write([]string{generatedid, originalurl}); err != nil {
		return err
	}

	w.Flush()

	return w.Error()
}

// func reverseString(s string) string {
// 	runes := []r"une(s)
// 	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
// 		runes[i], runes[j] = runes[j], runes[i]
// 	}
// 	return string(runes)
// }

func main() {

	var inputurl string

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintln(w, "Please type a URL:")

		queryURL := r.URL.Query().Get("url")

		if queryURL == "" {
			fmt.Fprintln(w, "No URL provided.")
			return
		}

		sUrl, err := shortenurl(inputurl)

		if err != nil {
			log.Println("Error is:", err)
		}

		err = dbinsert(queryURL, sUrl) // Pass queryURL to dbinsert
		if err != nil {
			log.Println("Error:", err)
			http.Error(w, "Error occurred while inserting into DB", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "\nShortened URL is: %s\n", sUrl)

	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
