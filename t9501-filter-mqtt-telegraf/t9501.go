package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

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

type T9501 struct {
	SensorName    string
	DataHex       []map[int]interface{}
	DataHexSensor map[string]interface{}
	Data          T9501Data
}

type T9501Data struct {
	FirmwareRevision       float32 // RO
	RelativeHumidity       float32 // RO
	Temperature            float32 // RO
	ModbusSlaveAddr        int16   // RW
	BaudRate               uint32  // RW
	ConfigFlag             uint16
	SerialProtocol         string // FIXED RTU
	ParityBit              string // RW
	DataBits               string // FROM PARITY
	Stopbits               string // FROM PARITY
	SerialResponseDelay    int16  // N/D
	SerialActiveOutputTime int16  // N/D
	Timestamp              int64
}

func NewT9501(debug *bool) *T9501 {
	*verbose = *debug
	return &T9501{}
}

func (t9501 *T9501) ParseSensorData(ctx interfaces.AppFunctionContext, data interface{}) (continuePipeline bool, result interface{}) {

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

	value := []map[string]interface{}{}
	err = json.Unmarshal([]byte(edgexMsg.Readings[0].Value), &value)
	if err != nil {
		log.Printf("Cannot unmarshal string to []map[string]interface{}{}. Error: %s", err)
		return false, nil
	}

	// Filter data by device of interest
	hit := false
	for _, entry := range value {
		if entry["SensorName"] == os.Getenv("T9501_SENSOR_NAME") {
			hit = true
			t9501.DataHexSensor = entry["DataHex"].(map[string]interface{})
		}
	}

	if hit {
		log.Printf("Handling %s data...", os.Getenv("T9501_SENSOR_NAME"))
	} else {
		log.Printf("Device data not available.")
		return false, nil
	}

	if len(t9501.DataHexSensor) == 0 {
		log.Println("Empty entry from message bus. Waiting for next message...")
		return false, nil
	}

	// Perform data processing to convert HEX data to understandable data
	err = t9501.ParseDataFromHex()
	if err != nil {
		log.Printf("Cannot convert data from HEX. Error: %s", err)
		return false, nil
	}

	t9501Message := make(map[string]interface{})
	t9501.Data.Timestamp = time.Now().UnixMilli()
	t9501Message["t9501"] = t9501.Data

	return true, t9501Message
}
