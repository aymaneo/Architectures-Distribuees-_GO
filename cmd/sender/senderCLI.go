package main

import (
	"TestProject/pkg/captorClass"
	"os"
	"time"
)

func main() {
	capteur, tempsAttente := captorClass.ArgsToCaptor(os.Args)
	if capteur != nil && tempsAttente != 0 {
		for 1 > 0 {
			capteur.Pub()
			capteur.NextValue()
			time.Sleep(time.Second * time.Duration(tempsAttente))
		}
	}
}
