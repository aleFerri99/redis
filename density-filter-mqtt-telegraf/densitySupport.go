package main

import (
	"math"
	"strconv"
	"time"
)

func (density *Density) ParseDataFromHex() error {

	// Store decoded data into data structure
	density.Data.DynamicViscosity = extractData(density.DataHexSensor["0x00"].(string)[0:4])
	density.Data.Density = extractData(density.DataHexSensor["0x00"].(string)[4:8])
	density.Data.Temperature = extractData(density.DataHexSensor["0x00"].(string)[8:])
	// Serial data cannot be retrived from sensor, so just set them as declared within documentation
	density.Data.BaudRate = 9600
	density.Data.ModbusSlaveAddr = 1
	density.Data.SerialProtocol = "Modbus RTU"
	density.Data.DataBits = "8 data bits"
	density.Data.ParityBit = "No parity"
	density.Data.Stopbits = "1 stop bit"
	// Not defined, so just put a non-sense value
	density.Data.SerialResponseDelay = -1
	density.Data.SerialActiveOutputTime = -1
	density.Data.Timestamp = time.Now().UnixMilli()

	return nil
}

func extractData(data string) float32 {
	n, err := strconv.ParseUint(data, 16, 32)
	if err != nil {
		panic(err)
	}
	n2 := uint32(n)
	return math.Float32frombits(n2)
}
