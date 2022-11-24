package tcp

import (
	"encoding/hex"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/edgexfoundry/device-modbus/parser"
	modbusclient "github.com/goburrow/modbus"
)

func OpenModbusCommunication(
	endpoint string,
	port int,
	slaveId int,
	timeout time.Duration,
	verbose bool,
) *modbusclient.TCPClientHandler {

	// Turn on the debug trace option, to see what is being transmitted
	handler := modbusclient.NewTCPClientHandler("endpoint:" + strconv.Itoa(port))
	handler.SlaveId = byte(slaveId)
	handler.Timeout = timeout
	if verbose {
		handler.Logger = log.New(os.Stdout, "mb: ", log.Lshortfile) // Necessary for debugging
	}

	return handler
}

func HandleSensorData(
	handler *modbusclient.TCPClientHandler,
	sensor *parser.Sensor,
	verbose bool,
) (map[int]interface{}, error) {

	clientResponse := make(map[int]interface{})

	// Connect as Modbus RTU client
	cerr := handler.Connect()

	if cerr != nil {
		log.Printf("TCP Connection Err: %s", cerr)
	}

	// Create MODBUS client
	client := modbusclient.NewClient(handler)
	var err error

	if verbose {
		log.Printf(
			`
			Handler configuration:
				Address: %s,
				SlaveID: %d,
				Timeout: %v
			`,
			handler.Address,
			handler.SlaveId,
			handler.Timeout,
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
