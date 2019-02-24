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

var mockData = make([][]Certificate, 2)

func init() {
	utils.AddMockCertificate(&certificates)
	createMultiDimensionalCertificateArray("1", "First Certificate", "John", 2019,
		"Blockchain note", &Transfer{})
}

// This test gets all certificates with the Path Variable {userId} of "John"
func TestGetUserCertificateHandler(t *testing.T) {
	req, err := http.NewRequest("GET", GetUserJohnCertificatesUrl, nil)
	if err != nil {
		t.Fatalf("Not able to create request: %v", err)
	}

	req = addAuthAndSetPathVariables(req, "John", "userId", "John")
	responseBodyBytes := getResponseBody(t, req, err, GetUserCertificatesHandler)

	var payload [][]Certificate
	jsonUnmarshalErr := json.Unmarshal([]byte(string(bytes.TrimSpace(responseBodyBytes))), &payload)
	if jsonUnmarshalErr != nil {
		t.Fatal(jsonUnmarshalErr)
	}

	assert.True(t, reflect.DeepEqual(payload, mockData))
}

// This test checks if new certificate is created when providing the Title, Year and Note
func TestCreateCertificateHandler(t *testing.T) {
	req, err := http.NewRequest("POST", GetCertificatesUrl, strings.NewReader(
		`{"Title": "New Certificate", "Year": 2019, "Note": "Art note"}`))
	if err != nil {
		t.Fatalf("Not able to create request: %v", err)
	}

	req = addAuthAndSetPathVariables(req, "John", "userId", "John")
	responseBodyBytes := getResponseBody(t, req, err, CreateCertificateHandler)

	var payload []Certificate
	jsonUnmarshalErr := json.Unmarshal([]byte(string(bytes.TrimSpace(responseBodyBytes))), &payload)
	if jsonUnmarshalErr != nil {
		t.Fatal(jsonUnmarshalErr)
	}

	assert.True(t, reflect.DeepEqual("New Certificate", payload[0].Title))
	assert.True(t, reflect.DeepEqual(2019, payload[0].Year))
	assert.True(t, reflect.DeepEqual("Art note", payload[0].Note))
}

func TestUpdateCertificateHandler(t *testing.T) {
	var updateCertificatePayload [1][1]Certificate
	updateCertificatePayload[0][0] = Certificate{Title: "First Certificate",
		OwnerId: "John", Year: 2019, Note: "Blockchain note",
		Transfer: &Transfer{}}

	mockData[0][0] = Certificate{Id: "1", Title: "First Certificate",
		OwnerId: "John", Year: 2019, Note: "Blockchain note",
		Transfer: &Transfer{}}

	jsonFormat, err := json.Marshal(&updateCertificatePayload)
	if err != nil {
		t.Fatalf("Not able to Marshal data into bytes: %v", err)
	}

	stringData := string(jsonFormat)
	req, err := http.NewRequest("PUT", GetCertificatesUrl, strings.NewReader(stringData))
	if err != nil {
		t.Fatalf("Not able to create request: %v", err)
	}

	responseBodyBytes := getResponseBody(t, req, err, UpdateCertificateHandler)

	var payload [][]Certificate
	jsonUnmarshalErr := json.Unmarshal([]byte(string(bytes.TrimSpace(responseBodyBytes))), &payload)
	if jsonUnmarshalErr != nil {
		t.Fatal(jsonUnmarshalErr)
	}

	if !reflect.DeepEqual(payload, mockData) {
		t.Fatalf("Error occured... expected payload: %v", mockData)
	}
}

func addAuthAndSetPathVariables(req *http.Request, auth string, pathKey string, pathValue string) *http.Request {
	req.Header.Add("Authorization", auth)
	req = mux.SetURLVars(req, map[string]string{pathKey: pathValue})
	return req
}

func TestGetCertificateHandler(t *testing.T) {
	req, err := http.NewRequest("GET", GetCertificatesUrl, nil)
	if err != nil {
		t.Fatalf("Not able to create request: %v", err)
	}

	responseBodyBytes := getResponseBody(t, req, err, GetCertificateHandler)

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
