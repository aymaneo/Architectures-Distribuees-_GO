package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type SensorData struct {
	Valeur    float64 `json:"valeur"`
	Timestamp int     `json:"time"`
	Mesure    string
	SensorID  string
	IATA      string
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
	config, err := loadConfig("./cmd/alertManager/config.json")
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
		var TopicArray []string = strings.Split(msg.Topic(), "/")
		data.IATA = TopicArray[2]
		data.Mesure = TopicArray[3]
		data.SensorID = TopicArray[4]
		checkThresholdAndAlert(client, data, config)
	})

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Erreur de connexion MQTT: %v", token.Error())
	}

	topics := map[string]byte{
		"/Airport/+/temperature/#": 0,
		"/Airport/+/vent/#":        0,
		"/Airport/+/pression/#":    0,
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

func checkThresholdAndAlert(client mqtt.Client, data SensorData, config *Config) bool {
	threshold, ok := config.Thresholds[data.Mesure]
	if !ok {
		fmt.Println("pas de seuil")
		return false // Pas de seuil défini pour ce type de mesure
	}
	if data.Valeur > threshold {
		alertMsg := fmt.Sprintf("Alerte ! %s pour %s à l'aéroport %s a dépassé le seuil. Valeur: %.2f - Timestamp UNIX: %d",
			data.Mesure, data.SensorID, data.IATA, data.Valeur, data.Timestamp)
		fmt.Println(alertMsg) // Afficher l'alerte dans la console

		// Publier l'alerte sur le topic d'alerte
		token := client.Publish(config.AlertTopic, 0, false, alertMsg)
		token.Wait()
		return true
	} else {
		fmt.Println("ne depasse pas ...")
		return false
	}
}
