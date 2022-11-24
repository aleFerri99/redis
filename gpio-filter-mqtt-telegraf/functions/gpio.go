package functions

import (
	"encoding/json"
	"log"
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

type gpioAndConfig struct {
	gpio       GPIO
	gpioConfig Config
}

type GPIO struct {
	Name      string
	Chip      string
	Line      int
	State     bool
	Timestamp int64
}

type Config struct {
	PumpTimer     time.Duration
	EnableClean   bool
	CleanTimer    time.Duration
	EnableReverse bool
	ReverseTimer  time.Duration
	GravityTimer  time.Duration
	CommandGap    time.Duration
}

var verbose bool

func NewGpio(debug *bool) *GPIO {
	verbose = *debug
	return &GPIO{}
}

func (gpio *GPIO) ParseSensorData(ctx interfaces.AppFunctionContext, data interface{}) (continuePipeline bool, result interface{}) {

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

	if verbose {
		log.Printf("EdgeX Core message:\n%v", data)
	}

	edgexMsg := &edgexMessage{}

	err := json.Unmarshal([]byte(data.(string)), edgexMsg)
	if err != nil {
		log.Printf("Cannot unmarshal data to edgexMessage struct. Error: %s", err)
		return false, nil
	}

	if verbose {
		prettyPrint, err := json.MarshalIndent(edgexMsg, "", "\t")
		if err != nil {
			log.Printf("Cannot pretty print EdgeX message. Error: %s", err)
		} else {
			log.Printf("Parsed EdgeX Core message: %v", string(prettyPrint))
		}
	}

	gpioMsg := &gpioAndConfig{}

	err = json.Unmarshal([]byte(edgexMsg.Readings[0].Value), &gpioMsg)
	if err != nil {
		log.Printf("Cannot unmarshal string to []map[string]interface{}{}. Error: %s", err)
		return false, nil
	}

	// Perform data processing to convert HEX data to understandable data

	gpioMessage := make(map[string]interface{})
	gpioMessage["gpio"] = gpioMsg.gpio
	gpioMessage["config"] = gpioMsg.gpioConfig

	return true, gpioMessage
}
