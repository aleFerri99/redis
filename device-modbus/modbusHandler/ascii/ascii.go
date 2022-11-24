package ascii

import (
	"encoding/hex"
	"log"
	"os"
	"time"

	"github.com/edgexfoundry/device-modbus/parser"
	modbusclient "github.com/goburrow/modbus"
)

func OpenModbusCommunication(
	serialPort string,
	slaveId int,
	baudRate int,
	dataBits int,
	stopBits int,
	parity string,
	timeout time.Duration,
	verbose bool,
) *modbusclient.ASCIIClientHandler {

	// Turn on the debug trace option, to see what is being transmitted
	handler := modbusclient.NewASCIIClientHandler(serialPort)
	handler.SlaveId = byte(slaveId)
	handler.BaudRate = baudRate
	handler.DataBits = dataBits
	handler.StopBits = stopBits
	handler.Timeout = timeout
	handler.Parity = parity
	if verbose {
		handler.Logger = log.New(os.Stdout, "mb: ", log.Lshortfile) // Necessary for debugging
	}

	return handler
}

func HandleSensorData(
	handler *modbusclient.ASCIIClientHandler,
	sensor *parser.Sensor,
	verbose bool,
) (map[int]interface{}, error) {

	clientResponse := make(map[int]interface{})

	// Connect as Modbus RTU client
	cerr := handler.Connect()

	if cerr != nil {
		log.Printf("ASCII Connection Err: %s", cerr)
	}

	// Create MODBUS client
	client := modbusclient.NewClient(handler)
	var err error

	if verbose {
		log.Printf(
			`
			Handler configuration:
				SlaveID: %d,
				BaudRate: %d,
				DataBits: %d,
				StopBits: %d,
				Timeout: %v,
				Parity: %s
			`,
			handler.SlaveId,
			handler.BaudRate,
			handler.DataBits,
			handler.StopBits,
			handler.Timeout,
			handler.Parity,
		)
	}

	for index, register := range sensor.Modbus.Registers {
		switch fc := sensor.Modbus.ReadWrite[index].FunctionCode; {
		case fc == 0x01:
			readResult, readErr := client.ReadCoils(uint16(register), uint16(sensor.Modbus.ReadLen[index]))
			if readErr != nil {
				log.Println(readErr)
				err = readErr
				break
			}
			log.Printf("Rx: %x", readResult)
			clientResponse[register] = hex.EncodeToString(readResult)
		case fc == 0x02:
			readResult, readErr := client.ReadDiscreteInputs(uint16(register), uint16(sensor.Modbus.ReadLen[index]))
			if readErr != nil {
				log.Println(readErr)
				err = readErr
				break
			}
			log.Printf("Rx: %x", readResult)
			clientResponse[register] = hex.EncodeToString(readResult)
		case fc == 0x03:
			readResult, readErr := client.ReadHoldingRegisters(uint16(register), uint16(sensor.Modbus.ReadLen[index]))
			if readErr != nil {
				log.Println(readErr)
				err = readErr
				break
			}
			log.Printf("Rx: %x", readResult)
			clientResponse[register] = hex.EncodeToString(readResult)
		case fc == 0x04:
			readResult, readErr := client.ReadInputRegisters(uint16(register), uint16(sensor.Modbus.ReadLen[index]))
			if readErr != nil {
				log.Println(readErr)
				err = readErr
				break
			}
			log.Printf("Rx: %x", readResult)
			clientResponse[register] = hex.EncodeToString(readResult)
		case fc == 0x05:
			writeResult, writeErr := client.WriteSingleCoil(uint16(register), uint16(sensor.Modbus.ReadWrite[index].Value[0]))
			if writeErr != nil {
				log.Println(writeErr)
				err = writeErr
				break
			}
			log.Printf("Rx: %x", writeResult)
			clientResponse[register] = hex.EncodeToString(writeResult)
		case fc == 0x06:
			writeResult, writeErr := client.WriteSingleRegister(uint16(register), uint16(sensor.Modbus.ReadWrite[index].Value[0]))
			if writeErr != nil {
				log.Println(writeErr)
				err = writeErr
				break
			}
			log.Printf("Rx: %x", writeResult)
			clientResponse[register] = hex.EncodeToString(writeResult)
		case fc == 0x0f:
			values := []byte{}
			for value := range sensor.Modbus.ReadWrite[index].Value {
				values = append(values, byte(value))
			}
			writeResult, writeErr := client.WriteMultipleCoils(
				uint16(register),
				uint16(len(sensor.Modbus.ReadWrite[index].Value)),
				values,
			)
			if writeErr != nil {
				log.Println(writeErr)
				err = writeErr
				break
			}
			log.Printf("Rx: %x", writeResult)
			clientResponse[register] = hex.EncodeToString(writeResult)
		case fc == 0x10:
			values := []byte{}
			for value := range sensor.Modbus.ReadWrite[index].Value {
				values = append(values, byte(value))
			}
			writeResult, writeErr := client.WriteMultipleRegisters(
				uint16(register),
				uint16(len(sensor.Modbus.ReadWrite[index].Value)),
				values,
			)
			if writeErr != nil {
				log.Println(writeErr)
				err = writeErr
				break
			}
			log.Printf("Rx: %x", writeResult)
			clientResponse[register] = hex.EncodeToString(writeResult)
		default:
			log.Printf("Hit unknown function code %d.", fc)
		}
	}

	// Release modbus resource
	handler.Close()

	return clientResponse, err
}
