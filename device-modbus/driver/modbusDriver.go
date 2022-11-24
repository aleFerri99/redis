package driver

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/edgexfoundry/device-modbus/modbusHandler/ascii"
	"github.com/edgexfoundry/device-modbus/modbusHandler/rtu"
	"github.com/edgexfoundry/device-modbus/modbusHandler/tcp"

	"github.com/edgexfoundry/device-modbus/parser"
)

func HandleSensors(sensorList *parser.Sensors, serialResource *string, verbose bool, readGap *time.Duration) []map[string]interface{} {

	if verbose {
		prettyPrint, err := json.MarshalIndent(sensorList, "", "   ")
		if err != nil {
			log.Printf("Cannot pretty print data struct to JSON. Error: %s", err)
			os.Exit(0)
		}
		log.Print(string(prettyPrint))
	}

	return openAndHandle(sensorList, serialResource, verbose, readGap)
}

func openAndHandle(sensorList *parser.Sensors, serialResource *string, verbose bool, readGap *time.Duration) []map[string]interface{} {

	readings := []map[string]interface{}{}
	for sensor, sensorProperties := range sensorList.Sensors {

		// Print some sensor oriented informations
		log.Printf("\n\nHandling sensor %s...\n\n", sensor)

		switch sensorProperties.Mode {
		case "rtu":
			reading := make(map[string]interface{})
			properties := &parser.SerialProperties{}
			byteData, err := json.Marshal(sensorProperties.ModbusProperty)
			if err != nil {
				log.Printf("Cannot marshal incoming data. Error: %s", err)
				continue
			}
			err = json.Unmarshal(byteData, properties)
			if err != nil {
				log.Printf("Cannot unmarshal incoming data. Error: %s", err)
				continue
			}
			handler := rtu.OpenModbusCommunication(
				*serialResource,
				int(properties.Address),
				properties.Baudrate,
				int(properties.DataBits),
				int(properties.StopBits),
				properties.Parity,
				time.Duration(properties.Timeout)*time.Second,
				verbose,
			)
			// Start sensor handling
			sensorReadings, err := rtu.HandleSensorData(handler, &sensorProperties, verbose)
			if err != nil {
				log.Printf("Error handling sensor data. Error: %s", err)
				continue
			}

			reading["SensorName"] = sensor
			reading["DataHex"] = sensorReadings
			readings = append(readings, reading)
			if verbose {
				prettyPrint, err := json.MarshalIndent(readings, "", "\t")
				if err != nil {
					log.Printf("Cannot pretty print sensors readings. Error: %s", err)
				}
				log.Printf("PRETTY PRINT: %s", string(prettyPrint))
			}
			// Push value to EdgeX core data
		case "ascii":
			reading := make(map[string]interface{})
			properties := &parser.AsciiProperties{}
			byteData, err := json.Marshal(sensorProperties.ModbusProperty)
			if err != nil {
				log.Printf("Cannot marshal incoming data. Error: %s", err)
				continue
			}
			err = json.Unmarshal(byteData, properties)
			if err != nil {
				log.Printf("Cannot unmarshal incoming data. Error: %s", err)
				continue
			}
			handler := ascii.OpenModbusCommunication(
				*serialResource,
				int(properties.Address),
				properties.Baudrate,
				int(properties.DataBits),
				int(properties.StopBits),
				properties.Parity,
				time.Duration(properties.Timeout)*time.Second,
				verbose,
			)

			// Start sensor handling
			sensorReadings, err := ascii.HandleSensorData(handler, &sensorProperties, verbose)
			if err != nil {
				log.Printf("Error handling sensor data. Error: %s", err)
				continue
			}

			reading["SensorName"] = sensor
			reading["DataHex"] = sensorReadings
			readings = append(readings, reading)
			if verbose {
				prettyPrint, err := json.MarshalIndent(readings, "", "\t")
				if err != nil {
					log.Printf("Cannot pretty print sensors readings. Error: %s", err)
				}
				log.Printf("PRETTY PRINT: %s", string(prettyPrint))
			}
			// Push value to EdgeX core data
		case "tcp":
			reading := make(map[string]interface{})
			properties := &parser.TcpProperties{}
			byteData, err := json.Marshal(sensorProperties.ModbusProperty)
			if err != nil {
				log.Printf("Cannot marshal incoming data. Error: %s", err)
				continue
			}
			err = json.Unmarshal(byteData, properties)
			if err != nil {
				log.Printf("Cannot unmarshal incoming data. Error: %s", err)
				continue
			}
			handler := tcp.OpenModbusCommunication(
				*serialResource,
				properties.Port,
				int(properties.Address),
				time.Duration(properties.Timeout)*time.Second,
				verbose,
			)

			// Start sensor handling
			sensorReadings, err := tcp.HandleSensorData(handler, &sensorProperties, verbose)
			if err != nil {
				log.Printf("Error handling sensor data. Error: %s", err)
				continue
			}

			reading["SensorName"] = sensor
			reading["DataHex"] = sensorReadings
			readings = append(readings, reading)
			if verbose {
				prettyPrint, err := json.MarshalIndent(readings, "", "\t")
				if err != nil {
					log.Printf("Cannot pretty print sensors readings. Error: %s", err)
				}
				log.Printf("PRETTY PRINT: %s", string(prettyPrint))
			}
			// Push value to EdgeX core data
		default:
			log.Printf("Unknown modbus mode: %s", sensorProperties.Mode)
			return nil
		}

		// Add the possibilty to pass a custom delay between each sensor as an ENV var
		time.Sleep(*readGap)
	}

	if verbose {
		log.Printf("Readings: %v", readings)
	}

	return readings
}
