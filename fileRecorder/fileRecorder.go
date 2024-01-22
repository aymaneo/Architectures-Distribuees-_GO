package main

import (
	"encoding/json"
	"errors"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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
	// First second of the last recorded date (in epoch format)
	var lastRecordedDate int64 = -1
	csvFilePath := filepath.FromSlash("../airport-data-1-1-1970.csv")

	// Loading configuration
	config, err := loadConfig(filepath.FromSlash("../sensors/config.json"))
	if err != nil {
		log.Fatalf("Erreur de chargement de la configuration: %v", err)
	}

	// Creating MQTT client
	opts := mqtt.NewClientOptions().
		AddBroker(fmt.Sprintf("%s:%d", config.MqttBroker, config.MqttPort)).
		SetClientID(config.ClientID)

	// Setting default MQTT message handler
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		data, err := unmarshalMQTTMessage(&msg)
		if err != nil {
			log.Fatal("MQTT message unmarshalling error")
		}

		// If this is the first message of the day
		if msgTimeStamp, _ := strconv.Atoi(data.Timestamp); int64(msgTimeStamp)-lastRecordedDate >= 86400 {
			// Update the last recorded date and the path of the .csv file to write to
			lastRecordedDate = int64(msgTimeStamp - (msgTimeStamp % 86400))
			lastRecordedTime := time.Unix(lastRecordedDate, 0)
			csvFilePath = fmt.Sprintf("../airport-data-%d-%d-%d.csv", lastRecordedTime.Day(),
				lastRecordedTime.Month(), lastRecordedTime.Year())
			csvFilePath = filepath.FromSlash(csvFilePath)

			// Writing CSV header in the new .csv file
			row := "Timestamp,AirportCode,SensorID,Measurement,Value\n"
			of, err := os.Create(csvFilePath)
			if err != nil {
				log.Fatalf("An error occurred when trying to create the file %s : %v", csvFilePath, err)
			}
			_, err = of.Write([]byte(row))
			if err != nil {
				log.Fatalf("An error occurred when trying to write to the file %s : %v", csvFilePath, err)
			}
			err = of.Close()
			if err != nil {
				log.Fatalf("An error occurred when trying to close the file %s : %v", csvFilePath, err)
			}
		}

		exportData(csvFilePath, ",", &data)
	})

	// Connecting to the MQTT broker
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("MQTT connection error: %v", token.Error())
	}
	log.Println("Connected to the MQTT broker, waiting for data...")

	topics := map[string]byte{
		"/Airport/#": 0,
	}
	if token := client.SubscribeMultiple(topics, nil); token.Wait() && token.Error() != nil {
		log.Fatalf("Subscription error: %v", token.Error())
	}

	//TODO : traduire tout le code en anglais (commentaires, messages de log, etc... inclus ?)

	// Waiting for the stop signals for a clean termination.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	client.Disconnect(1000)
	//TODO : mettre en place les logs
	log.Printf("Déconnecté du broker MQTT.") //TODO : vérifier si cette instruction utilise le bon niveau de log
}

func unmarshalMQTTMessage(msg *mqtt.Message) (SensorData, error) {
	var data SensorData
	err := json.Unmarshal((*msg).Payload(), &data)
	if err != nil {
		log.Printf("Erreur de décodage : %v\n", err)
		return SensorData{}, errors.New("an error occurred when trying to unmarshal an MQTT message")
	}
	msgTopicElems := strings.Split((*msg).Topic(), "/")
	data.SensorID, _ = strconv.Atoi(msgTopicElems[3])
	data.AirportCode = msgTopicElems[1]
	data.Measurement = msgTopicElems[2]
	return data, nil
}

func exportData(destFilePath string, separator string, data *SensorData) {
	// Extract data from the SensorData struct
	row := strings.Join(
		[]string{
			data.Timestamp, data.AirportCode, strconv.Itoa(data.SensorID),
			data.Measurement,
			strconv.FormatFloat(data.Value, 'f', -1, 64),
		}, separator) + "\n"
	log.Println("Message data saved : " + row)
	of, err := os.OpenFile(destFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("An error occurred when trying to open the file "+destFilePath+" : ", err)
	}
	_, err = of.Write([]byte(row))
	if err != nil {
		log.Fatal("An error occurred when trying to write to the file "+destFilePath+" : ", err)
	}
	err = of.Close()
	if err != nil {
		log.Fatal("An error occurred when trying to close the file "+destFilePath+" : ", err)
	}
}
