package receiverClass

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

var dbToken string = "xqBlL61ya6BtT6P_F3TrLLfbTOwXnUF9ZTX2dcZAOBcHOtA7QCNuFo8HphUsYJgTLdXZozjs22P-EtVxMe6N1g=="
var dbURL string = "http://158.178.194.137:8086"
var org = "IMT"
var bucket = "projetGo"
var broker = "158.178.194.137"
var port = 1883
var IdClient = "go_mqtt_client_for_tests"

func TestCreateInfluxClient(t *testing.T) {
	var influxClient = CreateInfluxClient(dbToken, dbURL, bucket, org)
	assert.Equal(t, org, influxClient.Org, "Client Influx non connecté")
}

func TestReceiver_CreateMosquittoClient(t *testing.T) {
	var influxClient = CreateInfluxClient(dbToken, dbURL, bucket, org)
	influxClient.CreateMosquittoClient(broker, port, IdClient)
	assert.True(t, influxClient.MqttClient.IsConnected(), "Client MQTT non connecté")
}

func TestReceiver_Sub(t *testing.T) {
	var influxClient = CreateInfluxClient(dbToken, dbURL, bucket, org)
	influxClient.CreateMosquittoClient(broker, port, IdClient)
	assert.Equal(t, 0, influxClient.Sub("/#"))
}

func TestInitReceiver(t *testing.T) {
	var influxClient = CreateInfluxClient(dbToken, dbURL, bucket, org)
	influxClient.CreateMosquittoClient(broker, port, IdClient)
	var totalReceiver = InitReceiver(influxClient.MqttClient, influxClient.InfluxClient)
	assert.NotNil(t, totalReceiver)
}

type MockMQTTMessage struct {
	mock.Mock
}

func (m *MockMQTTMessage) Duplicate() bool {
	//TODO implement me
	//panic("implement me")
	return false
}

func (m *MockMQTTMessage) Qos() byte {
	//TODO implement me
	//panic("implement me")
	return 1
}

func (m *MockMQTTMessage) Retained() bool {
	//TODO implement me
	//panic("implement me")
	return false
}

func (m *MockMQTTMessage) MessageID() uint16 {
	//TODO implement me
	//panic("implement me")
	return 0
}

func (m *MockMQTTMessage) Ack() {
	//TODO implement me
	//panic("implement me")
}

func (m *MockMQTTMessage) Payload() []byte {
	args := m.Called()
	return args.Get(0).([]byte)
}

func (m *MockMQTTMessage) Topic() string {
	args := m.Called()
	return args.Get(0).(string)
}

type MockReceiver struct {
	mock.Mock
}

func (m *MockReceiver) messagePubHandler(client mqtt.Client, msg mqtt.Message) {
	m.Called(client, msg)
}

func TestMessageHandler(t *testing.T) {
	var influxClient = CreateInfluxClient(dbToken, dbURL, bucket, org)
	influxClient.CreateMosquittoClient(broker, port, IdClient)
	var totalReceiver = InitReceiver(influxClient.MqttClient, influxClient.InfluxClient)

	mockMsg := new(MockMQTTMessage)
	mockReceiver := new(MockReceiver)
	mockMsg.On("Payload").Return([]byte(`{"valeur": 42.0, "time": 1611600000}`))
	mockMsg.On("Topic").Return("/Airport/TLS/Capteur/123456")
	mockReceiver.On("messagePubHandler", mock.Anything, mock.Anything)
	totalReceiver.messagePubHandler(nil, mockMsg)
	mockReceiver.messagePubHandler(nil, mockMsg)
	mockMsg.AssertExpectations(t)
	mockReceiver.AssertExpectations(t)
}
