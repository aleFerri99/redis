package parser

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Sensors struct {
	Sensors map[string]Sensor `yaml:"sensors"`
}

type Sensor struct {
	Name           string      `yaml:"name"`
	Poll           bool        `yaml:"poll"`
	Mode           string      `yaml:"mode"`
	ModbusProperty interface{} `yaml:"modbus-properties"`
	Modbus         Modbus      `yaml:"modbus"`
}

type SerialProperties struct {
	Address  uint   `yaml:"address"`
	Baudrate int    `yaml:"baudrate"`
	DataBits uint   `yaml:"databits"`
	StopBits uint   `yaml:"stopbits"`
	Parity   string `yaml:"parity"`
	Timeout  int    `yaml:"timeout"`
}

type AsciiProperties struct {
	Address  uint   `yaml:"address"`
	Baudrate int    `yaml:"baudrate"`
	DataBits uint   `yaml:"databits"`
	StopBits uint   `yaml:"stopbits"`
	Parity   string `yaml:"parity"`
	Timeout  int    `yaml:"timeout"`
}

type TcpProperties struct {
	Address uint `yaml:"address"`
	Port    int  `yaml:"port"`
	Timeout int  `yaml:"timeout"`
}

type Modbus struct {
	Registers []int       `yaml:"registers"`
	ReadLen   []int       `yaml:"read-length"`
	ReadWrite []ReadWrite `yaml:"read-write"`
}

type ReadWrite struct {
	FunctionCode int   `yaml:"function-code"`
	Value        []int `yaml:"value"`
}

func (modbusRtuList *Sensors) Parse(fileName string, verbose bool) error {

	if verbose {
		log.Println(`Parser default options:
	NAME: "",
	POLL: false,
	SERIAL-PROPERTIES:
		ADDRESS: 0x00
		BAUDRATE: 9600
		DATABITS: 8
		STOPBITS: 1
		PARITY: 'N'
		TIMEOUT: 5
	MODBUS_RTU:
		REGISTERS: []
		READ-LENGTH: []
		READ_WRITE:
			FUNCTION_CODE: 0
			VALUE: []
		`)
	}

	yamlFile, err := os.ReadFile(fileName)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, &modbusRtuList)
	if err != nil {
		log.Printf("Cannot unmarshal YAML file. Error: %s", err)
		return err
	}

	// Perform some file integrity checks
	for sensor, properties := range modbusRtuList.Sensors {
		if len(properties.Modbus.ReadLen) !=
			len(properties.Modbus.Registers) ||
			len(properties.Modbus.ReadLen) !=
				len(properties.Modbus.ReadWrite) {
			log.Printf(`
					Wrong YAML definition for sensor %s.
						len('registers') = %d
						len('read-length') = %d
						len('function-codes') = %d
						`,
				sensor,
				len(properties.Modbus.Registers),
				len(properties.Modbus.ReadLen),
				len(properties.Modbus.ReadWrite),
			)

			return fmt.Errorf(
				`
					Wrong YAML definition for sensor %s.
						len('registers') = %d
						len('read-length') = %d
						len('function-codes') = %d
						`,
				sensor,
				len(properties.Modbus.Registers),
				len(properties.Modbus.ReadLen),
				len(properties.Modbus.ReadWrite),
			)
		}
	}

	return nil
}
