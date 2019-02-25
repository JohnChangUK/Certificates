package api

import (
	"bytes"
	"encoding/json"
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

const GetCertificatesUrl = "localhost:8000/certificates"
const GetUserJohnCertificatesUrl = "/users/John/certificates"

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
	responseBodyBytes := getResponseBody(t, req, err, GetUserCertificatesHandler)

	var payload []Certificate
	jsonUnmarshalErr := json.Unmarshal([]byte(string(bytes.TrimSpace(responseBodyBytes))), &payload)
	if jsonUnmarshalErr != nil {
		t.Fatal(jsonUnmarshalErr)
	}

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
	responseBodyBytes := getResponseBody(t, req, err, CreateCertificateHandler)

	var payload Certificate
	jsonUnmarshalErr := json.Unmarshal([]byte(string(bytes.TrimSpace(responseBodyBytes))), &payload)
	if jsonUnmarshalErr != nil {
		t.Fatal(jsonUnmarshalErr)
	}

	assert.True(t, reflect.DeepEqual("New Certificate", payload.Title))
	assert.True(t, reflect.DeepEqual(2019, payload.Year))
	assert.True(t, reflect.DeepEqual("Art note", payload.Note))
}

func TestUpdateCertificateHandler(t *testing.T) {
	req, err := http.NewRequest("PUT", GetCertificatesUrl, strings.NewReader(
		`{"Title": "Updated Certificate", "Year": 3000, "Note": "Updated note"}`))
	if err != nil {
		t.Fatalf("Not able to create request: %v", err)
	}

	req = addAuthAndSetPathVariables(req, "John", "id", "1")
	responseBodyBytes := getResponseBody(t, req, err, UpdateCertificateHandler)

	var payload Certificate
	jsonUnmarshalErr := json.Unmarshal([]byte(string(bytes.TrimSpace(responseBodyBytes))), &payload)
	if jsonUnmarshalErr != nil {
		t.Fatal(jsonUnmarshalErr)
	}

	assert.True(t, reflect.DeepEqual("Updated Certificate", payload.Title))
	assert.True(t, reflect.DeepEqual(3000, payload.Year))
	assert.True(t, reflect.DeepEqual("Updated note", payload.Note))
}

func TestGetCertificateHandler(t *testing.T) {
	req, err := http.NewRequest("GET", GetCertificatesUrl, nil)
	if err != nil {
		t.Fatalf("Not able to create request: %v", err)
	}

	req = addAuthAndSetPathVariables(req, "John", "id", "1")
	responseBodyBytes := getResponseBody(t, req, err, GetCertificateHandler)

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

func addAuthAndSetPathVariables(req *http.Request, auth string, pathKey string, pathValue string) *http.Request {
	req.Header.Add("Authorization", auth)
	req = mux.SetURLVars(req, map[string]string{pathKey: pathValue})
	return req
}

func getResponseBody(t *testing.T, req *http.Request, err error, f func(http.ResponseWriter,
	*http.Request)) []byte {
	recorder := startMockHttpServer(req, f)
	response := recorder.Result()
	checkIfStatusCodeIs200(response, t)

	responseBodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("Could not read response %v", err)
	}

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
