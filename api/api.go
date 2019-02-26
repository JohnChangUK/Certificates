package api

import (
	. "github.com/JohnChangUK/Certificates/constants"
	. "github.com/JohnChangUK/Certificates/model"
	. "github.com/JohnChangUK/Certificates/utils"
	"github.com/gorilla/mux"
	"github.com/rs/xid"
	"log"
	"net/http"
	"time"
)

var certificates []Certificate

// This API retrieves all the user's certificates by matching the certificate's owner Id
func GetUserCertificatesHandler(w http.ResponseWriter, req *http.Request) {
	params, authorization := GetParamsAndSetContentTypeToJson(w, req)
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
func CreateCertificateHandler(w http.ResponseWriter, req *http.Request) {
	_, authorization := GetParamsAndSetContentTypeToJson(w, req)

	var certificate Certificate
	DecodeFromJson(w, req, &certificate)

	certificate.CreatedAt = time.Now()
	certificate.OwnerId = authorization
	certificate.Id = xid.New().String()
	certificate.Transfer = &Transfer{}
	certificates = append(certificates, certificate)

	EncodeToJson(w, certificate)
}

// Updates the certificate by certificate Id
func UpdateCertificateHandler(w http.ResponseWriter, req *http.Request) {
	params, authorization := GetParamsAndSetContentTypeToJson(w, req)

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
func DeleteCertificateHandler(w http.ResponseWriter, req *http.Request) {
	params, authorization := GetParamsAndSetContentTypeToJson(w, req)

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
func CreateTransferHandler(w http.ResponseWriter, req *http.Request) {
	params, authorization := GetParamsAndSetContentTypeToJson(w, req)

	for index, certificate := range certificates {
		if certificate.Id == params["id"] && certificate.OwnerId == authorization {
			certificates = append(certificates[:index], certificates[index+1:]...)

			var user User
			DecodeFromJson(w, req, &user)

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
func AcceptTransferHandler(w http.ResponseWriter, req *http.Request) {
	params, authorization := GetParamsAndSetContentTypeToJson(w, req)

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
func DeclineTransferHandler(w http.ResponseWriter, req *http.Request) {
	params, authorization := GetParamsAndSetContentTypeToJson(w, req)

	for _, certificate := range certificates {
		if certificate.Id == params["id"] && certificate.OwnerId == authorization {
			*certificate.Transfer = Transfer{}

			EncodeToJson(w, Transfer{}.TransferDeclined())
			return
		}
	}
}

// Retrieves all certificates
func GetCertificatesHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	EncodeToJson(w, certificates)
}

// Retrieves a certificate by matching the certificate Id with the Params[Id]
func GetCertificateHandler(w http.ResponseWriter, req *http.Request) {
	params, _ := GetParamsAndSetContentTypeToJson(w, req)

	for _, certificate := range certificates {
		if certificate.Id == params["id"] {
			EncodeToJson(w, certificate)
			return
		}
	}

	EncodeToJson(w, &Certificate{})
}

func StartHttpServer() {
	MockCertificates(&certificates)
	router := mux.NewRouter()

	router.HandleFunc("/users/{userId}/certificates", GetUserCertificatesHandler).Methods("GET")
	router.HandleFunc("/certificates", CreateCertificateHandler).Methods("POST")
	router.HandleFunc("/certificates/{id}", UpdateCertificateHandler).Methods("PUT")
	router.HandleFunc("/certificates/{id}", DeleteCertificateHandler).Methods("DELETE")
	router.HandleFunc("/certificates/{id}/transfers", CreateTransferHandler).Methods("POST")
	router.HandleFunc("/certificates/{id}/transfers", AcceptTransferHandler).Methods("PUT")
	router.HandleFunc("/certificates/{id}/transfers", DeclineTransferHandler).Methods("PATCH")
	router.HandleFunc("/certificates", GetCertificatesHandler).Methods("GET")
	router.HandleFunc("/certificates/{id}", GetCertificateHandler).Methods("GET")

	err := http.ListenAndServe(":8000", router)

	if err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}
