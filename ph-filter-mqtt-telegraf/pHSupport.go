package main

import (
	"fmt"
	"log"
	"strconv"
)

func (ph *Ph) ParseDataFromHex() error {

	for register, hexValue := range ph.DataHexSensor {
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

			ph.Data.Temperature = data
		case 0x01:
			phHigh, err := strconv.ParseInt(hexValue.(string)[:2], 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}
			phLow, err := strconv.ParseInt(hexValue.(string)[2:], 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}
			data := (float32(phHigh)*256 + float32(phLow)) / 100
			ph.Data.PhValue = data

		case 0x02:
			phCalibRawAdHigh, err := strconv.ParseInt(hexValue.(string)[:2], 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}
			phCalibRawAdLow, err := strconv.ParseInt(hexValue.(string)[2:], 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}
			var data float32
			if hexValue.(string)[:2] != "ff" { // Positive ph calibration
				data = (float32(phCalibRawAdHigh)*256 + float32(phCalibRawAdLow)) / 100
			} else { //Negative ph calibration
				data = ((float32(phCalibRawAdHigh)*256 + float32(phCalibRawAdLow)) - 65535 - 1) / 100
			}
			ph.Data.PhCalibrationRawAd = data

		case 0x20:
			temperatureComp, err := strconv.ParseInt(hexValue.(string), 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}
			if temperatureComp == 1 {
				ph.Data.TempCompensationStatus = "Turn off temperature compensation"
			} else {
				ph.Data.TempCompensationStatus = "Turn on temperature compensation"
			}

		case 0x30:
			phCalibRawAd0High, err := strconv.ParseInt(hexValue.(string)[:2], 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}
			phCalibRawAd0Low, err := strconv.ParseInt(hexValue.(string)[2:], 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}
			var data float32
			if hexValue.(string)[:2] != "ff" { // Positive ph calibration
				data = (float32(phCalibRawAd0High)*256 + float32(phCalibRawAd0Low)) / 100
			} else { //Negative ph calibration
				data = ((float32(phCalibRawAd0High)*256 + float32(phCalibRawAd0Low)) - 65535 - 1) / 100
			}
			ph.Data.PhCalibrationRawAd0 = data

		case 0x31:
			phCalibRawAd1High, err := strconv.ParseInt(hexValue.(string)[:2], 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}
			phCalibRawAd1Low, err := strconv.ParseInt(hexValue.(string)[2:], 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}
			var data float32
			if hexValue.(string)[:2] != "ff" { // Positive ph calibration
				data = (float32(phCalibRawAd1High)*256 + float32(phCalibRawAd1Low)) / 100
			} else { //Negative ph calibration
				data = ((float32(phCalibRawAd1High)*256 + float32(phCalibRawAd1Low)) - 65535 - 1) / 100
			}
			ph.Data.PhCalibrationRawAd1 = data

		case 0x32:
			phCalibRawAd2High, err := strconv.ParseInt(hexValue.(string)[:2], 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}
			phCalibRawAd2Low, err := strconv.ParseInt(hexValue.(string)[2:], 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}
			var data float32
			if hexValue.(string)[:2] != "ff" { // Positive ph calibration
				data = (float32(phCalibRawAd2High)*256 + float32(phCalibRawAd2Low)) / 100
			} else { //Negative ph calibration
				data = ((float32(phCalibRawAd2High)*256 + float32(phCalibRawAd2Low)) - 65535 - 1) / 100
			}
			ph.Data.PhCalibrationRawAd2 = data

		case 0x0200:
			modbusAddr, err := strconv.ParseInt(hexValue.(string), 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}
			ph.Data.ModbusSlaveAddr = int16(modbusAddr)

		case 0x0201:
			baudRate, err := strconv.ParseInt(hexValue.(string), 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}
			switch baudRate {
			case 0:
				ph.Data.BaudRate = 1200

			case 1:
				ph.Data.BaudRate = 2400

			case 2:
				ph.Data.BaudRate = 4800

			case 3:
				ph.Data.BaudRate = 9600

			case 4:
				ph.Data.BaudRate = 19200

			case 5:
				ph.Data.BaudRate = 38400
			}

		case 0x0202:
			serialProtocol, err := strconv.ParseInt(hexValue.(string), 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}
			if serialProtocol == 0 {
				ph.Data.SerialProtocol = "Modbus RTU"
			} else {
				ph.Data.SerialProtocol = "Modbus ASCII"
			}

		case 0x0203:
			parity, err := strconv.ParseInt(hexValue.(string), 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}
			if parity == 1 {
				ph.Data.ParityBit = "Even parity"
			} else if parity == 2 {
				ph.Data.ParityBit = "Odd parity"
			} else {
				ph.Data.ParityBit = "No parity"
			}

		case 0x0204:
			data, err := strconv.ParseInt(hexValue.(string), 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}
			if data == 1 {
				ph.Data.DataBits = "8 data bits"
			} else {
				ph.Data.DataBits = "Error, this device only supports 8 data bits format!"
			}

		case 0x0205:
			stop, err := strconv.ParseInt(hexValue.(string), 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}
			if stop == 0 {
				ph.Data.Stopbits = "1 stop bit"
			} else {
				ph.Data.Stopbits = "2 stop bits"
			}

		case 0x0206:
			responseDelay, err := strconv.ParseInt(hexValue.(string), 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}
			ph.Data.SerialResponseDelay = int16(responseDelay * 10)

		case 0x0207:
			outputInterval, err := strconv.ParseInt(hexValue.(string), 16, 64)
			if err != nil {
				log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
			}
			ph.Data.SerialResponseDelay = int16(outputInterval)

		default:
			return fmt.Errorf("default case entered for register: %d", reg)
		}

	}

	return nil
}
