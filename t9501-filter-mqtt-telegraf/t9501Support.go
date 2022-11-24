package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
)

func (t9501 *T9501) ParseDataFromHex() error {

	for register, hexValue := range t9501.DataHexSensor {
		reg, err := strconv.Atoi(register)
		if err != nil {
			log.Printf("Cannot parse string to integer. Error: %s", err)
			return err
		}
		switch reg {
		case 0x1389:
			fwVersion, err := strconv.ParseInt(string(hexValue.(string)), 16, 64)
			if err != nil {
				fmt.Printf("Cannot parse Hex data int. Error: %s", err)
				return err
			}
			t9501.Data.FirmwareRevision = float32(fwVersion) / 100
		case 0x138D:
			temperature, err := strconv.ParseInt(string(hexValue.(string)), 16, 64)
			if err != nil {
				fmt.Printf("Cannot parse Hex data int. Error: %s", err)
				return err
			}
			t9501.Data.Temperature = float32(temperature) / 100
		case 0x138E:
			relHumidity, err := strconv.ParseInt(string(hexValue.(string)), 16, 64)
			if err != nil {
				fmt.Printf("Cannot parse Hex data int. Error: %s", err)
				return err
			}
			t9501.Data.RelativeHumidity = float32(relHumidity) / 100
		case 0x0FA4:
			configFlag, err := strconv.ParseInt(string(hexValue.(string)), 16, 64)
			if err != nil {
				fmt.Printf("Cannot parse Hex data int. Error: %s", err)
				return err
			}
			t9501.Data.ConfigFlag = uint16(configFlag)
		case 0x0FA5:
			slaveAddress, err := strconv.ParseInt(string(hexValue.(string)), 16, 64)
			if err != nil {
				fmt.Printf("Cannot parse Hex data int. Error: %s", err)
				return err
			}
			t9501.Data.ModbusSlaveAddr = int16(slaveAddress)
		case 0x0FA6:
			baud, err := strconv.ParseInt(string(hexValue.(string)), 16, 64)
			if err != nil {
				fmt.Printf("Cannot parse Hex data int. Error: %s", err)
				return err
			}
			switch baud {
			case 1:
				t9501.Data.BaudRate = 9600
			case 2:
				t9501.Data.BaudRate = 19200
			case 3:
				t9501.Data.BaudRate = 38400
			case 4:
				t9501.Data.BaudRate = 57600
			case 5:
				t9501.Data.BaudRate = 115200
			default:
				fmt.Printf("Unknown baud rate value %d", baud)
				return errors.New("Unknown baud rate value " + strconv.Itoa(int(baud)))
			}
		case 0x0FA9:
			parity, err := strconv.ParseInt(string(hexValue.(string)), 16, 64)
			if err != nil {
				fmt.Printf("Cannot parse Hex data int. Error: %s", err)
				return err
			}
			switch parity {
			case 0:
				t9501.Data.SerialProtocol = "Modbus RTU"
				t9501.Data.ParityBit = "None"
				t9501.Data.DataBits = "8 data bits"
				t9501.Data.Stopbits = "2 stop bits"
				t9501.Data.SerialActiveOutputTime = -1
				t9501.Data.SerialResponseDelay = -1
			case 1:
				t9501.Data.SerialProtocol = "Modbus RTU"
				t9501.Data.ParityBit = "Odd"
				t9501.Data.DataBits = "8 data bits"
				t9501.Data.Stopbits = "1 stop bit"
				t9501.Data.SerialActiveOutputTime = -1
				t9501.Data.SerialResponseDelay = -1
			case 2:
				t9501.Data.SerialProtocol = "Modbus RTU"
				t9501.Data.ParityBit = "Even"
				t9501.Data.DataBits = "8 data bits"
				t9501.Data.Stopbits = "1 stop bit"
				t9501.Data.SerialActiveOutputTime = -1
				t9501.Data.SerialResponseDelay = -1
			default:
				fmt.Printf("Unknown Modbus RTU configuration %d", parity)
				return fmt.Errorf("unknown Modbus RTU configuration %s", strconv.Itoa(int(parity)))
			}
		default:
			fmt.Printf("Unknown register %d", reg)
			return fmt.Errorf("unknown register %d", reg)
		}

	}

	return nil
}
