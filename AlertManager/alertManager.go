package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"

	mqtt "github.com/eclipse/paho.mqtt.golang"
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
	file, err := ioutil.ReadFile(path)
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
	config, err := loadConfig("./AlertManager/config.json")
	if err != nil {
		log.Fatalf("Erreur de chargement de la configuration: %v", err)
	}

	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("%s:%d", config.MqttBroker, config.MqttPort)).SetClientID(config.ClientID)
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		var data SensorData
		err := json.Unmarshal(msg.Payload(), &data)
		if err != nil {
			log.Printf("Erreur de décodage JSON: %v\n", err)
			return
		}
		checkThresholdAndAlert(client, data, config)
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

	// Attente des signaux d'arrêt pour une terminaison propre.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	client.Disconnect(1000)
	fmt.Println("Déconnecté du broker MQTT.")
}

func checkThresholdAndAlert(client mqtt.Client, data SensorData, config *Config) {
	threshold, ok := config.Thresholds[data.Measurement]
	if !ok {
		return // Pas de seuil défini pour ce type de mesure
	}

	if data.Value > threshold {
		alertMsg := fmt.Sprintf("Alerte ! %s pour %d à l'aéroport %s a dépassé le seuil. Valeur: %.2f - Timestamp: %s",
			data.Measurement, data.SensorID, data.AirportCode, data.Value, data.Timestamp)
		fmt.Println(alertMsg) // Afficher l'alerte dans la console

		// Publier l'alerte sur le topic d'alerte
		token := client.Publish(config.AlertTopic, 0, false, alertMsg)
		token.Wait()
	}
}
