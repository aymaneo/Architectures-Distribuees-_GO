package receiverClass

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"strings"
	"time"
)

func CreateInfluxClient(dbToken string, dbURL string, bucket string, org string) *Receiver {
	client := influxdb2.NewClient(dbURL, dbToken)
	rec := new(Receiver)
	rec.InfluxClient = client
	rec.Org = org
	rec.Bucket = bucket
	return rec
}

type Data struct {
	Valeur float64 `json:"valeur"`
	Time   int64   `json:"time"`
}

type Intermediate struct {
	InfluxClient influxdb2.Client
}

func (rec Receiver) messagePubHandler(client mqtt.Client, msg mqtt.Message) {
	var clientInflux influxdb2.Client = rec.InfluxClient
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	var value = msg.Payload()
	var TopicArray []string = strings.Split(msg.Topic(), "/")
	//fmt.Printf(TopicArray[0] + "\n")
	fmt.Printf(TopicArray[1] + "\n") //Airport : valeur fixe
	if TopicArray[1] == "Airport" {  // Si on a le bon topic (pas de messages d'info)
		fmt.Printf(TopicArray[2] + "\n") //code Aita
		fmt.Printf(TopicArray[3] + "\n") //type capteur
		fmt.Printf(TopicArray[4] + "\n") //Id
		//fmt.Printf(value + "\n")
		//fmt.Printf(time.Now().String() + "\n")
		//fmt.Println(value)
		var data Data
		err := json.Unmarshal(value, &data)
		if err != nil {
			fmt.Println(err.Error())
		}

		floatValue := data.Valeur
		t := data.Time
		fmt.Println(floatValue)
		fmt.Println(t)
		writeAPI := clientInflux.WriteAPI("IMT", "projetGo")
		p := influxdb2.NewPointWithMeasurement(TopicArray[3]).
			AddTag("TypeData", TopicArray[1]).
			AddTag("ville", TopicArray[2]).
			AddTag("idCapteur", TopicArray[4]).
			AddField("valeur", floatValue).
			SetTime(time.Unix(t, 0))
		writeAPI.WritePoint(p)
		// Flush writes
		writeAPI.Flush()
	}
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func (rec *Receiver) CreateMosquittoClient(broker string, port int, clientId string) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID(clientId)
	fmt.Println(broker)
	fmt.Println(port)
	fmt.Println(clientId)
	//opts.SetUsername("emqx")
	//opts.SetPassword("public")
	opts.SetDefaultPublishHandler(rec.messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	rec.MqttClient = client
}

func (rec Receiver) Sub(topic string) {
	rec.MqttClient.Subscribe(topic, 1, rec.messagePubHandler)
	fmt.Printf("Subscribed to topic: %s", topic)
}

type Receiver struct {
	MqttClient   mqtt.Client
	InfluxClient influxdb2.Client
	Org          string
	Bucket       string
}

func InitReceiver(clientMqtt mqtt.Client, clientInflux influxdb2.Client) *Receiver {
	ob := new(Receiver)
	ob.MqttClient = clientMqtt
	ob.InfluxClient = clientInflux
	return ob
}
