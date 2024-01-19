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
