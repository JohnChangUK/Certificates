package api

import (
	"bytes"
	"encoding/json"
	. "github.com/JohnChangUK/Certificates/constants"
	. "github.com/JohnChangUK/Certificates/model"
	"github.com/JohnChangUK/Certificates/utils"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

const GetCertificatesUrl = "/certificates"
const CertificatesUrlWithId = "/certificates/{id}"
const GetUserJohnCertificatesUrl = "/users/John/certificates"
const CreateTransferUrl = "/certificates/{id}/transfers"

var mockData = make([]Certificate, 1)

func init() {
	utils.AddMockCertificate(&certificates)
	mockData[0] = Certificate{Id: "1", Title: "First Certificate", OwnerId: "John", Year: 2019, Note: "Blockchain note", Transfer: &Transfer{}}
}

// This test gets all certificates with the Path Variable {userId} of "John"
func TestGetUserCertificateHandler(t *testing.T) {
	req, err := http.NewRequest("GET", GetUserJohnCertificatesUrl, nil)
	if err != nil {
		t.Fatalf("Not able to create request: %v", err)
	}

	req = addAuthAndSetPathVariables(req, "John", "userId", "John")
	responseBodyBytes := getResponseBody(t, req, GetUserCertificatesHandler)

	payload := unmarshalCertificatesArray(responseBodyBytes, t)

	assert.True(t, reflect.DeepEqual(payload, mockData))
}

// This test checks if new certificate is created when providing the Title, Year and Note
func TestCreateCertificateHandler(t *testing.T) {
	req, err := http.NewRequest("POST", GetCertificatesUrl, strings.NewReader(
		`{"title": "New Certificate", "year": 2019, "note": "Art note"}`))
	if err != nil {
		t.Fatalf("Not able to create request: %v", err)
	}

	req = addAuthAndSetPathVariables(req, "John", "userId", "John")
	payload := unmarshalCertificate(t, req, CreateCertificateHandler)

	assert.True(t, reflect.DeepEqual("New Certificate", payload.Title))
	assert.True(t, reflect.DeepEqual(2019, payload.Year))
	assert.True(t, reflect.DeepEqual("Art note", payload.Note))
}

func TestUpdateCertificateHandler(t *testing.T) {
	req, err := http.NewRequest("PUT", CertificatesUrlWithId, strings.NewReader(
		`{"Title": "Updated Certificate", "Year": 3000, "Note": "Updated note"}`))
	if err != nil {
		t.Fatalf("Not able to create request: %v", err)
	}

	req = addAuthAndSetPathVariables(req, "John", "id", "1")
	payload := unmarshalCertificate(t, req, UpdateCertificateHandler)

	assert.True(t, reflect.DeepEqual("Updated Certificate", payload.Title))
	assert.True(t, reflect.DeepEqual(3000, payload.Year))
	assert.True(t, reflect.DeepEqual("Updated note", payload.Note))
}

func TestDeleteCertificateHandler(t *testing.T) {
	// This is to reset the number of certificates in memory, aka add 1 single mock certificate
	utils.AddMockCertificate(&certificates)

	req, err := http.NewRequest("DELETE", CertificatesUrlWithId, nil)
	if err != nil {
		t.Fatalf("Not able to create request: %v", err)
	}

	req = addAuthAndSetPathVariables(req, "John", "id", "1")
	responseBodyBytes := getResponseBody(t, req, DeleteCertificateHandler)

	payload := unmarshalCertificatesArray(responseBodyBytes, t)

	assert.True(t, reflect.DeepEqual([]Certificate{}, payload))
}

func TestCreateTransferHandler(t *testing.T) {
	utils.AddMockCertificate(&certificates)

	req, err := http.NewRequest("POST", CreateTransferUrl, strings.NewReader(
		`{"id": "UserB", "email": "userb@gmail.com", "name": "User B"}`))
	if err != nil {
		t.Fatalf("Not able to create request: %v", err)
	}

	req = addAuthAndSetPathVariables(req, "John", "id", "1")
	payload := unmarshalCertificate(t, req, CreateTransferHandler)

	assert.True(t, reflect.DeepEqual("userb@gmail.com", payload.Transfer.To))
	assert.True(t, reflect.DeepEqual(Pending, payload.Transfer.Status))
}

func TestAcceptTransferHandler(t *testing.T) {
	utils.AddMockCertificate(&certificates)

	req, err := http.NewRequest("PUT", CreateTransferUrl, strings.NewReader(
		`{"id": "UserB", "email": "userb@gmail.com", "name": "User B"}`))
	if err != nil {
		t.Fatalf("Not able to create request: %v", err)
	}

	req = addAuthAndSetPathVariables(req, "John", "id", "1")
	transfer := unmarshalTransfer(t, req, AcceptTransferHandler)

	assert.True(t, reflect.DeepEqual("userb@gmail.com", transfer.To))
	assert.True(t, reflect.DeepEqual(Complete, transfer.Status))
}

func TestCancelTransferHandler(t *testing.T) {
	utils.AddMockCertificate(&certificates)

	req, err := http.NewRequest("PATCH", CreateTransferUrl, strings.NewReader(
		`{"id": "UserB", "email": "userb@gmail.com", "name": "User B"}`))
	if err != nil {
		t.Fatalf("Not able to create request: %v", err)
	}

	req = addAuthAndSetPathVariables(req, "John", "id", "1")
	transfer := unmarshalTransfer(t, req, DeclineTransferHandler)

	assert.True(t, reflect.DeepEqual(Declined, transfer.Status))
}

func TestGetCertificateHandler(t *testing.T) {
	req, err := http.NewRequest("GET", GetCertificatesUrl, nil)
	if err != nil {
		t.Fatalf("Not able to create request: %v", err)
	}

	req = addAuthAndSetPathVariables(req, "John", "id", "1")
	responseBodyBytes := getResponseBody(t, req, GetCertificateHandler)

	var payload Certificate
	jsonUnmarshalErr := json.Unmarshal([]byte(string(bytes.TrimSpace(responseBodyBytes))), &payload)
	if jsonUnmarshalErr != nil {
		t.Fatal(jsonUnmarshalErr)
	}

	var correctCertificate Certificate
	for _, certificate := range certificates {
		if "1" == certificate.Id {
			correctCertificate = certificate
		}
	}

	assert.True(t, reflect.DeepEqual(correctCertificate, payload))
}

func unmarshalCertificate(t *testing.T, req *http.Request, f func(http.ResponseWriter, *http.Request)) Certificate {
	responseBodyBytes := getResponseBody(t, req, f)
	var payload Certificate
	jsonUnmarshalErr := json.Unmarshal([]byte(string(bytes.TrimSpace(responseBodyBytes))), &payload)
	if jsonUnmarshalErr != nil {
		t.Fatal(jsonUnmarshalErr)
	}

	return payload
}

func unmarshalCertificatesArray(responseBodyBytes []byte, t *testing.T) []Certificate {
	var payload []Certificate
	jsonUnmarshalErr := json.Unmarshal([]byte(string(bytes.TrimSpace(responseBodyBytes))), &payload)
	if jsonUnmarshalErr != nil {
		t.Fatal(jsonUnmarshalErr)
	}

	return payload
}

func unmarshalTransfer(t *testing.T, req *http.Request, f func(http.ResponseWriter, *http.Request)) Transfer {
	responseBodyBytes := getResponseBody(t, req, f)
	var transfer Transfer
	jsonUnmarshalErr := json.Unmarshal([]byte(string(bytes.TrimSpace(responseBodyBytes))), &transfer)
	if jsonUnmarshalErr != nil {
		t.Fatal(jsonUnmarshalErr)
	}

	return transfer
}

func addAuthAndSetPathVariables(req *http.Request, auth string, pathKey string, pathValue string) *http.Request {
	req.Header.Add("Authorization", auth)
	req = mux.SetURLVars(req, map[string]string{pathKey: pathValue})
	return req
}

func getResponseBody(t *testing.T, req *http.Request, f func(http.ResponseWriter,
	*http.Request)) []byte {
	recorder := startMockHttpServer(req, f)
	response := recorder.Result()
	checkIfStatusCodeIs200(response, t)

	responseBodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("Could not read response %v", err)
	}

	defer response.Body.Close()
	return responseBodyBytes
}

func checkIfStatusCodeIs200(response *http.Response, t *testing.T) {
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got: %v ", response.Status)
	}
}

func startMockHttpServer(req *http.Request, f func(http.ResponseWriter, *http.Request)) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(f)
	handler.ServeHTTP(recorder, req)

	return recorder
}
