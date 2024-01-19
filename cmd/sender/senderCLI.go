package main

import (
	"TestProject/pkg/captorClass"
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {

	if len(os.Args) != 11 {
		fmt.Println("Attention il manque des paramètres ! Voici comment utiliser la commande : ")
		//                                      1            2      3         4          5          6       7       8         9           10
		fmt.Println("go run senderCLI.go [typeMesure] [unité] [ValMax] [ValMin] [incrément] [ValDeft] [IATA] [idMQTT] [MQTTURI] [tmpAttSec]")
	} else {
		uprange, _ := strconv.ParseFloat(os.Args[3], 64)
		lowrange, _ := strconv.ParseFloat(os.Args[4], 64)
		incr, _ := strconv.ParseFloat(os.Args[5], 64)
		defaultvalue, _ := strconv.ParseFloat(os.Args[6], 64)
		var Type = captorClass.InitCaptorType(os.Args[1], os.Args[2], uprange, lowrange, incr, defaultvalue)

		var capteur = captorClass.InitCaptor(os.Args[7], Type, os.Args[8], os.Args[9])
		var waitTime, _ = strconv.ParseInt(os.Args[10], 10, 64)
		for 1 > 0 {
			capteur.Pub()
			capteur.NextValue()
			time.Sleep(time.Second * time.Duration(waitTime))
		}
	}
	//	var Vent = captorClass.InitCaptorType("vent", "km/h", 120.0, 0.0, 0.5, 12.0)
	//  var capteurQuatre = captorClass.InitCaptor(Nantes, Vent, "412345678", "tcp://158.178.194.137:1883")
}
