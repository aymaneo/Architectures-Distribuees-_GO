package main

import (
	"encoding/json"
	"errors"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
)

type SensorData struct {
	SensorID    int     `json:"sensor_id"`
	AirportCode string  `json:"airport_code"`
	Measurement string  `json:"measurement"`
	Value       float64 `json:"value"`
	Timestamp   string  `json:"timestamp"`
}

type Config struct {
	MqttBroker string             `json:"mqtt_broker"`
	MqttPort   int                `json:"mqtt_port"`
	ClientID   string             `json:"client_id"`
	Thresholds map[string]float64 `json:"thresholds"`
	AlertTopic string             `json:"alert_topic"`
}

func loadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func main() {
	var allSensorData []SensorData
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("Erreur de chargement de la configuration: %v", err)
	}

	opts := mqtt.NewClientOptions().
		AddBroker(fmt.Sprintf("%s:%d", config.MqttBroker, config.MqttPort)).
		SetClientID(config.ClientID)
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		var data SensorData
		err := json.Unmarshal(msg.Payload(), &data)
		if err != nil {
			log.Printf("Erreur de décodage JSON: %v\n", err)
			return
		}
		allSensorData = append(allSensorData, data)
		fmt.Println("allSensorData : ", allSensorData)
	})

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Erreur de connexion MQTT: %v", token.Error())
	}
	fmt.Println("Connecté au broker MQTT, en attente de données...")

	topics := map[string]byte{
		"aeroport/+/Temperature":          0,
		"aeroport/+/Wind speed":           0,
		"aeroport/+/Atmospheric pressure": 0,
	}
	if token := client.SubscribeMultiple(topics, nil); token.Wait() && token.Error() != nil {
		log.Fatalf("Erreur de souscription: %v", token.Error())
	}

	//TODO : traduire tout le code en anglais (commentaires, messages de log, etc... inclus ?)

	// Attente des signaux d'arrêt pour une terminaison propre.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	client.Disconnect(1000)

	freeFileName := searchFreeFileName("test", ".csv")
	exportData(freeFileName, ",", allSensorData)

	//TODO : mettre en place les logs
	log.Printf("Déconnecté du broker MQTT.") //TODO : vérifier si cette instruction utilise le bon niveau de log
}

func searchFreeFileName(nameBase string, fileExtension string) string {
	fullFileName := nameBase + fileExtension
	for i := 0; ; i++ {
		// If file already exists
		if _, err := os.Stat(fullFileName); !errors.Is(err, os.ErrNotExist) {
			// Look for another file name that is free
			fullFileName = nameBase + strconv.Itoa(i) + fileExtension
		} else {
			return fullFileName
		}
	}
}

func exportData(destFilePath string, separator string, data []SensorData) {
	row := strings.Join([]string{"SensorID", "AirportCode", "Timestamp",
		"Value", "Measurement"}, separator)
	err := os.WriteFile(destFilePath, []byte(row), 0666)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range data {
		// Extract data from the SensorData struct for this row
		row = strings.Join([]string{strconv.Itoa(v.SensorID), v.AirportCode,
			v.Timestamp,
			strconv.FormatFloat(v.Value, 4, -1, 64),
			v.Measurement}, separator)
		err := os.WriteFile(destFilePath, []byte(row), 0666)
		if err != nil {
			log.Fatal(err)
		}
	}
}
