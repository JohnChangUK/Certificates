package api

import (
	"bytes"
	"encoding/json"
	. "github.com/JohnChangUK/Certificates/model"
	"github.com/JohnChangUK/Certificates/utils"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
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

	responseBodyBytes := getResponseBody(t, req, err)

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
	var createCertificatePayload [1][1]Certificate
	createCertificatePayload[0][0] = Certificate{Id: "1", Title: "First Certificate",
		OwnerId: "John", Year: 2019, Note: "Blockchain note",
		Transfer: &Transfer{}}

	mockData[0][0] = Certificate{Id: "1", Title: "First Certificate",
		OwnerId: "John", Year: 2019, Note: "Blockchain note",
		Transfer: &Transfer{}}

	jsonFormat, err := json.Marshal(&createCertificatePayload)
	if err != nil {
		t.Fatalf("Not able to Marshal data into bytes: %v", err)
	}

	stringData := string(jsonFormat)
	req, err := http.NewRequest("POST", GetCertificatesUrl, strings.NewReader(stringData))
	if err != nil {
		t.Fatalf("Not able to create request: %v", err)
	}

	responseBodyBytes := getResponseBody(t, req, err)

	var payload [][]Certificate
	jsonUnmarshalErr := json.Unmarshal([]byte(string(bytes.TrimSpace(responseBodyBytes))), &payload)
	if jsonUnmarshalErr != nil {
		t.Fatal(jsonUnmarshalErr)
	}

	if !reflect.DeepEqual(payload, mockData) {
		t.Fatalf("Error occured... expected payload: %v", mockData)
	}
}

func TestGetCertificateHandler(t *testing.T) {
	req, err := http.NewRequest("GET", GetCertificatesUrl, nil)
	if err != nil {
		t.Fatalf("Not able to create request: %v", err)
	}

	responseBodyBytes := getResponseBody(t, req, err)

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

func getResponseBody(t *testing.T, req *http.Request, err error) []byte {
	recorder := startMockHttpServer(req)
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

func startMockHttpServer(req *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(GetCertificatesHandler)
	handler.ServeHTTP(recorder, req)

	return recorder
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
