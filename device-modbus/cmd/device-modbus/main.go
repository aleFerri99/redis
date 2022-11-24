// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2017-2018 Canonical Ltd
// Copyright (C) 2018-2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

// This package provides a simple example of a device service.
package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/edgexfoundry/device-modbus"
	"github.com/edgexfoundry/device-modbus/driver"
	"github.com/edgexfoundry/device-modbus/parser"
	"github.com/edgexfoundry/device-sdk-go/v2/pkg/startup"
)

const (
	serviceName string = "device-modbus"
)

var (
	verbose        = flag.Bool("verbose", false, "Add/Remove debug logs")
	serialResource = flag.String("serialProtocol", os.Getenv("SERIAL_RESOURCE"), "Serial resource used to open MODBUS communication")
	readGap        = flag.Duration("readGap", time.Duration(5)*time.Second, "Time span between modbus handling of consecutive sensors. Default: 5 seconds")
	confdir        = flag.String("confdir", "", "Path to EdgeX DS configuration files")
	err            error
)

func main() {

	// Get env vars
	*verbose, err = strconv.ParseBool(os.Getenv("VERBOSE"))
	if err != nil {
		log.Printf("Cannot parse %s to bool. Taking default value -> false...", os.Getenv("VERBOSE"))
		*verbose = false
	}

	*readGap, err = time.ParseDuration(os.Getenv("READ_GAP"))
	if err != nil {
		log.Printf("Cannot parse %s to time.Duration. Taking defualt value of 5 seconds...", os.Getenv("READ_GAP"))
		*readGap = time.Duration(5) * time.Second
	}
	sd := driver.SimpleDriver{}
	// Assign ENV variables to device struct
	sd.Verbose = *verbose
	sd.SerialResource = serialResource
	sd.ReadGap = readGap
	// Add YAML file parsing here, to avoid perform the same parsing step at each Device AutoEvent

	log.Printf("Using serial resource %s", *serialResource)
	log.Printf("Reading MODBUS configuration from %s", os.Getenv("MODBUS_FILE"))

	sd.ModbusConfig = &parser.Sensors{}
	sd.ModbusConfig.Parse(os.Getenv("MODBUS_FILE"), *verbose)
	if *verbose {
		prettyprint, err := json.MarshalIndent(&sd.ModbusConfig, "", "\t")
		if err != nil {
			log.Printf("Failed to pretty print modbus configration file. Error: %s", err)
			prettyprint = []byte("ERROR")
		}
		log.Printf("Pretty print MODBUS configuration file:\n%s", string(prettyprint))
	}
	startup.Bootstrap(serviceName, device.Version, &sd)
}
