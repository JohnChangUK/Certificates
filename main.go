package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

import . "./models"

var certificates []Certificate

func getCertificates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(certificates); err != nil {
		log.Fatal("Error encoding to JSON: ", err)
	}
}

func getCertificate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, certificate := range certificates {
		if certificate.Id == params["id"] {
			if err := json.NewEncoder(w).Encode(certificate); err != nil {
				log.Fatal("Error encoding to JSON: ", err)
			}
			return
		}
	}
	if err := json.NewEncoder(w).Encode(&Certificate{}); err != nil {
		log.Fatal("Error encoding to empty Certificate to JSON: ", err)
	}
}

func getUserCertificate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, certificate := range certificates {
		if certificate.Id == params["userId"] {
			if err := json.NewEncoder(w).Encode(certificate); err != nil {
				log.Fatal("Error encoding to JSON: ", err)
			}
			return
		}
	}
	if err := json.NewEncoder(w).Encode(&Certificate{}); err != nil {
		log.Fatal("Error encoding to empty Certificate to JSON: ", err)
	}
}

func createCertificate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var certificate Certificate
	_ = json.NewDecoder(r.Body).Decode(&certificate)
	certificate.CreatedAt = time.Now()
	certificate.Id = strconv.Itoa(rand.Intn(10000000))
	certificates = append(certificates, certificate)
	if err := json.NewEncoder(w).Encode(certificate); err != nil {
		log.Fatal("Error encoding to JSON: ", err)
	}
}

func updateCertificate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, certificate := range certificates {
		if certificate.Id == params["id"] {
			certificates = append(certificates[:index], certificates[index+1:]...)
			var certificate Certificate
			_ = json.NewDecoder(r.Body).Decode(&certificate)
			certificate.Id = params["id"]
			certificates = append(certificates, certificate)
			if err := json.NewEncoder(w).Encode(certificate); err != nil {
				log.Fatal("Error encoding to JSON: ", err)
			}
		}
	}
	if err := json.NewEncoder(w).Encode(certificates); err != nil {
		log.Fatal("Error encoding to JSON: ", err)
	}
}

func deleteCertificate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, certificate := range certificates {
		if certificate.Id == params["id"] {
			certificates = append(certificates[:index], certificates[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(certificates)
}

func main() {
	user := User{"userId1", "user1@gmail.com", "User1"}
	user2 := User{"userId2", "user2@gmail.com", "User2"}
	certificates = append(certificates, Certificate{"1", "First Certificate", time.Time{},
		"John", 2019, "Art note",
		&Transfer{user.Email, "Processing Transfer"}},
		Certificate{"2", "Second Certificate", time.Now(),
			"Jim", 2010, "Painting note",
			&Transfer{user2.Email, "Complete"}})

	router := mux.NewRouter()
	router.HandleFunc("/certificates", getCertificates).Methods("GET")
	router.HandleFunc("/certificates/{id}", getCertificate).Methods("GET")
	router.HandleFunc("/users/{userId}/certificates", getUserCertificate).Methods("GET")
	router.HandleFunc("/certificates", createCertificate).Methods("POST")
	router.HandleFunc("/certificates/{id}", updateCertificate).Methods("PUT")
	router.HandleFunc("/certificates/{id}", deleteCertificate).Methods("DELETE")
	router.HandleFunc("/certificates/{id}/transfers", deleteCertificate).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8000", router))
}
