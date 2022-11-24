package main

import (
	"fmt"
	"log"
	"strconv"
)

func (conductivity *Conductivity) ParseDataFromHex() error {

	for register, hexValue := range conductivity.DataHexSensor {
		reg, err := strconv.Atoi(register)
		if err != nil {
			log.Printf("Cannot parse string to integer. Error: %s", err)
			return err
		}
		switch reg {
		case 0x00:
			temperatureHigh, err := strconv.ParseInt(hexValue.(string)[:2], 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}
			temperatureLow, err := strconv.ParseInt(hexValue.(string)[2:], 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}
			var data float32
			if hexValue.(string)[:2] != "ff" { // Positive temperatures
				data = (float32(temperatureHigh)*256 + float32(temperatureLow)) / 100
			} else { //Negative temperatures
				data = ((float32(temperatureHigh)*256 + float32(temperatureLow)) - 65535 - 1) / 100
			}

			conductivity.Data.Temperature = data
		case 0x02:
			conductivityHigh, err := strconv.ParseInt(hexValue.(string)[:2], 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}
			conductivityLow, err := strconv.ParseInt(hexValue.(string)[2:], 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}

			conductivity.Data.ElectricalConductivity = (float32(conductivityHigh)*256 + float32(conductivityLow))

		case 0x03:
			salinityHigh, err := strconv.ParseInt(hexValue.(string)[:2], 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}
			salinityLow, err := strconv.ParseInt(hexValue.(string)[2:], 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}

			conductivity.Data.Salinity = (float32(salinityHigh)*256 + float32(salinityLow))

		case 0x04:
			tdsHigh, err := strconv.ParseInt(hexValue.(string)[:2], 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}
			tdsLow, err := strconv.ParseInt(hexValue.(string)[2:], 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}

			conductivity.Data.Tds = (float32(tdsHigh)*256 + float32(tdsLow))

		case 0x08:
			ecRawAdHigh, err := strconv.ParseInt(hexValue.(string)[:2], 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}
			ecRawAdLow, err := strconv.ParseInt(hexValue.(string)[2:], 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}

			conductivity.Data.EcRawAd = (float32(ecRawAdHigh)*256 + float32(ecRawAdLow))

		case 0x20:
			tempCmp, err := strconv.ParseInt(hexValue.(string), 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}

			if tempCmp == 0 {
				conductivity.Data.TempCompensationStatus = "External temperature Sensor"
			} else if tempCmp == 1 {
				conductivity.Data.TempCompensationStatus = "Onboard temperature Sensor"
			} else {
				conductivity.Data.TempCompensationStatus = "Disabled"
			}

		case 0x22:
			ecTempCoff, err := strconv.ParseInt(hexValue.(string), 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}

			conductivity.Data.EcTempCoff = int16(ecTempCoff)

		case 0x23:
			salinityCoff, err := strconv.ParseInt(hexValue.(string), 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}

			conductivity.Data.EcTempCoff = int16(salinityCoff)

		case 0x24:
			tdsCoff, err := strconv.ParseInt(hexValue.(string), 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}

			conductivity.Data.EcTempCoff = int16(tdsCoff)

		case 0x25:
			electrodeConstant, err := strconv.ParseInt(hexValue.(string), 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}

			conductivity.Data.EcTempCoff = int16(electrodeConstant)

		case 0x30:
			cecCalib1413, err := strconv.ParseInt(hexValue.(string), 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}

			conductivity.Data.EcTempCoff = int16(cecCalib1413)

		case 0x31:
			ecCalib12880, err := strconv.ParseInt(hexValue.(string), 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}

			conductivity.Data.EcTempCoff = int16(ecCalib12880)

		case 0x0200:
			modbusAddr, err := strconv.ParseInt(hexValue.(string), 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}
			conductivity.Data.ModbusSlaveAddr = int16(modbusAddr)

		case 0x0201:
			baudRate, err := strconv.ParseInt(hexValue.(string), 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}
			switch baudRate {
			case 0:
				conductivity.Data.BaudRate = 1200

			case 1:
				conductivity.Data.BaudRate = 2400

			case 2:
				conductivity.Data.BaudRate = 4800

			case 3:
				conductivity.Data.BaudRate = 9600

			case 4:
				conductivity.Data.BaudRate = 19200

			case 5:
				conductivity.Data.BaudRate = 38400
			}

		case 0x0202:
			serialProtocol, err := strconv.ParseInt(hexValue.(string), 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}
			if serialProtocol == 0 {
				conductivity.Data.SerialProtocol = "Modbus RTU"
			} else {
				conductivity.Data.SerialProtocol = "Modbus ASCII"
			}

		case 0x0203:
			parity, err := strconv.ParseInt(hexValue.(string), 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}
			if parity == 1 {
				conductivity.Data.ParityBit = "Even parity"
			} else if parity == 2 {
				conductivity.Data.ParityBit = "Odd parity"
			} else {
				conductivity.Data.ParityBit = "No parity"
			}

		case 0x0204:
			data, err := strconv.ParseInt(hexValue.(string), 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}
			if data == 1 {
				conductivity.Data.DataBits = "8 data bits"
			} else {
				conductivity.Data.DataBits = "Error, this device only supports 8 data bits format!"
			}

		case 0x0205:
			stop, err := strconv.ParseInt(hexValue.(string), 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}
			if stop == 0 {
				conductivity.Data.Stopbits = "1 stop bit"
			} else {
				conductivity.Data.Stopbits = "2 stop bits"
			}

		case 0x0206:
			responseDelay, err := strconv.ParseInt(hexValue.(string), 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}
			conductivity.Data.SerialResponseDelay = int16(responseDelay * 10)

		case 0x0207:
			outputInterval, err := strconv.ParseInt(hexValue.(string), 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}
			conductivity.Data.SerialResponseDelay = int16(outputInterval)

		// Questo default rompe tutto, indagare sul perchè
		default:
			return fmt.Errorf("default case entered for register: %d", reg)
		}
	}

	return nil
}
