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

type Conductivity struct {
	SensorName    string
	DataHex       []map[int]interface{}
	DataHexSensor map[string]interface{}
	Data          ConductivityData
}

type ConductivityData struct {
	Temperature            float32 // Type: RO, Register Address (HEX/DEC): 0x0000/0, Function Number: 3/4, Range: -4000-8000 (=-40.00-80.00 °C), Default Value: N/A
	ElectricalConductivity float32 // Type: RO, Register Address (HEX/DEC): 0x0002/2, Function Number: 3/4, Range: 0-2000 (=0-2000us/cm), Default Value: N/A
	Salinity               float32 // Type: RO, Register Address (HEX/DEC): 0x0003/3, Function Number: 3/4, Range: 0-2000 (=0-2000mg/L), Default Value: N/A
	Tds                    float32 // Type: RO, Register Address (HEX/DEC): 0x0004/4, Function Number: 3/4, Range: 0-2000 (=0-2000mg/L), Default Value: N/A
	EcRawAd                float32 // Type: RO, Register Address (HEX/DEC): 0x0008/8, Function Number: 3/4, Range: 0-4000, Default Value: N/A
	TempCompensationStatus string  // Type: RW, Register Address (HEX/DEC): 0x0020/32, Function Number: 3/6/16, Range: 0: External sensor, 1: Onboard sensor, 2: Disabled, Default Value: 0
	EcTempCoff             int16   // Type: RW, Register Address (HEX/DEC): 0x0022/34, Function Number: 3/6/16, Range: 0-100 (=0.0%-10.00%), Default Value: 20 (=20%)
	SalinityCoff           int16   // Type: RW, Register Address (HEX/DEC): 0x0023/35, Function Number: 3/6/16, Range: 0-100 (=0.00-1.00), Default Value: 55 (=0.55)
	TdsCoff                int16   // Type: RW, Register Address (HEX/DEC): 0x0024/36, Function Number: 3/6/16, Range: 0-100 (=0.00-1.00), Default Value: 5 (=0.5)
	ElectrodeConst         float32 // Type: RW, Register Address (HEX/DEC): 0x0025/37, Function Number: 3/6/16, Range: 500-1500 (=0.500-1.500), Default Value: 1000 (=1.000)
	EcCalib_1413           float32 // Type: RW, Register Address (HEX/DEC): 0x0030/48, Function Number: 3/6/16, Range: Immerse electrode in 1413 us/cm solution for a while and write 0xFFFF into register to perform auto calibration, Default Value: 223
	EcCalib_12880          float32 // Type: RW, Register Address (HEX/DEC): 0x0031/49, Function Number: 3/6/16, Range: Immerse electrode in 12880 us/cm solution for a while and write 0xFFFF into register to perform auto calibration, Default Value: 1851
	ModbusSlaveAddr        int16   // Type: RW, Register Address (HEX/DEC): 0x0200/512, Function Number: 3/6/16, Range: 0-255, Default Value: 1 or 30
	BaudRate               int32   // Type: RW, Register Address (HEX/DEC): 0x0201/513, Function Number: 3/6/16, Range: 0: 1200bps, 1: 2400bps, 2: 4800bps, 3: 9600bps, 4: 19200bps, 5: 38400bps, Default Value: 3: 9600bps
	SerialProtocol         string  // Type: RW, Register Address (HEX/DEC): 0x0202/514, Function Number: 3/6/16, Range: 0: Modbus RTU, 1: Modbus ASCII, Default Value: 0: Modbus RTU
	ParityBit              string  // Type: RW, Register Address (HEX/DEC): 0x0203/515, Function Number: 3/6/16, Range: 0: No parity, 1: Even parity, 2: Odd parity, Default Value: 0: No parity
	DataBits               string  // Type: RW, Register Address (HEX/DEC): 0x0204/516, Function Number: 3/6/16, Range: 1: 8 Data bits, Default Value: 1: 8 data bits
	Stopbits               string  // Type: RW, Register Address (HEX/DEC): 0x0205/517, Function Number: 3/6/16, Range: 0: 1 Stop bit, 1: 2 Stop bits, Default Value: 0: 1 Stop bit
	SerialResponseDelay    int16   // Type: RW, Register Address (HEX/DEC): 0x0206/518, Function Number: 3/6/16, Range: 0-250 (=0-2500ms), Default Value: 0 (=DISABLED)
	SerialActiveOutputTime int16   // Type: RW, Register Address (HEX/DEC): 0x0207/519, Function Number: 3/6/16, Range: 0-250 (=0-250s), Default Value: 0 (=DISABLED)
	Timestamp              int64
}

func NewConductivity(debug *bool) *Conductivity {
	*verbose = *debug
	return &Conductivity{}
}

func (Conductivity *Conductivity) ParseSensorData(ctx interfaces.AppFunctionContext, data interface{}) (continuePipeline bool, result interface{}) {

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
		if entry["SensorName"] == os.Getenv("CONDUCTIVITY_SENSOR_NAME") {
			hit = true
			Conductivity.DataHexSensor = entry["DataHex"].(map[string]interface{})
		}
	}

	if hit {
		log.Printf("Handling %s data...", os.Getenv("CONDUCTIVITY_SENSOR_NAME"))
	} else {
		log.Printf("Device data not available.")
		return false, nil
	}

	if len(Conductivity.DataHexSensor) == 0 {
		log.Println("Empty entry from message bus. Waiting for next message...")
		return false, nil
	}

	// Perform data processing to convert HEX data to understandable data
	err = Conductivity.ParseDataFromHex()
	if err != nil {
		log.Printf("Cannot convert data from HEX. Error: %s", err)
		return false, nil
	}

	ConductivityMessage := make(map[string]interface{})
	Conductivity.Data.Timestamp = time.Now().Unix()
	ConductivityMessage["Conductivity"] = Conductivity.Data

	return true, ConductivityMessage
}
