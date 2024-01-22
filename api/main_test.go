package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetMeasurements(t *testing.T) {
	// Example test data
	iata := "NTE"
	mesure := "temperature"
	start := "2023-01-20-10"
	end := "2024-01-20-12"

	path := "/api/mesure/" + mesure + "/" + iata + "/" + start + "/" + end
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetMeasurements)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Add checks for the response body here...
}

func TestGetMean(t *testing.T) {
	iata := "NTE"
	mesure := "temperature"
	start := "2024-01-09"

	path := "/api/mesure/" + mesure + "/" + iata + "/" + start
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetMean)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Add checks for the response body here...
}

func TestGetAllMeans(t *testing.T) {
	iata := "NTE"
	start := "2024-01-09"

	path := "/api/allMeans/" + iata + "/" + start
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAllMeans)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Add checks for the response body here...
}
