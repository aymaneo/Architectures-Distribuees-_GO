package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetMeasurements(t *testing.T) {

	mesures := "vent"
	iata := "NTE"
	start := "2024-01-09-00"
	end := "2024-01-10-00"

	// Create a request URL
	path := "/api/mesure/" + mesures + "/" + iata + "/" + start + "/" + end
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Create a new mux router
	r := mux.NewRouter()
	r.HandleFunc("/api/mesure/{mesures}/{iata}/{start}/{end}", GetMeasurements)

	// Simulate a request
	r.ServeHTTP(rr, req)

	// Check if  status code is what we expect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Deserialize the responsee into a MeasurementResponse
	var measurements []MeasurementResponse
	err = json.Unmarshal(rr.Body.Bytes(), &measurements)
	if err != nil {
		t.Fatal("Could not parse json:", err)
	}

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

	r := mux.NewRouter()
	r.HandleFunc("/api/allMeans/{iata}/{start}", GetAllMeans)

	r.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expectedBody := `[{"measurement":"vent","mean_value":24.036577434564993}]`

	if rr.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expectedBody)
	}
}
