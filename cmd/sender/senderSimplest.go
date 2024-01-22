package main

import (
	"TestProject/pkg/captorClass"
	"fmt"
	"time"
)

func main() {
	//var Toulouse = captorClass.InitLoc("France", "Toulouse")
	//var Paris = captorClass.InitLoc("France", "Paris")
	var Toulouse = "TLS"
	var Nantes = "NTE"
	var Paris = "CDG"

	var Temperature = captorClass.InitCaptorType("temperature", "°C", 50.0, -10.0, 0.15, 22.0)
	var Luminosite = captorClass.InitCaptorType("luminosite", "Lux", 150.0, 10.0, 0.15, 70)
	var Pression = captorClass.InitCaptorType("pression", "hPa", 1200.0, 800, 0.015, 1010)
	var Vent = captorClass.InitCaptorType("vent", "km/h", 120.0, 0.0, 0.5, 12.0)

	//Pour les id des capteurs : Ils sont utilisés lors de la connection au broker MQTT, il faut qu'ils soient uniques
	//entre eux, sinon ils se fermeront la connection tout à tour

	var capteurUn = captorClass.InitCaptor(Toulouse, Temperature, "13234421", "tcp://158.178.194.137:1883")
	var capteurDeux = captorClass.InitCaptor(Paris, Luminosite, "22124Z1", "tcp://158.178.194.137:1883")
	var capteurTrois = captorClass.InitCaptor(Toulouse, Pression, "3456789", "tcp://158.178.194.137:1883")
	var capteurQuatre = captorClass.InitCaptor(Nantes, Vent, "412345678", "tcp://158.178.194.137:1883")

	var increment int64 = 0

	for increment < 1000 {
		capteurUn.NextValue()
		capteurUn.Pub()
		capteurDeux.NextValue()
		capteurDeux.Pub()
		capteurTrois.NextValue()
		capteurTrois.Pub()
		capteurQuatre.NextValue()
		capteurQuatre.Pub()
		capteurQuatre.Print()
		fmt.Printf("sending data\n")
		increment++
		time.Sleep(time.Second * 5)
	}

	capteurUn.Print()
	capteurDeux.Print()
}
