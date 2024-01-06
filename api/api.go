package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
)

func main() {
	config := LoadConfig()

	router := mux.NewRouter()

	router.HandleFunc("/api/mesures/{iata}/{start}/{end}", AirportWeatherMeasures).Methods("GET")
}

type HourMeasurement struct {
	Heure  string      `json:"heure"`
	Mesure Measurement `json:"mesure"`
}

type DayMeasurement struct {
	Day        string            `json:"day"`
	HourMesure []HourMeasurement `json:"hourMesure"`
}

type MeasuresResultat struct {
	Temperature []DayMeasurement `json:"temperature"`
	Pressure    []DayMeasurement `json:"pressure"`
	WindSpeed   []DayMeasurement `json:"windSpeed"`
}

type Measurement struct {
	IDCapteur string  `json:"idCapteur"`
	Value     float64 `json:"value"`
}


func AirportWeatherMeasures(responseWriter http.ResponseWriter, request *http.Request) {
    // on recupére la config
    configuration := LoadConfig()

    // on extrait les parametre de l'url
    parameters := mux.Vars(request)
    airportCode := parameters["iata"]

    // on parse les donnees de format time.
    const timeFormat = "2000-01-01-01"

    startTime, errStart := time.Parse(timeFormat, parameters["start"])
    endTime, errEnd := time.Parse(timeFormat, parameters["end"])
    if errStart != nil || errEnd != nil {
        http.Error(responseWriter, "Invalid time format", http.StatusBadRequest)
        return
    }

    // TODO create db connectionn (InfluxDB???)
    
    // crée les tableaux pour stocker les mesures.
    var temps, pressures, windSpeeds []DayMeasurement

    // ajuste le start time au debut de la journee 
    dayStart := startTime.Truncate(24 * time.Hour)

    //for day := adjustedStart; day.Before(endTime); day = day.Add(24 * time.Hour) {
        //weatherData := [3]DayMeasurement{}

        // mesures à recuperer.
        measurements := []string{"temperature", "pressure", "wind_speed"}

        for index, measurementType := range measurements {
            key := fmt.Sprintf("/%s/%s/%d/%02d/%02d", airportCode, measurementType, day.Year(), day.Month(), day.Day())
            weatherData[index] = DayMeasurement{
                Date: fmt.Sprintf("%02d/%02d/%d", day.Day(), day.Month(), day.Year()),
                HourlyData: RetrieveDailyData(key, startTime, endTime, day, redisConnection),
            }
        }

        temps = append(temps, weatherData[0])
        pressures = append(pressures, weatherData[1])
        windSpeeds = append(windSpeeds, weatherData[2])
    }

    weatherResult := MeasuresResult{
        Temperatures: temps,
        Pressures:    pressures,
        WindSpeeds:   windSpeeds,
    }

    // encode en JSON
    if err := json.NewEncoder(responseWriter).Encode(weatherResult); err != nil {
        http.Error(responseWriter, "Failed to encode JSON: "+err.Error(), http.StatusInternalServerError)
    }
}

