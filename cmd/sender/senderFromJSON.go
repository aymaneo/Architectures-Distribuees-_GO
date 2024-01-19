package main

import (
	"TestProject/pkg/captorClass"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type captorData struct {
	Name         string  `json:"name"`
	Unit         string  `json:"unit"`
	Uprange      float64 `json:"uprange"`
	Lowrange     float64 `json:"lowrange"`
	Incr         float64 `json:"incr"`
	DefaultValue float64 `json:"defaultValue"`
	Airport      string  `json:"Airport"`
	BrokerURI    string  `json:"BrokerURI"`
	MQTTId       string  `json:"MQTTId"`
}

type CaptorConfig struct {
	Captors []captorData `json:"captors"`
}

type ListeCapteurs struct {
	Cpt []*captorClass.Captor
}

func loadConfig(path string) (CaptorConfig, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return CaptorConfig{}, err
	}
	var captorConfig CaptorConfig
	err = json.Unmarshal(file, &captorConfig)
	if err != nil {
		return CaptorConfig{}, err
	}
	return captorConfig, nil
}

func main() {
	var timeInSecs int = 15
	var captorConfig, err = loadConfig("./config.json")
	if err != nil {
		fmt.Println(err)
	}
	var listeCapteurs ListeCapteurs
	for i, s := range captorConfig.Captors {
		fmt.Println(i, s.Airport)
		listeCapteurs.Cpt = append(listeCapteurs.Cpt, captorClass.InitCaptor(s.Airport, captorClass.InitCaptorType(s.Name, s.Unit, s.Uprange, s.Lowrange, s.Incr, s.DefaultValue), s.MQTTId, s.BrokerURI))
	}

	for 1 < 1000 {
		fmt.Println("Sending Data")
		time.Sleep(time.Second * time.Duration(timeInSecs))
		for _, s := range listeCapteurs.Cpt {
			s.Pub()
			s.NextValue()
			time.Sleep(time.Millisecond * 100)
		}
	}
}
