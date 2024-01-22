package main

import (
	_ "aymane.com/main/api/docs"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"net/url"
	"time"
)

type MeasurementResponse struct {
	Timestamp string  `json:"timestamp"`
	Value     float64 `json:"value"`
}

type MeanMeasurementResponse struct {
	Measurement string  `json:"measurement"`
	MeanValue   float64 `json:"mean_value"`
}

func main() {
	fmt.Println("server running")
	r := mux.NewRouter()

	r.HandleFunc("/api/mesure/{mesures}/{iata}/{start}/{end}", GetMeasurements).Methods("GET")

	r.HandleFunc("/api/mesure/{mesures}/{iata}/{start}", GetMean).Methods("GET")

	r.HandleFunc("/api/allMeans/{iata}/{start}", GetAllMeans).Methods("GET")

	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler())

	log.Fatal(http.ListenAndServe(":8080", r))
}

// GetMeasurements godoc
// @Summary Get Measurements
// @Description Retrieve measurements for a given sensor, location, and time range
// @Tags measurements
// @Accept  json
// @Produce  json
// @Param   mesures path string true "Type of measurement"
// @Param   iata path string true "IATA code for the location"
// @Param   start path string true "Start time in YYYY-MM-DD-HH format"
// @Param   end path string true "End time in YYYY-MM-DD-HH format"
// @Success 200 {array} MeasurementResponse
// @Failure 500 {object} object "Internal Server Error"
// @Router /api/mesure/{mesures}/{iata}/{start}/{end} [get]
func GetMeasurements(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	var measurements []MeasurementResponse

	vars := mux.Vars(r)
	iata := vars["iata"]
	mesures := vars["mesures"]
	const layout = "2006-01-02-15"

	decodeStartDate, _ := url.QueryUnescape(vars["start"])
	decodeEndDate, _ := url.QueryUnescape(vars["end"])

	start, _ := time.Parse(layout, decodeStartDate)
	end, _ := time.Parse(layout, decodeEndDate)

	fmt.Println(start, end, iata)
	formattedDateStart := start.Format(time.RFC3339)
	formattedDateEnd := end.Format(time.RFC3339)
	fmt.Println("Formatted Date:", formattedDateStart, formattedDateEnd)

	client := influxdb2.NewClient("http://158.178.194.137:8086", "xqBlL61ya6BtT6P_F3TrLLfbTOwXnUF9ZTX2dcZAOBcHOtA7QCNuFo8HphUsYJgTLdXZozjs22P-EtVxMe6N1g==")
	queryAPI := client.QueryAPI("IMT")

	query := fmt.Sprintf(`from(bucket: "projetGo") `+
		`|> range(start: %s, stop: %s) `+
		`|> filter(fn: (r) => r._measurement == "%s" and r.ville == "%s")`, formattedDateStart, formattedDateEnd, mesures, iata)

	fmt.Println(query)
	// Query the database
	result, err := queryAPI.Query(context.Background(), query)
	fmt.Println("bonjour")
	if err == nil {

		// Iterate through the result set
		for result.Next() {
			fmt.Printf("Average: %v\n", result.Record().Values())
			measurement := MeasurementResponse{
				Timestamp: result.Record().Time().Format(time.RFC3339),
				Value:     result.Record().Value().(float64), // Type assert based on your actual data type
			}
			measurements = append(measurements, measurement)
		}

		if result.Err() != nil {
			fmt.Printf("Query parsing error: %s\n", result.Err().Error())
			http.Error(w, result.Err().Error(), http.StatusInternalServerError)
			return

		}
	} else {
		panic(err)
	}

	// Marshal the response into JSON
	jsonResponse, err := json.Marshal(measurements)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set Content-Type header and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)

	// Close the client
	client.Close()

}

// GetMean godoc
// @Summary Get Mean Measurement
// @Description Retrieve the mean measurement for a given sensor and location on a specific date
// @Tags measurements
// @Accept  json
// @Produce  json
// @Param   mesures path string true "Type of measurement"
// @Param   iata path string true "IATA code for the location"
// @Param   start path string true "Date in YYYY-MM-DD format"
// @Success 200 {object} MeasurementResponse
// @Failure 400 {object} object "Bad Request"
// @Failure 500 {object} object "Internal Server Error"
// @Router /api/mesure/{mesures}/{iata}/{start} [get]
func GetMean(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	fmt.Printf("api-2")

	vars := mux.Vars(r)
	iata := vars["iata"]
	mesures := vars["mesures"]
	const layout = "2006-01-02"

	decodeDate, _ := url.QueryUnescape(vars["start"])

	start, err := time.Parse(layout, decodeDate)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing start time: %v", err), http.StatusBadRequest)
		return
	}

	end := start.Add(24 * time.Hour)

	formattedDateStart := start.Format(time.RFC3339)
	formattedDateEnd := end.Format(time.RFC3339)

	client := influxdb2.NewClient("http://158.178.194.137:8086", "xqBlL61ya6BtT6P_F3TrLLfbTOwXnUF9ZTX2dcZAOBcHOtA7QCNuFo8HphUsYJgTLdXZozjs22P-EtVxMe6N1g==")
	queryAPI := client.QueryAPI("IMT")

	query := fmt.Sprintf(`from(bucket: "projetGo") `+
		`|> range(start: %s, stop: %s) `+
		`|> filter(fn: (r) => r._measurement == "%s" and r.ville == "%s") `+
		`|> mean()`, formattedDateStart, formattedDateEnd, mesures, iata)

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		http.Error(w, fmt.Sprintf("Query execution error: %v", err), http.StatusInternalServerError)
		return
	}

	var measurement MeasurementResponse
	for result.Next() {
		measurement = MeasurementResponse{
			Timestamp: formattedDateStart,
			Value:     result.Record().Value().(float64),
		}

	}
	if result.Err() != nil {
		http.Error(w, fmt.Sprintf("Query parsing error: %s\n", result.Err().Error()), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(measurement)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
	client.Close()
}

// GetAllMeans godoc
// @Summary Get All Mean Measurements
// @Description Retrieve mean measurements for all sensors at a given location on a specific date
// @Tags measurements
// @Accept  json
// @Produce  json
// @Param   iata path string true "IATA code for the location"
// @Param   start path string true "Date in YYYY-MM-DD format"
// @Success 200 {array} MeanMeasurementResponse
// @Failure 400 {object} object "Bad Request"
// @Failure 500 {object} object "Internal Server Error"
// @Router /api/allMeans/{iata}/{start} [get]
func GetAllMeans(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	iata := vars["iata"]
	const layout = "2006-01-02"

	decodeDate, _ := url.QueryUnescape(vars["start"])

	fmt.Println("start: " + decodeDate)
	fmt.Println("decodedStart: " + vars["start"])
	start, err := time.Parse(layout, decodeDate)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing start time: %v", err), http.StatusBadRequest)
		return
	}

	end := start.Add(24 * time.Hour)
	formattedDateStart := start.Format(time.RFC3339)
	formattedDateEnd := end.Format(time.RFC3339)

	client := influxdb2.NewClient("http://158.178.194.137:8086", "xqBlL61ya6BtT6P_F3TrLLfbTOwXnUF9ZTX2dcZAOBcHOtA7QCNuFo8HphUsYJgTLdXZozjs22P-EtVxMe6N1g==")
	queryAPI := client.QueryAPI("IMT")

	var results []MeanMeasurementResponse

	// Define measurements
	measurements := []string{"vent", "temperature", "pression"}

	// Perform queries for each measurement
	for _, mesure := range measurements {
		query := fmt.Sprintf(`from(bucket: "projetGo") `+
			`|> range(start: %s, stop: %s) `+
			`|> filter(fn: (r) => r._measurement == "%s" and r.ville == "%s") `+
			`|> mean()`, formattedDateStart, formattedDateEnd, mesure, iata)

		result, err := queryAPI.Query(context.Background(), query)
		if err != nil {
			http.Error(w, fmt.Sprintf("Query execution error for %s: %v", mesure, err), http.StatusInternalServerError)
			return
		}

		if result.Next() {
			value, ok := result.Record().Value().(float64)
			if !ok {
				http.Error(w, fmt.Sprintf("Error processing result for %s", mesure), http.StatusInternalServerError)
				return
			}
			results = append(results, MeanMeasurementResponse{
				Measurement: mesure,
				MeanValue:   value,
			})
		}
	}

	jsonResponse, err := json.Marshal(results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
	client.Close()
}
