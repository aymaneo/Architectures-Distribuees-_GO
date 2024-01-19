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
	config, err := loadConfig("../sensors/config.json")
	if err != nil {
		log.Fatalf("Erreur de chargement de la configuration: %v", err)
	}
	opts := mqtt.NewClientOptions().
		AddBroker(fmt.Sprintf("%s:%d", config.MqttBroker, config.MqttPort)).
		SetClientID(config.ClientID)

	freeFileName := searchFreeFileName("../sensorData", ".csv")
	row := "SensorID,AirportCode,Timestamp,Value,Measurement\n"
	of, err := os.Create(freeFileName)
	if err != nil {
		log.Fatalf("Erreur lors de la tentative de création du fichier %s : %v", freeFileName, err)
	}
	_, err = of.Write([]byte(row))
	if err != nil {
		log.Fatalf("Erreur lors de la tentative d'écriture dans le fichier %s : %v", freeFileName, err)
	}
	err = of.Close()
	if err != nil {
		log.Fatalf("Erreur lors de la tentative de fermeture du fichier %s : %v", freeFileName, err)
	}

	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		var data SensorData
		err := json.Unmarshal(msg.Payload(), &data)
		if err != nil {
			log.Printf("Erreur de décodage JSON: %v\n", err)
			return
		}
		exportData(freeFileName, ",", &data)
	})

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Erreur de connexion MQTT: %v", token.Error())
	}
	log.Println("Connecté au broker MQTT, en attente de données...")

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

func exportData(destFilePath string, separator string, data *SensorData) {
	// Extract data from the SensorData struct
	row := strings.Join([]string{strconv.Itoa(data.SensorID), data.AirportCode,
		data.Timestamp, strconv.FormatFloat(data.Value, 'f', -1, 64),
		data.Measurement}, separator) + "\n"
	log.Println(row)
	of, err := os.OpenFile(destFilePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("openerr : ", err)
	}
	_, err = of.WriteString(row)
	if err != nil {
		log.Fatal("writeerr :", err)
	}
	err = of.Close()
	if err != nil {
		log.Fatal("closeerr :", err)
	}
}
