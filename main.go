package main

import (
	"encoding/json"
	. "github.com/JohnChangUK/verisart/model"
	"github.com/JohnChangUK/verisart/utils"
	"github.com/gorilla/mux"
	"github.com/rs/xid"
	"log"
	"net/http"
	"time"
)

var certificates []Certificate
var users map[string]User

func getCertificates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(certificates); err != nil {
		log.Println("Error encoding to JSON: ", err)
	}
}

func getCertificate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, certificate := range certificates {
		if certificate.Id == params["id"] {
			if err := json.NewEncoder(w).Encode(certificate); err != nil {
				log.Println("Error encoding to JSON: ", err)
			}
			return
		}
	}
	if err := json.NewEncoder(w).Encode(&Certificate{}); err != nil {
		log.Println("Error encoding to empty Certificate to JSON: ", err)
	}
}

func getUserCertificate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, certificate := range certificates {
		if certificate.Id == params["userId"] {
			if err := json.NewEncoder(w).Encode(certificate); err != nil {
				log.Println("Error encoding to JSON: ", err)
			}
			return
		}
	}
	if err := json.NewEncoder(w).Encode(&Certificate{}); err != nil {
		log.Println("Error encoding to empty Certificate to JSON: ", err)
	}
}

func createCertificate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userId := r.Header.Get("Authorization")
	var certificate Certificate
	_ = json.NewDecoder(r.Body).Decode(&certificate)
	certificate.CreatedAt = time.Now()
	certificate.OwnerId = userId
	certificate.Id = xid.New().String()
	certificate.Transfer = &Transfer{}
	certificates = append(certificates, certificate)
	if err := json.NewEncoder(w).Encode(certificate); err != nil {
		log.Println("Error encoding to JSON: ", err)
	}
}

func updateCertificate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userId := r.Header.Get("Authorization")
	params := mux.Vars(r)
	for index, certificate := range certificates {
		if certificate.Id == params["id"] {
			dateCreated := certificate.CreatedAt
			certificates = append(certificates[:index], certificates[index+1:]...)
			var certificate Certificate
			_ = json.NewDecoder(r.Body).Decode(&certificate)
			certificate.OwnerId = userId
			certificate.Id = params["id"]
			certificate.CreatedAt = dateCreated
			certificates = append(certificates, certificate)
			if err := json.NewEncoder(w).Encode(certificate); err != nil {
				log.Println("Error encoding to JSON: ", err)
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
				log.Println("Error encoding to JSON: ", err)
			}
			return
		}
	}
}

func updateTransfer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Certificate.OwnerId must be changed the new User.Id
	userId := r.Header.Get("Authorization")
	params := mux.Vars(r)
	var oldCertificate Certificate
	for index, certificate := range certificates {
		if certificate.Id == params["id"] {
			oldCertificate = certificate
			dateCreated := certificate.CreatedAt
			certificates = append(certificates[:index], certificates[index+1:]...)
			var transfer Transfer
			_ = json.NewDecoder(r.Body).Decode(&transfer)
			oldCertificate.OwnerId = userId
			oldCertificate.CreatedAt = dateCreated
		}
	}
	if err := json.NewEncoder(w).Encode(oldCertificate); err != nil {
		log.Println("Error encoding to JSON: ", err)
	}
	return
}

func main() {
	certificates = utils.MockCertificates(certificates)
	startHttpServer()
}

func startHttpServer() {
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
