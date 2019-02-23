package main

import (
	"encoding/json"
	. "github.com/JohnChangUK/verisart/model"
	. "github.com/JohnChangUK/verisart/utils"
	"github.com/gorilla/mux"
	"github.com/rs/xid"
	"log"
	"net/http"
	"time"
)

var certificates []Certificate
var transfers = make(map[string]Transfer)
var certificatesTranfer = make(map[*Certificate]map[Transfer]User)
var userCertificates = make(map[User]Certificate)

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

func getUserCertificate(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	userId := req.Header.Get("Authorization")
	var usersCertificates []Certificate

	for _, certificate := range certificates {
		if userId == params["userId"] && userId == certificate.OwnerId {
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
	err := json.NewDecoder(req.Body).Decode(&certificate)
	req.Body.Close()

	if err != nil {
		InvalidBody(w, req, err)
		return
	}

	certificate.CreatedAt = time.Now()
	certificate.OwnerId = userId
	certificate.Id = xid.New().String()
	certificate.Transfer = &Transfer{}
	certificates = append(certificates, certificate)

	EncodeToJson(w, certificate)
}

func updateCertificate(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userId := req.Header.Get("Authorization")
	params := mux.Vars(req)

	for index, certificate := range certificates {
		if certificate.Id == params["id"] {
			dateCreated := certificate.CreatedAt
			certificates = append(certificates[:index], certificates[index+1:]...)
			var certificate Certificate
			_ = json.NewDecoder(req.Body).Decode(&certificate)
			certificate.OwnerId = userId
			certificate.Id = params["id"]
			certificate.CreatedAt = dateCreated
			certificates = append(certificates, certificate)

			EncodeToJson(w, certificate)
			return
		}
	}
}

func deleteCertificate(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)

	for index, certificate := range certificates {
		if certificate.Id == params["id"] {
			certificates = append(certificates[:index], certificates[index+1:]...)
			break
		}
	}

	EncodeToJson(w, certificates)
}

func createTransfer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	userId := req.Header.Get("Authorization")

	for index, certificate := range certificates {
		if certificate.Id == params["id"] && certificate.OwnerId == userId {
			certificates = append(certificates[:index], certificates[index+1:]...)
			//var certificate Certificate
			var user User
			// Create a new Certificate with new User Details
			// Only changed when other person ACCEPTS the transfer
			_ = json.NewDecoder(req.Body).Decode(&user)
			user.Id = userId
			transfers[user.Name] = Transfer{To: user.Email, Status: "TRANSFER_PENDING"}
			certificate.Transfer = &Transfer{To: user.Email, Status: "TRANSFER_PENDING"}
			//certificatesTranfer[&certificate] = transfers
			certificates = append(certificates, certificate)

			EncodeToJson(w, certificate)
			return
		}
	}
}

func acceptTransfer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	userId := req.Header.Get("Authorization")

	for index, certificate := range certificates {
		if certificate.Id == params["id"] {
			certificates = append(certificates[:index], certificates[index+1:]...)
			var user User
			_ = json.NewDecoder(req.Body).Decode(&user)
			user.Id = userId
			certificate.Transfer = &Transfer{To: user.Email, Status: "TRANSFER_COMPLETE"}
			certificate.OwnerId = user.Id
			certificates = append(certificates, certificate)
			EncodeToJson(w, certificate)
			return
		}
	}
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

	err := http.ListenAndServe(":8000", router)
	if err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}
