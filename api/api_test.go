package api

import (
	"bytes"
	"encoding/json"
	"github.com/JohnChangUK/verisart/model"
	"github.com/JohnChangUK/verisart/utils"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

const certificatesEndpoint = "localhost:8000/certificates"

var mockData = make([][]model.Certificate, 1)

func Before() {
	utils.AddMockCertificate(&certificates)
	createMultiDimensionalCertificateArray()
}

func TestGetCertificateHandler(t *testing.T) {
	runGetRequest(t, certificatesEndpoint, GetCertificatesHandler)
}

func runGetRequest(t *testing.T, url string, handlerFunc http.HandlerFunc) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("Not able to create request: %v", err)
	}

	recorder := httptest.NewRecorder()

	Before()
	handler := http.HandlerFunc(handlerFunc)
	handler.ServeHTTP(recorder, req)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got: %v ", response.Status)
	}

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("Could not read response %v", err)
	}

	payload := string(bytes.TrimSpace(b))
	log.Print(payload)

	var msg [][]model.Certificate
	err2 := json.Unmarshal([]byte(string(payload)), &msg)
	if err2 != nil {
		t.Fatal(err2)
	}

	jsonMsg, _ := json.Marshal(&msg)
	jsonMockData, _ := json.Marshal(&mockData)

	if string(jsonMsg) != string(jsonMockData) {
		t.Fatalf("Error occured... expected payload: %v", mockData)
	}
}

/**
* Creates a multi dimensional array of Certificates
* Payload of HTTP Handlers return a structure of [][]Certificate
 */
func createMultiDimensionalCertificateArray() {
	mockData = make([][]model.Certificate, 1)
	for i := 0; i < 1; i++ {
		mockData[i] = make([]model.Certificate, 0, 1)
		vector := make([]model.Certificate, 1)
		for j := 0; j < 1; j++ {
			vector[j] = model.Certificate{Id: "1", Title: "First Certificate",
				OwnerId: "John", Year: 2019, Note: "Blockchain",
				Transfer: &model.Transfer{}}
			mockData[i] = append(mockData[i], vector[j])
		}
	}
}
