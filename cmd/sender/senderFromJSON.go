package main

import (
	"TestProject/pkg/captorClass"
	"fmt"
	"time"
)

func main() {
	var timeInSecs int = 15

	var listeConf, _ = captorClass.LoadConfig("./config.json")
	var liste = captorClass.ListOfCaptors(listeConf)

	for 1 < 1000 {
		fmt.Println("Sending Data")
		time.Sleep(time.Second * time.Duration(timeInSecs))
		for _, s := range liste.Cpt {
			s.Pub()
			s.NextValue()
			time.Sleep(time.Millisecond * 100)
		}
	}
}
