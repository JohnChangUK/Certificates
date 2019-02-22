package main

import (
	"encoding/json"
	. "github.com/JohnChangUK/verisart/models"
	. "github.com/JohnChangUK/verisart/utils"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var certificates []Certificate
var users map[string]User

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
			//certificate.Id = params["id"]
			certificates = append(certificates, certificate)
			if err := json.NewEncoder(w).Encode(certificate); err != nil {
				log.Fatal("Error encoding to JSON: ", err)
			}
			return
		}
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

func createTransfer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, certificate := range certificates {
		if certificate.Id == params["id"] {
			certificates = append(certificates[:index], certificates[index+1:]...)
			//var certificate Certificate
			var user User
			// Create a new Certificate with new User Details
			// Only changed when other person ACCEPTS the transfer
			_ = json.NewDecoder(r.Body).Decode(&user)
			certificate.Transfer.To = user.Email
			certificate.Transfer.Status = "TRANSFER_IN_PROGRESS"
			certificates = append(certificates, certificate)
			if err := json.NewEncoder(w).Encode(certificate); err != nil {
				log.Fatal("Error encoding to JSON: ", err)
			}
			return
		}
	}
}

func updateTransfer(w http.ResponseWriter, r *http.Request) {

}

func main() {
	certificates = MockCertificates(certificates)

	router := mux.NewRouter()
	router.HandleFunc("/certificates", getCertificates).Methods("GET")
	router.HandleFunc("/certificates/{id}", getCertificate).Methods("GET")
	router.HandleFunc("/users/{userId}/certificates", getUserCertificate).Methods("GET")
	router.HandleFunc("/certificates", createCertificate).Methods("POST")
	router.HandleFunc("/certificates/{id}", updateCertificate).Methods("PUT")
	router.HandleFunc("/certificates/{id}", deleteCertificate).Methods("DELETE")
	router.HandleFunc("/certificates/{id}/transfers", createTransfer).Methods("POST")
	router.HandleFunc("/certificates/{id}/transfers", updateTransfer).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8000", router))
}
