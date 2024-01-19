package captorClass

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"math/rand"
	"time"
)

func createClientOptions(brokerURI string, clientId string) *mqtt.ClientOptions {
	opt := mqtt.NewClientOptions()
	opt.AddBroker(brokerURI)
	opt.SetClientID(clientId)
	return opt
}

func connect(brokerURI string, clientId string) mqtt.Client {
	opt := createClientOptions(brokerURI, clientId)
	client := mqtt.NewClient(opt)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		//log.Fatal(err)
		//fmt.Println("error !")
		return nil
	} else {
		return client
	}
}

//Name : nom du capteur ex: "Temperature"
//Unit : unité de la valeur ex: "°C"
//UpperRange : Valeur à ne pas dépasser pour le noise
//LowerRange : Valeur à ne pas dépasser pour le noise
//IncrementRange : plage de variation pour la fonction noise : Value = Value + [ -IncrementRange, +IncrementRange]
//DefaultValue : Valeur par défaut du capteur ex: 1020 mPa pour la pression

type CaptorType struct {
	Name           string
	Unit           string
	UpperRange     float64
	LowerRange     float64
	IncrementRange float64
	DefaultValue   float64
}

type ConnectionParam struct {
	BrokerURI string
	Topic     string
	CapteurID string
}

type Captor struct {
	Aita      string
	CapType   CaptorType
	ConParams ConnectionParam
	Mqtt      mqtt.Client
	Value     float64
}

type Data struct {
	Valeur float64 `json:"valeur"`
	Time   int64   `json:"time"`
}

func (cap *Captor) NextValue() {
	var newValue = cap.Value + (rand.Float64() * cap.CapType.IncrementRange * 2.0) - cap.CapType.IncrementRange
	for (newValue > cap.CapType.UpperRange) || (newValue < cap.CapType.LowerRange) {
		//newValue = newValue + (rand.Float64() * cap.CapType.IncrementRange * 2.0) - cap.CapType.IncrementRange
		if newValue > cap.CapType.UpperRange {
			newValue = newValue - cap.CapType.IncrementRange*rand.Float64()
		}
		if newValue < cap.CapType.LowerRange {
			newValue = newValue + cap.CapType.IncrementRange*rand.Float64()
		}
	}
	cap.Value = newValue
}

func (cap Captor) Print() int {
	fmt.Println("------------Capteur -------------")
	fmt.Println("Type de Capteur   :", cap.CapType.Name)
	fmt.Println("Valeur du capteur :", cap.Value, cap.CapType.Unit)
	fmt.Println("---------------------------------")
	return 0
}

func (cap Captor) Pub() int {
	//text := fmt.Sprintf("%f", cap.Value)
	dataToSend := Data{Valeur: cap.Value, Time: time.Now().Unix()}
	fmt.Println(dataToSend)
	jsn, err := json.Marshal(dataToSend)
	if err != nil {
		fmt.Println("Erreur de sérialisation JSON :", err)
		return 10
	}
	cap.Mqtt.Publish(cap.ConParams.Topic, 0, false, jsn)
	return 0
}

func InitCaptor(CodeAita string, CapType *CaptorType, idCaptor string, brokerURI string) *Captor {
	ob := new(Captor)
	ob.Aita = CodeAita
	ob.CapType = *CapType
	ob.Value = ob.CapType.DefaultValue

	ConOb := new(ConnectionParam)
	ConOb.BrokerURI = brokerURI
	ConOb.CapteurID = idCaptor
	ConOb.Topic = "/Airport/" + ob.Aita + "/" + ob.CapType.Name + "/" + idCaptor

	ob.ConParams = *ConOb

	ob.Mqtt = connect(brokerURI, idCaptor)
	return ob
}

func InitCaptorType(name string, unit string, uprange float64, lowrange float64, incr float64, defaultvalue float64) *CaptorType {
	ob := new(CaptorType)
	ob.Name = name
	ob.Unit = unit
	ob.UpperRange = uprange
	ob.LowerRange = lowrange
	ob.IncrementRange = incr
	ob.DefaultValue = defaultvalue
	return ob
}
