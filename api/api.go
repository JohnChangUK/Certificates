package api

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

// This API retrieves all the user's certificates by matching the certificate's owner Id
func GetUserCertificate(w http.ResponseWriter, req *http.Request) {
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

// Creates a new certificate
func CreateCertificate(w http.ResponseWriter, req *http.Request) {
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

// Updates the certificate by certificate Id
func UpdateCertificate(w http.ResponseWriter, req *http.Request) {
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

// Deletes the certificate by by matching the certificate Id with the Params[Id]
func DeleteCertificate(w http.ResponseWriter, req *http.Request) {
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

/**
* This API is called when User A creates a Transfer request
* The Transfer status is changed to PENDING along with User B's email
 */
func CreateTransfer(w http.ResponseWriter, req *http.Request) {
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

/**
* This API is called when User B accepts the Transfer request.
* The certificate owner Id is changed to User B's Id
* Transfer struct with a COMPLETE status is sent back as the payload
 */
func AcceptTransfer(w http.ResponseWriter, req *http.Request) {
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

			EncodeToJson(w, *Transfer{}.TransferComplete(user.Email))
			return
		}
	}
}

/**
* This API is called when User B declines the Transfer request.
* Original contents of the certificate is kept the same.
* Transfer struct with a DECLINED status is sent back as the payload
 */
func CancelTransfer(w http.ResponseWriter, req *http.Request) {
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

// Retrieves all certificates
func GetCertificates(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	EncodeToJson(w, certificates)
}

// Retrieves a certificate by matching the certificate Id with the Params[Id]
func GetCertificate(w http.ResponseWriter, req *http.Request) {
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

func StartHttpServer() {
	certificates = MockCertificates(certificates)
	router := mux.NewRouter()

	router.HandleFunc("/certificates", GetCertificates).Methods("GET")
	router.HandleFunc("/certificates/{id}", GetCertificate).Methods("GET")
	router.HandleFunc("/users/{userId}/certificates", GetUserCertificate).Methods("GET")
	router.HandleFunc("/certificates", CreateCertificate).Methods("POST")
	router.HandleFunc("/certificates/{id}", UpdateCertificate).Methods("PUT")
	router.HandleFunc("/certificates/{id}", DeleteCertificate).Methods("DELETE")
	router.HandleFunc("/certificates/{id}/transfers", CreateTransfer).Methods("POST")
	router.HandleFunc("/certificates/{id}/transfers", AcceptTransfer).Methods("PUT")
	router.HandleFunc("/certificates/{id}/transfers", CancelTransfer).Methods("PATCH")

	err := http.ListenAndServe(":8000", router)

	if err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}
