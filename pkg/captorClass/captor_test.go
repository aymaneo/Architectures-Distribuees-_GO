package captorClass

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitCaptorType(t *testing.T) {
	captorType := InitCaptorType("Temperature", "°C", 50.0, -10.0, 0.15, 22.0)
	assert.Equal(t, "Temperature", captorType.Name, "Pas le bon nom ! ")
}

func TestInitCaptor(t *testing.T) {
	captor := InitCaptor("TLS", InitCaptorType("Temperature", "°C", 50.0, -10.0, 0.15, 22.0), "123456789", "tcp://158.178.194.137:1883")
	assert.Equal(t, "TLS", captor.Aita, "Pas le bon code AITA ! ")
	assert.NotNil(t, captor)
}

func TestCaptor_Pub(t *testing.T) {
	captor := InitCaptor("TLS", InitCaptorType("Temperature", "°C", 50.0, -10.0, 0.15, 22.0), "123456789", "tcp://158.178.194.137:1883")
	assert.Equal(t, 0, captor.Pub(), "Erreur dans l'envoi de la data")
}

func TestCaptor_NextValue(t *testing.T) {
	captor := InitCaptor("TLS", InitCaptorType("Temperature", "°C", 50.0, -10.0, 0.15, 22.0), "123456789", "tcp://158.178.194.137:1883")
	captor.Value = 55.0
	captor.NextValue()
	//assert.NotEqual(t, 22.0, captor.Value, "Valeur inchangée même après le NextValue()")
	captor.Value = -15.0
	captor.NextValue()
	assert.NotEqual(t, 22.0, captor.Value, "Valeur inchangée même après le NextValue()")
}

func TestCaptor_Print(t *testing.T) {
	captor := InitCaptor("TLS", InitCaptorType("Temperature", "°C", 50.0, -10.0, 0.15, 22.0), "123456789", "tcp://158.178.194.137:1883")
	strvalue := captor.Print()
	assert.Equal(t, 0, strvalue, "Error on printing !")
}

func TestInitCaptor2(t *testing.T) {
	captor := InitCaptor("TLS", InitCaptorType("Temperature", "°C", 50.0, -10.0, 0.15, 22.0), "123456789", "tcp://158.178.194.137:183")
	assert.Nil(t, captor.Mqtt)
	//assert.
}

func TestArgsToCaptor(t *testing.T) {
	//oubli dans un premier cas de l'unité : °C
	var inputValue = []string{"senderCLI.go", "temperature", "50.0", "-10.0", "0.15", "22.0", "LIS", "11111222", "tcp://158.178.194.137:1883", "10"}
	captor, waitTime := ArgsToCaptor(inputValue)
	assert.Nil(t, captor)
	assert.NotEqual(t, int64(10), waitTime, "Fail au test 1")

	var inputValue2 = []string{"senderCLI.go", "temperature", "°C", "50.0", "-10.0", "0.15", "22.0", "LIS", "11111222", "tcp://158.178.194.137:1883", "10"}
	captor2, waitTime2 := ArgsToCaptor(inputValue2)
	assert.NotNil(t, captor2)
	assert.Equal(t, int64(10), waitTime2, "Fail au test 2")
}

func TestLoadConfig(t *testing.T) {
	failedCnf, err := LoadConfig("./inexistantconfig.json")

	assert.NotNil(t, err)
	assert.Empty(t, failedCnf) //Vide car pas de conf

	successedCnf, err2 := LoadConfig("./config_test.json")
	assert.Nil(t, err2)
	assert.NotEmpty(t, successedCnf)
	assert.Equal(t, 2, len(successedCnf.Captors))
}

func TestListOfCaptors(t *testing.T) {
	data, _ := LoadConfig("./config_test.json")
	list := ListOfCaptors(data)
	assert.NotEmpty(t, list)

}
