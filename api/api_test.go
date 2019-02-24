package api

import (
	"bytes"
	"encoding/json"
	. "github.com/JohnChangUK/verisart/model"
	"github.com/JohnChangUK/verisart/utils"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const GetCertificatesUrl = "localhost:8000/certificates"
const GetUserJohnCertificatesUrl = "localhost:8000/users/John/certificates"

var mockData = make([][]Certificate, 2)

func init() {
	utils.AddMockCertificate(&certificates)
	createMultiDimensionalCertificateArray("1", "First Certificate", "John", 2019,
		"Blockchain note", &Transfer{})
}

// This test gets all certificates with the User Id 'John'
func TestGetUserCertificateHandler(t *testing.T) {
	req, err := http.NewRequest("GET", GetUserJohnCertificatesUrl, nil)
	if err != nil {
		t.Fatalf("Not able to create request: %v", err)
	}

	responseBodyBytes := startServerAndGetResponse(t, req, err)

	var payload [][]Certificate
	jsonUnmarshalErr := json.Unmarshal([]byte(string(bytes.TrimSpace(responseBodyBytes))), &payload)
	if jsonUnmarshalErr != nil {
		t.Fatal(jsonUnmarshalErr)
	}

	jsonPayload, _ := json.Marshal(&payload)
	jsonMockData, _ := json.Marshal(&mockData)

	if string(jsonPayload) != string(jsonMockData) {
		t.Fatalf("Error occured... expected payload: %v", mockData)
	}
}

func TestCreateCertificateHandler(t *testing.T) {
	jsonFormat, err := json.Marshal(&mockData)
	if err != nil {
		t.Fatalf("Not able to Marshal data into bytes: %v", err)
	}

	stringData := string(jsonFormat)
	req, err := http.NewRequest("POST", GetCertificatesUrl, strings.NewReader(stringData))
	if err != nil {
		t.Fatalf("Not able to create request: %v", err)
	}

	responseBodyBytes := startServerAndGetResponse(t, req, err)

	var payload [][]Certificate
	jsonUnmarshalErr := json.Unmarshal([]byte(string(bytes.TrimSpace(responseBodyBytes))), &payload)
	if jsonUnmarshalErr != nil {
		t.Fatal(jsonUnmarshalErr)
	}

	jsonPayload, _ := json.Marshal(&payload)
	jsonMockData, _ := json.Marshal(&mockData)

	if string(jsonPayload) != string(jsonMockData) {
		t.Fatalf("Error occured... expected payload: %v", mockData)
	}
}

func TestGetCertificateHandler(t *testing.T) {
	req, err := http.NewRequest("GET", GetCertificatesUrl, nil)
	if err != nil {
		t.Fatalf("Not able to create request: %v", err)
	}

	responseBodyBytes := startServerAndGetResponse(t, req, err)

	var payload [][]Certificate
	jsonUnmarshalErr := json.Unmarshal([]byte(string(bytes.TrimSpace(responseBodyBytes))), &payload)
	if jsonUnmarshalErr != nil {
		t.Fatal(jsonUnmarshalErr)
	}

	jsonPayload, _ := json.Marshal(&payload)
	jsonMockData, _ := json.Marshal(&mockData)

	if string(jsonPayload) != string(jsonMockData) {
		t.Fatalf("Error occured... expected payload: %v", mockData)
	}
}

func startServerAndGetResponse(t *testing.T, req *http.Request, err error) []byte {
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(GetCertificatesHandler)
	handler.ServeHTTP(recorder, req)
	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got: %v ", response.Status)
	}

	responseBodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("Could not read response %v", err)
	}

	return responseBodyBytes
}

/**
* Creates a multi dimensional array of Certificates
* Payload of HTTP Handlers return a structure of [][]Certificate
 */
func createMultiDimensionalCertificateArray(id string, title string, ownerId string, year int, note string, transfer *Transfer) {
	mockData = make([][]Certificate, 1)

	for i := 0; i < 1; i++ {
		mockData[i] = make([]Certificate, 0, 2)
		vector := make([]Certificate, 2)
		for j := 0; j < 1; j++ {
			vector[j] = Certificate{Id: id, Title: title, OwnerId: ownerId, Year: year, Note: note, Transfer: transfer}
			mockData[i] = append(mockData[i], vector[j])
		}
	}
}
