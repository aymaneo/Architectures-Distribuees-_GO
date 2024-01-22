package main

import (
	"TestProject/pkg/receiverClass"
)

func main() {
	//Ici on passe par un objet receiver pour pouvoir dupliquer facilement le processus
	//comme pour Captor, afin de pouvoir ajouter de la redondance
	//A terme le but sera de passer par un fichier init.json qui paramettra le nombre de
	//capteur et de receiver
	var dbToken string = "xqBlL61ya6BtT6P_F3TrLLfbTOwXnUF9ZTX2dcZAOBcHOtA7QCNuFo8HphUsYJgTLdXZozjs22P-EtVxMe6N1g=="
	var dbURL string = "http://158.178.194.137:8086"
	var org = "IMT"
	var bucket = "projetGo"
	var broker = "158.178.194.137"
	var port = 1883
	var IdClient = "go_mqtt_client"

	var influxClient = receiverClass.CreateInfluxClient(dbToken, dbURL, org, bucket)
	influxClient.CreateMosquittoClient(broker, port, IdClient)

	influxClient.Sub("/#")
	for 1 > 0 {

	}
	influxClient.MqttClient.Disconnect(250)
}
