package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// SensorData représente les données envoyées par un capteur.
type SensorData struct {
	SensorID    int     `json:"sensor_id"`
	AirportCode string  `json:"airport_code"`
	Measurement string  `json:"measurement"`
	Value       float64 `json:"value"`
	Timestamp   string  `json:"timestamp"`
}

type Config struct {
	MqttBroker       string     `json:"mqtt_broker"`
	MqttPort         int        `json:"mqtt_port"`
	ClientID         string     `json:"client_id"`
	Qos              byte       `json:"qos"`
	Interval         int        `json:"interval"`
	SensorID         int        `json:"sensor_id"`
	AirportCode      string     `json:"airport_code"`
	TemperatureRange ValueRange `json:"temperature_range"`
	WindSpeedRange   ValueRange `json:"wind_speed_range"`
	PressureRange    ValueRange `json:"pressure_range"`
}

type ValueRange struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

// Fonction pour charger la configuration
func loadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	err = json.NewDecoder(file).Decode(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// Fonction pour simuler les données du capteur
func simulateSensorData(sensorID int, airportCode, measurement string, valueRange ValueRange) SensorData {
	return SensorData{
		SensorID:    sensorID,
		AirportCode: airportCode,
		Measurement: measurement,
		Value:       valueRange.Min + rand.Float64()*(valueRange.Max-valueRange.Min),
		Timestamp:   strconv.FormatInt(time.Now().Unix(), 10),
	}
}

func main() {
	config, err := loadConfig("config.json")
	if err != nil {
		fmt.Println("Erreur de chargement de la configuration:", err)
		os.Exit(1)
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("%s:%d", config.MqttBroker, config.MqttPort))
	opts.SetClientID(config.ClientID)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalln("Erreur de connexion :", token.Error())
	}

	sensorID := rand.Int() // ID de capteur aléatoire

	for {

		// Simulation des données pour chaque mesure
		tempData := simulateSensorData(sensorID, config.AirportCode, "Temperature", config.TemperatureRange)
		windData := simulateSensorData(sensorID, config.AirportCode, "Wind speed", config.WindSpeedRange)
		pressureData := simulateSensorData(sensorID, config.AirportCode, "Atmospheric pressure", config.PressureRange)

		// Publier les données pour chaque mesure
		publishSensorData(client, tempData, config.Qos)
		publishSensorData(client, windData, config.Qos)
		publishSensorData(client, pressureData, config.Qos)

		// Attendre avant la prochaine publication
		time.Sleep(time.Duration(config.Interval) * time.Second)
	}
}

// publishSensorData publie les données du capteur sur le topic MQTT correspondant.
func publishSensorData(client mqtt.Client, data SensorData, qos byte) {
	// Sérialisation des données du capteur en JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Erreur de sérialisation JSON :", err)
		return
	}

	topic := fmt.Sprintf("/aeroport/%s/%s", data.AirportCode, data.Measurement)
	token := client.Publish(topic, qos, false, jsonData)
	token.Wait()

	fmt.Printf("Données publiées : %s\n", string(jsonData))
}
