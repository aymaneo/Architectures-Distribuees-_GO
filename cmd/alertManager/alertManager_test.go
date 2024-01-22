package main

import (
	"TestProject/pkg/captorClass"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	file, _ := loadConfig("./config.json")
	assert.NotNil(t, file)
}

func TestCheckThresholdAndAlert(t *testing.T) {
	blankCaptor := captorClass.InitCaptorType("test", "test", 60.0, 10.0, 0.1, 22)
	blankClient := captorClass.InitCaptor("Test", blankCaptor, "testForTresh", "mqtt://158.178.194.137:1883")
	cnf := &Config{
		MqttBroker: "",
		MqttPort:   0,
		ClientID:   "",
		Thresholds: map[string]float64{
			"temperature": 20.0,
			"vent":        45.0,
			"pression":    1100.0,
		},
		AlertTopic: "/Alertes",
	}
	data1 := SensorData{
		Valeur:    45.0,
		Timestamp: 1705957491,
		Mesure:    "temperature",
		SensorID:  "TEST1",
		IATA:      "TLS",
	}

	fmt.Println(checkThresholdAndAlert(blankClient.Mqtt, data1, cnf))
	assert.True(t, checkThresholdAndAlert(blankClient.Mqtt, data1, cnf)) //alerte émise car la température dépasse le treshold
	data2 := SensorData{
		Valeur:    5.0,
		Timestamp: 1705957491,
		Mesure:    "temperature",
		SensorID:  "TEST1",
		IATA:      "TLS",
	}
	assert.False(t, checkThresholdAndAlert(blankClient.Mqtt, data2, cnf)) //alerte non émise car la température ne dépasse pas le treshold

	data3 := SensorData{
		Valeur:    5.0,
		Timestamp: 1705957491,
		Mesure:    "WeirdMesure",
		SensorID:  "TEST1",
		IATA:      "TLS",
	}
	assert.False(t, checkThresholdAndAlert(blankClient.Mqtt, data3, cnf))
}
