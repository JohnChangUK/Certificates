package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type Certificate struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	OwnerId   string    `json:"owner_id"`
	Year      int       `json:"year"`
	Note      string    `json:"note"`
	Transfer  *Transfer `json:"transfer"`
}

type Transfer struct {
	To     string `json:"email"`
	Status string `json:"status"`
}

type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

var certificates []Certificate

func getCertificates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(certificates)
}

func getCertificate(w http.ResponseWriter, r *http.Request) {

}

func createCertificate(w http.ResponseWriter, r *http.Request) {

}

func updateCertificate(w http.ResponseWriter, r *http.Request) {

}

func deleteCertificate(w http.ResponseWriter, r *http.Request) {

}

func main() {
	certificates = append(certificates, Certificate{"1", "First Certificate", time.Now(),
		"John", 2019, "Art note",
		&Transfer{"john@gmail.com", "Complete"}})

	router := mux.NewRouter()
	router.HandleFunc("/certificates", getCertificates).Methods("GET")
	router.HandleFunc("/certificates/{id}", getCertificate).Methods("GET")
	router.HandleFunc("/certificates", createCertificate).Methods("POST")
	router.HandleFunc("/certificates/{id}", updateCertificate).Methods("PUT")
	router.HandleFunc("/certificates/{id}", deleteCertificate).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}
