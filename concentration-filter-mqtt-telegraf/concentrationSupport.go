package main

import (
	"fmt"
	"time"
)

func (concentration *Concentration) ParseDataFromHex() error {

	for register, value := range concentration.DataSensor {
		// reg, err := strconv.Atoi(register)
		// if err != nil {
		// 	log.Printf("Cannot parse string to integer. Error: %s", err)
		// 	return err
		// }
		switch register {
		case "temp":
			concentration.Data.Temperature = value.(float64)
		case "freq":
			//var data []float64

			concentration.Data.Frequence = value.([]interface{})
		case "phase":
			// var data []float64
			// data = hexValue.([]float64)
			concentration.Data.Phase = value.([]interface{})
		case "imp":
			// var data []float64
			// data = hexValue.([]float64)
			concentration.Data.Impedance = value.([]interface{})
		case "name":
			concentration.SensorName = value.(string)

		// case 0x02:
		// 	phCalibRawAdHigh, err := strconv.ParseInt(hexValue.(string)[:2], 16, 64)
		// 	if err != nil {
		// 		log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
		// 	}
		// 	phCalibRawAdLow, err := strconv.ParseInt(hexValue.(string)[2:], 16, 64)
		// 	if err != nil {
		// 		log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
		// 	}
		// 	var data float32
		// 	if hexValue.(string)[:2] != "ff" { // Positive ph calibration
		// 		data = (float32(phCalibRawAdHigh)*256 + float32(phCalibRawAdLow)) / 100
		// 	} else { //Negative ph calibration
		// 		data = ((float32(phCalibRawAdHigh)*256 + float32(phCalibRawAdLow)) - 65535 - 1) / 100
		// 	}
		// 	conc.Data.Impedance = data

		// case 0x30:
		// 	phCalibRawAd0High, err := strconv.ParseInt(hexValue.(string)[:2], 16, 64)
		// 	if err != nil {
		// 		log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
		// 	}
		// 	phCalibRawAd0Low, err := strconv.ParseInt(hexValue.(string)[2:], 16, 64)
		// 	if err != nil {
		// 		log.Printf("Cannot parse HEX string to int. Error: %x\n", err)
		// 	}
		// 	var data float32
		// 	if hexValue.(string)[:2] != "ff" { // Positive ph calibration
		// 		data = (float32(phCalibRawAd0High)*256 + float32(phCalibRawAd0Low)) / 100
		// 	} else { //Negative ph calibration
		// 		data = ((float32(phCalibRawAd0High)*256 + float32(phCalibRawAd0Low)) - 65535 - 1) / 100
		// 	}
		// 	conc.Data.Phase = data

		default:
			return fmt.Errorf("default case entered for register: %s", register)
		}
		concentration.Data.Timestamp = time.Now().Unix()

	}

	return nil
}
