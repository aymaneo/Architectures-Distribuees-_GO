package main

import (
	"TestProject/pkg/receiverClass"
)

func main() {
	var dbToken string = "xqBlL61ya6BtT6P_F3TrLLfbTOwXnUF9ZTX2dcZAOBcHOtA7QCNuFo8HphUsYJgTLdXZozjs22P-EtVxMe6N1g=="
	var dbURL string = "http://158.178.194.137:8086"
	var org = "IMT"
	var bucket = "projetGo"
	var broker = "158.178.194.137"
	var port = 1883
	var IdClient = "go_mqtt_client"

	var influxClient = receiverClass.CreateInfluxClient(dbToken, dbURL, org, bucket)
	influxClient.CreateMosquittoClient(broker, port, IdClient)

	influxClient.Sub("/Airport/#")
	for 1 > 0 {

	}
	influxClient.MqttClient.Disconnect(250)
}
