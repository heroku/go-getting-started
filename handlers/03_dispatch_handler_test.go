package handlers

import (
	"bytes"
	"golang-fifa-world-cup-web-service/data"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCorrectHTTPGetMethodDispatch(t *testing.T) {
	setup()

	req, _ := http.NewRequest("GET", "/winners", nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(WinnersHandler)
	handler.ServeHTTP(rr, req)

	if body := rr.Body.String(); body == "" {
		t.Error("Did not properly dispatch HTTP GET")
	}
}

func TestCorrectHTTPPostMethodDispatch(t *testing.T) {
	setup()

	var jsonStr = []byte(`{"country":"Croatia", "year": 2030}`)
	req, _ := http.NewRequest("POST", "/winners", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-ACCESS-TOKEN", data.AccessToken)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(WinnersHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Error("Did not properly dispatch HTTP POST", status)
	}
}

func TestCorrectHTTPUnsupportedMethodDispatch(t *testing.T) {
	setup()

	req, _ := http.NewRequest("PUT", "/winners", nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(WinnersHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Error("Did not properly catch unsupported HTTP methods", status)
	}
}
