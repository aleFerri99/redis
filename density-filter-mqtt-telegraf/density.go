package main

import (
	"encoding/json"
	"log"
	"os"

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

type Density struct {
	SensorName    string
	DataHex       []map[int]interface{}
	DataHexSensor map[string]interface{}
	Data          DensityData
}

type DensityData struct {
	DynamicViscosity       float32 // HEX to single-precision float number
	Density                float32 // HEX to single-precision float number
	Temperature            float32 // HEX to single-precision float number
	ModbusSlaveAddr        int16   // Type: RW, Register Address (HEX/DEC): 0x0200/512, Function Number: 3/6/16, Range: 0-255, Default Value: 1
	BaudRate               int32   // Type: RW, Register Address (HEX/DEC): 0x0201/513, Function Number: 3/6/16, Range: 0: 1200bps, 1: 2400bps, 2: 4800bps, 3: 9600bps, 4: 19200bps, 5: 38400bps, Default Value: 3: 9600bps
	SerialProtocol         string  // Type: RW, Register Address (HEX/DEC): 0x0202/514, Function Number: 3/6/16, Range: 0: Modbus RTU, 1: Modbus ASCII, Default Value: 0: Modbus RTU
	ParityBit              string  // Type: RW, Register Address (HEX/DEC): 0x0203/515, Function Number: 3/6/16, Range: 0: No parity, 1: Even parity, 2: Odd parity, Default Value: 0: No parity
	DataBits               string  // Type: RW, Register Address (HEX/DEC): 0x0204/516, Function Number: 3/6/16, Range: 1: 8 Data bits, Default Value: 1: 8 data bits
	Stopbits               string  // Type: RW, Register Address (HEX/DEC): 0x0205/517, Function Number: 3/6/16, Range: 0: 1 Stop bit, 1: 2 Stop bits, Default Value: 0: 1 Stop bit
	SerialResponseDelay    int16   // Type: RW, Register Address (HEX/DEC): 0x0206/518, Function Number: 3/6/16, Range: 0-250 (=0-2500ms), Default Value: 0 (=DISABLED)
	SerialActiveOutputTime int16   // Type: RW, Register Address (HEX/DEC): 0x0207/519, Function Number: 3/6/16, Range: 0-250 (=0-250s), Default Value: 0 (=DISABLED)
	Timestamp              int64
}

func NewDensity(debug *bool) *Density {
	*verbose = *debug
	return &Density{}
}

func (Density *Density) ParseSensorData(ctx interfaces.AppFunctionContext, data interface{}) (continuePipeline bool, result interface{}) {

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
		if entry["SensorName"] == os.Getenv("DENSITY_SENSOR_NAME") {
			hit = true
			Density.DataHexSensor = entry["DataHex"].(map[string]interface{})
		}
	}

	if hit {
		log.Printf("Handling %s data...", os.Getenv("DENSITY_SENSOR_NAME"))
	} else {
		log.Printf("Device data not available.")
		return false, nil
	}

	if len(Density.DataHexSensor) == 0 {
		log.Println("Empty entry from message bus. Waiting for next message...")
		return false, nil
	}

	// Perform data processing to convert HEX data to understandable data
	err = Density.ParseDataFromHex()
	if err != nil {
		log.Printf("Cannot convert data from HEX. Error: %s", err)
		return false, nil
	}

	DensityMessage := make(map[string]interface{})
	DensityMessage["Density"] = Density.Data

	return true, DensityMessage
}
