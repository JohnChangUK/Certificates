package utils

import (
	"encoding/json"
	"fmt"
	. "github.com/JohnChangUK/verisart/model"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func EncodeToJson(w http.ResponseWriter, v ...interface{}) {
	err := json.NewEncoder(w).Encode(v)

	if err != nil {
		log.Fatal("Error encoding to JSON: ", err)
	}
}

func DecodeFromJson(w http.ResponseWriter, req *http.Request, v interface{}) {
	err := json.NewDecoder(req.Body).Decode(&v)
	req.Body.Close()

	if err != nil {
		InvalidBody(w, req, err)
		return
	}
}

func MockCertificates(certificates []Certificate) []Certificate {
	certificates = append(certificates, Certificate{Id: "1", Title: "First Certificate", CreatedAt: time.Now(),
		OwnerId: "John", Year: 2019, Note: "Blockchain",
		Transfer: &Transfer{}},
		Certificate{Id: "2", Title: "Second Certificate", CreatedAt: time.Now(),
			OwnerId: "John", Year: 3000, Note: "Art",
			Transfer: &Transfer{}},
		Certificate{Id: "3", Title: "Third Certificate", CreatedAt: time.Now(),
			OwnerId: "Jim", Year: 2010, Note: "Painting",
			Transfer: &Transfer{}})

	return certificates
}

func InvalidBody(w http.ResponseWriter, req *http.Request, err error) {
	badRequest(w, fmt.Sprintf("Invalid request: %s, error : %s\n", req.URL, err))
}

func badRequest(w http.ResponseWriter, message string) {
	log.Error(message)
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, message)
}
