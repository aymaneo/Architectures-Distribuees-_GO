package captorClass

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"io/ioutil"
	"math/rand"
	"strconv"
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

func ArgsToCaptor(liste []string) (*Captor, int64) {
	if len(liste) != 11 {
		fmt.Println("Attention il manque des paramètres ! Voici comment utiliser la commande : ")
		//                                      1            2      3         4          5          6       7       8         9           10
		fmt.Println("go run senderCLI.go [typeMesure] [unité] [ValMax] [ValMin] [incrément] [ValDeft] [IATA] [idMQTT] [MQTTURI] [tmpAttSec]")
		//fmt.Println(liste)
		return nil, 0
	} else {
		uprange, _ := strconv.ParseFloat(liste[3], 64)
		lowrange, _ := strconv.ParseFloat(liste[4], 64)
		incr, _ := strconv.ParseFloat(liste[5], 64)
		defaultvalue, _ := strconv.ParseFloat(liste[6], 64)
		var Type = InitCaptorType(liste[1], liste[2], uprange, lowrange, incr, defaultvalue)

		var capteur = InitCaptor(liste[7], Type, liste[8], liste[9])
		var waitTime, _ = strconv.ParseInt(liste[10], 10, 64)
		return capteur, waitTime
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

//partie pour avoir des capteurs à partir d'un fichier JSON

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
	Cpt []*Captor
}

func LoadConfig(path string) (CaptorConfig, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return CaptorConfig{}, err
	}
	var captorConfig CaptorConfig
	err2 := json.Unmarshal(file, &captorConfig)
	if err2 != nil {
		return CaptorConfig{}, err2
	}
	return captorConfig, nil
}

func ListOfCaptors(config CaptorConfig) ListeCapteurs {
	var listeCapteurs ListeCapteurs
	for i, s := range config.Captors {
		fmt.Println(i, s.Airport)
		listeCapteurs.Cpt = append(listeCapteurs.Cpt, InitCaptor(s.Airport, InitCaptorType(s.Name, s.Unit, s.Uprange, s.Lowrange, s.Incr, s.DefaultValue), s.MQTTId, s.BrokerURI))
	}
	return listeCapteurs
}
