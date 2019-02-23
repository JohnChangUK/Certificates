package main

import (
	. "github.com/JohnChangUK/verisart/model"
	. "github.com/JohnChangUK/verisart/utils"
	"github.com/gorilla/mux"
	"github.com/rs/xid"
	"log"
	"net/http"
	"time"
)

var certificates []Certificate

func getUserCertificate(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	authorization := req.Header.Get("Authorization")
	var usersCertificates []Certificate

	for _, certificate := range certificates {
		if authorization == params["userId"] && authorization == certificate.OwnerId {
			usersCertificates = append(usersCertificates, certificate)
		}
	}

	EncodeToJson(w, usersCertificates)
	return
}

func createCertificate(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userId := req.Header.Get("Authorization")

	var certificate Certificate
	DecodeFromJson(w, req, &certificate)

	certificate.CreatedAt = time.Now()
	certificate.OwnerId = userId
	certificate.Id = xid.New().String()
	certificate.Transfer = &Transfer{}
	certificates = append(certificates, certificate)

	EncodeToJson(w, certificate)
}

func updateCertificate(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	authorization := req.Header.Get("Authorization")

	for index, certificate := range certificates {
		if certificate.Id == params["id"] && certificate.OwnerId == authorization {
			dateCreated := certificate.CreatedAt
			certificates = append(certificates[:index], certificates[index+1:]...)

			var certificate Certificate
			DecodeFromJson(w, req, &certificate)

			certificate.OwnerId = authorization
			certificate.Id = params["id"]
			certificate.CreatedAt = dateCreated
			certificate.Transfer = &Transfer{}
			certificates = append(certificates, certificate)

			EncodeToJson(w, certificate)
			return
		}
	}
}

func deleteCertificate(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	authorization := req.Header.Get("Authorization")

	for index, certificate := range certificates {
		if certificate.Id == params["id"] && certificate.OwnerId == authorization {
			certificates = append(certificates[:index], certificates[index+1:]...)
			break
		}
	}

	EncodeToJson(w, certificates)
}

func createTransfer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	authorization := req.Header.Get("Authorization")

	for index, certificate := range certificates {
		if certificate.Id == params["id"] && certificate.OwnerId == authorization {
			certificates = append(certificates[:index], certificates[index+1:]...)

			var user User
			DecodeFromJson(w, req, &user)

			user.Id = authorization
			certificate.Transfer = &Transfer{To: user.Email, Status: Pending}
			certificates = append(certificates, certificate)

			EncodeToJson(w, certificate)
			return
		}
	}
}

func acceptTransfer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	authorization := req.Header.Get("Authorization")

	for index, certificate := range certificates {
		if certificate.Id == params["id"] && certificate.OwnerId == authorization {
			certificates = append(certificates[:index], certificates[index+1:]...)

			var user User
			DecodeFromJson(w, req, &user)

			certificate.Transfer = &Transfer{}
			certificate.OwnerId = user.Id
			certificates = append(certificates, certificate)

			EncodeToJson(w, Transfer{}.TransferComplete(user.Email))
			return
		}
	}
}

func cancelTransfer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	authorization := req.Header.Get("Authorization")

	for _, certificate := range certificates {
		if certificate.Id == params["id"] && certificate.OwnerId == authorization {
			certificate.Transfer = &Transfer{}

			EncodeToJson(w, Transfer{}.TransferDeclined())
			return
		}
	}
}

func getCertificates(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	EncodeToJson(w, certificates)
}

func getCertificate(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)

	for _, certificate := range certificates {
		if certificate.Id == params["id"] {
			EncodeToJson(w, certificate)
			return
		}
	}

	EncodeToJson(w, &Certificate{})
}

func main() {
	certificates = MockCertificates(certificates)
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
	router.HandleFunc("/certificates/{id}/transfers", acceptTransfer).Methods("PUT")
	router.HandleFunc("/certificates/{id}/transfers", cancelTransfer).Methods("PATCH")

	err := http.ListenAndServe(":8000", router)

	if err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}
