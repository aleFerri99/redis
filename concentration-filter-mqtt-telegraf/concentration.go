package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/interfaces"
)

type edgexMessage struct {
	ApiVersion  string          `json:"apiVersion"`
	Id          string          `json:"id"`
	DeviceName  string          `json:"deviceName"`
	ProfileName string          `json:"profileName"`
	SourceName  string          `json:"sourceName"`
	Origin      int64           `json:"origin"`
	Readings    []edgexReadings `json:"readings"`
}

type edgexReadings struct {
	Id           string `json:"id"`
	Origin       int64  `json:"origin"`
	DeviceName   string `json:"deviceName"`
	ResourceName string `json:"resourceName"`
	ProfileName  string `json:"profileName"`
	ValueType    string `json:"valueType"`
	BinaryValue  []byte `json:"binaryValue"`
	MediaType    string `json:"mediaType"`
	Value        string `json:"value"`
}

type Concentration struct {
	SensorName string
	DataSensor map[string]interface{}
	Data       ConcentrationData
}

type ConcentrationData struct {
	Temperature float64
	Frequence   []interface{}
	Impedance   []interface{}
	Phase       []interface{}
	Timestamp   int64
}

func NewConcentration(debug *bool) *Concentration {
	*verbose = *debug
	return &Concentration{}
}

func (concentration *Concentration) ParseSensorData(ctx interfaces.AppFunctionContext, data interface{}) (continuePipeline bool, result interface{}) {

	/*
		EdgeX return a message struct of type:
		"apiVersion":"v2",
		"id":"7877cc63-31f9-43de-9084-676cfc3550cf",
		"deviceName":"Modbus-Device01",
		"profileName":"Modbus-Device",
		"sourceName":"Modbus",
		"origin":1660137654310302925,
		"readings":[{
			"id":"7be456df-5705-4a3c-b9ab-8aac16a237b0",
			"origin":1660137654310302925,
			"deviceName":"Modbus-Device01",
			"resourceName":"Modbus",
			"profileName":"Modbus-Device",
			"valueType":"String",
			"binaryValue":null,
			"mediaType": "",
			"value":"{}" <-- Part of interest
			}]
		}
	*/

	if *verbose {
		log.Printf("EdgeX Core message:\n%v", data)
	}

	edgexMsg := &edgexMessage{}

	err := json.Unmarshal([]byte(data.(string)), edgexMsg)
	fmt.Println(edgexMsg)
	if err != nil {
		log.Printf("Cannot unmarshal data to edgexMessage struct. Error: %s", err)
		return false, nil
	}

	if *verbose {
		prettyPrint, err := json.MarshalIndent(edgexMsg, "", "\t")
		if err != nil {
			log.Printf("Cannot pretty print EdgeX message. Error: %s", err)
		} else {
			log.Printf("Parsed EdgeX Core message: %v", string(prettyPrint))
		}
	}

	// value := []map[string]interface{}{}
	fmt.Printf("%v \n", edgexMsg.Readings[0].Value)
	err = json.Unmarshal([]byte(edgexMsg.Readings[0].Value), &concentration.DataSensor)
	fmt.Println(concentration.DataSensor)
	if err != nil {
		log.Printf("Cannot unmarshal string to map[string]interface{}{}. Error: %s", err)
		return true, nil
	}

	// // Filter data by device of interest
	// hit := false
	// for _, entry := range value {
	// 	if entry["SensorName"] == "ConcentrationentrationSensorDevice" {
	// 		hit = true
	// 		concentration.DataHexSensor = entry["DataHex"].(map[string]interface{})
	// 	}
	// }

	// if hit {
	// 	log.Printf("Handling %s data...", "ConcentrationentrationSensorDevice")
	// } else {
	// 	log.Printf("Device data not available.")
	// 	return false, nil
	// }

	if len(concentration.DataSensor) == 0 {
		log.Println("Empty entry from message bus. Waiting for next message...")
		return false, nil
	}

	// Perform data processing to convert HEX data to understandable data
	err = concentration.ParseDataFromHex()
	if err != nil {
		log.Printf("Cannot convert data from HEX. Error: %s", err)
		return false, nil
	}

	concentrationMessage := make(map[string]interface{})
	concentrationMessage["concentration"] = concentration.Data
	fmt.Println(concentrationMessage)

	return true, concentrationMessage
}
