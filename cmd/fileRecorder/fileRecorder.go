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
	Value       float64 `json:"valeur"`
	Timestamp   string  `json:"time"`
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

func getOpenFile(filePath string, header string) (*os.File, error) {
	var of *os.File
	var err error
	// if the file does not exist
	if _, err = os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		of, err = os.Create(filePath)
		if err != nil {
			return nil, errors.New("an error occurred when trying to create the file " +
				filePath)
		}

		// Writing CSV header in the .csv file
		_, err = of.Write([]byte(header))
		if err != nil {
			return nil, errors.New("an error occurred when trying to write to the file " +
				filePath)
		}
	} else {
		of, err = os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, errors.New("an error occurred when trying to open the file " +
				filePath)
		}
	}
	return of, err
}

func main() {
	// First second of the last recorded date (in epoch format)
	var csvFilePath string = ""

	// Loading configuration
	config, err := loadConfig(filepath.FromSlash("config.json"))
	if err != nil {
		log.Fatalf("configuration loading error: %v", err)
	}

	// Creating MQTT client
	opts := mqtt.NewClientOptions().
		AddBroker(fmt.Sprintf("%s:%d", config.MqttBroker, config.MqttPort)).
		SetClientID(config.ClientID)

	// Setting default MQTT message handler
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		sensorData, err := unmarshalMQTTMessage(&msg)
		if err != nil {
			log.Fatal("MQTT message unmarshalling error")
		}

		// Identifying where to write the data (what day and what town is it from ?)
		msgTimestamp, err := strconv.Atoi(sensorData.Timestamp)
		if err != nil {
			log.Fatalf("error when trying to parse sensorData timestamp : %v", err)
		}
		var msgDate time.Time = time.Unix(int64(msgTimestamp), 0)
		csvFilePath = fmt.Sprintf("../%s-airport-data-%d-%d-%d.csv",
			sensorData.AirportCode, msgDate.Day(), msgDate.Month(), msgDate.Year())
		csvFilePath = filepath.FromSlash(csvFilePath)

		exportData(csvFilePath, ",", &sensorData)
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
		log.Fatalf("subscription error : %v", token.Error())
	}

	// Waiting for the stop signals for a clean termination.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	client.Disconnect(1000)
	log.Printf("Deconnected from MQTT broker.")
}

func unmarshalMQTTMessage(msg *mqtt.Message) (SensorData, error) {
	var sensorData SensorData
	var msgJsonData struct {
		Value     float64 `json:"valeur"`
		Timestamp int     `json:"time"`
	}

	err := json.Unmarshal((*msg).Payload(), &msgJsonData)
	if err != nil {
		log.Printf("MQTT message unmarshalling error : %v\n", err)
		return SensorData{}, errors.New("an error occurred " +
			"when trying to unmarshal an MQTT message")
	}

	// Building the SensorData structure with the received data
	msgTopicElems := strings.Split((*msg).Topic(), "/")
	sensorData.SensorID, err = strconv.Atoi(msgTopicElems[4])
	if err != nil {
		log.Printf("SensorID parsing error : %v", err)
		return SensorData{}, errors.New("an error occurred " +
			"when trying to parse an incoming SensorID")
	}
	sensorData.AirportCode = msgTopicElems[2]
	sensorData.Timestamp = strconv.Itoa(msgJsonData.Timestamp)
	sensorData.Value = msgJsonData.Value
	sensorData.Measurement = msgTopicElems[3]

	return sensorData, nil
}

func exportData(destFilePath string, separator string, data *SensorData) {
	of, err := getOpenFile(destFilePath,
		strings.Join(
			[]string{"Timestamp", "AirportCode", "SensorID", "Measurement", "Value"},
			separator)+"\n")
	if err != nil {
		log.Fatalf("an error occurred when trying to get open file %s : %v",
			destFilePath, err)
	}
	// Extract data to be written from the SensorData struct
	row := strings.Join(
		[]string{
			data.Timestamp, data.AirportCode, strconv.Itoa(data.SensorID),
			data.Measurement,
			strconv.FormatFloat(data.Value, 'f', -1, 64),
		}, separator) + "\n"
	_, err = of.Write([]byte(row))
	if err != nil {
		log.Fatalf("an error occurred when trying to write to the file %s : %v",
			destFilePath, err)
	}
	err = of.Close()
	if err != nil {
		log.Fatalf("an error occurred when trying to close the file %s : %v",
			destFilePath, err)
	}
	log.Printf("Message data %s recorded to %s", row, destFilePath)
}
