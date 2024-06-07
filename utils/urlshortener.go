package utils

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

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

// func dbinsert(originalurl, generatedid string) error {
// 	file, err := os.OpenFile("dummydb.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 	if err != nil {
// 		return err
// 	}

// 	w := csv.NewWriter(file)
// 	if err := w.Write([]string{generatedid, originalurl}); err != nil {
// 		return err
// 	}

// 	w.Flush()

// 	return w.Error()
// }



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

	if _, err := svc.db.Exec("INSERT INTO urls (original_link, short_id) VALUES (?, ?)", queryURL, sUrl); err != nil {
		log.Println("Error:", err)
		http.Error(w, "Error occurred while inserting into DB", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{
		"SURL": sUrl,
	}); err != nil {
		log.Println("Error encoding JSON:", err)
		http.Error(w, "Error occurred while encoding JSON", http.StatusInternalServerError)
		return
	}
}

func (svc *service) RedirectHandler(w http.ResponseWriter, r *http.Request) {
    shortID := r.URL.Query().Get("short_id") //gets short_id param from retrieve?short_id
    if shortID == "" {
        http.Error(w, "No short ID provided.", http.StatusBadRequest) //when no short_id is provided
        return
    }

    var originalURL string // variable to store original URl
    // line to search the row and get orignal_link from the "urls" table by passing short_id 
	// as a query and then scan the original_url and assings it to originalURL 
	err := svc.db.QueryRow("SELECT original_link FROM urls WHERE short_id = ?", shortID).Scan(&originalURL) 
    if err != nil {
        log.Println("Error:", err)
        http.Error(w, "Error occurred while retrieving original URL", http.StatusInternalServerError) //send error if nothing is found
        return
    }

	//this line redirects to the originalURL
    http.Redirect(w, r, originalURL, http.StatusFound)
}