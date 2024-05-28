package utils

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/teris-io/shortid"
)

type service struct {
	db *sql.DB
}

func idgen() (string, error) {
	sid, err := shortid.New(1, shortid.DefaultABC, 2342)
	if err != nil {
		return "", err
	}

	id, err := sid.Generate()
	if err != nil {
		return "", err
	}

	return id, err
}

func dbinsert(originalurl, generatedid string) error {
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

func NewService(db *sql.DB) *service {
	return &service{db: db}
}

func (svc *service) ShortenURLHandler(w http.ResponseWriter, r *http.Request) {
	queryURL := r.URL.Query().Get("url")
	if queryURL == "" {
		http.Error(w, "No URL provided.", http.StatusBadRequest)
		return
	}

	sUrl, err := idgen()
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, "Error occurred while generating shortid", http.StatusInternalServerError)
		return
	}

	if err := dbinsert(queryURL, sUrl); err != nil {
		log.Println("Error:", err)
		http.Error(w, "Error occurred while inserting into DB", http.StatusInternalServerError)
		return
	}

	// defer db.close()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{
		"SURL": sUrl,
	}); err != nil {
		log.Println("Error encoding JSON:", err)
		http.Error(w, "Error occurred while encoding JSON", http.StatusInternalServerError)
		return
	}
}
