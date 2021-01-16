package handlers

import (
	"bytes"
	"encoding/json"
	"golang-fifa-world-cup-web-service/data"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddNewWinnerHandlerReturnsUnauthorizedForInvalidAccessToken(t *testing.T) {
	setup()

	req, _ := http.NewRequest("POST", "/winners", nil)
	req.Header.Set("X-ACCESS-TOKEN", data.AccessToken+"bla")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AddNewWinner)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Error("Did not return status 401 - Unauthorized for invalid Access Token")
	}
}

func TestAddNewWinnerHandlerReturnsCreatedForValidAccessToken(t *testing.T) {
	var jsonStr = []byte(`{"country":"Croatia", "year": 2030}`)
	req, _ := http.NewRequest("POST", "/winners", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-ACCESS-TOKEN", data.AccessToken)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AddNewWinner)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Error("Did not return status 201 - Created for valid Access Token")
	}
}

func TestAddNewWinnerHandlerAddsNewWinnerWithValidData(t *testing.T) {
	setup()

	var jsonStr = []byte(`{"country":"Croatia", "year": 2030}`)
	req, _ := http.NewRequest("POST", "/winners", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-ACCESS-TOKEN", data.AccessToken)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AddNewWinner)
	handler.ServeHTTP(rr, req)

	allWinners, _ := data.ListAllJSON()
	var winners data.Winners
	json.Unmarshal([]byte(allWinners), &winners)

	if len(winners.Winners) != 22 {
		t.Error("Did not properly add new winner to the list")
	}
}

func TestAddNewWinnerHandlerReturnsUnprocessableEntityForEmptyPayload(t *testing.T) {
	setup()

	// Invalid because empty
	var jsonStr = []byte(``)
	req, _ := http.NewRequest("POST", "/winners", bytes.NewBuffer(jsonStr))
	req.Header.Set("X-ACCESS-TOKEN", data.AccessToken)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AddNewWinner)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Error("Did not properly validate winner payload")
	}
}

func TestAddNewWinnerHandlerDoesNotAddInvalidNewWinner(t *testing.T) {
	setup()

	// Invalid entry because year is in the past.
	var jsonStr = []byte(`{"country":"Croatia", "year": 1984}`)
	req, _ := http.NewRequest("POST", "/winners", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-ACCESS-TOKEN", data.AccessToken)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AddNewWinner)
	handler.ServeHTTP(rr, req)

	allWinners, _ := data.ListAllJSON()
	var winners data.Winners
	json.Unmarshal([]byte(allWinners), &winners)

	if rr.Code == http.StatusOK || len(winners.Winners) != 21 {
		t.Error("Added invalid winner to list")
	}
}
