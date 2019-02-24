package api

import (
	"github.com/JohnChangUK/verisart/model"
	"github.com/JohnChangUK/verisart/utils"
	"net/http"
	"net/http/httptest"
	"testing"
)

const certificatesEndpoint = "/certificates"

var mockCertificates = utils.MockCertificates(make([]model.Certificate, 0))

func runGetRequest(t *testing.T, url string, handlerFunc http.HandlerFunc) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()

	handler := http.HandlerFunc(handlerFunc)
	handler.ServeHTTP(responseRecorder, req)
}

func GetAllCertificatesTest(t *testing.T) {
	runGetRequest(t, certificatesEndpoint, GetCertificates)
}
