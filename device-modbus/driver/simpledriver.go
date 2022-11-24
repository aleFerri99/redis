// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 Canonical Ltd
// Copyright (C) 2018-2021 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

// This package provides a simple example implementation of
// ProtocolDriver interface.
package driver

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/edgexfoundry/device-modbus/parser"
	"github.com/edgexfoundry/device-sdk-go/v2/pkg/interfaces"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/models"

	"github.com/edgexfoundry/device-sdk-go/v2/example/config"
	sdkModels "github.com/edgexfoundry/device-sdk-go/v2/pkg/models"
	"github.com/edgexfoundry/device-sdk-go/v2/pkg/service"
)

type SimpleDriver struct {
	lc             logger.LoggingClient
	asyncCh        chan<- *sdkModels.AsyncValues
	deviceCh       chan<- []sdkModels.DiscoveredDevice
	ModbusConfig   *parser.Sensors
	Verbose        bool
	SerialResource *string
	ReadGap        *time.Duration
	modbus         []map[string]interface{}
	serviceConfig  *config.ServiceConfig
}

type channelMsg struct {
	deviceName string
	req        sdkModels.CommandRequest
	command    bool
}

const (
	MAX_RETRY = 5
)

var command = make(chan channelMsg, 2)

// Initialize performs protocol-specific initialization for the device
// service.
func (s *SimpleDriver) Initialize(lc logger.LoggingClient, asyncCh chan<- *sdkModels.AsyncValues, deviceCh chan<- []sdkModels.DiscoveredDevice) error {
	s.lc = lc
	s.asyncCh = asyncCh
	s.deviceCh = deviceCh
	s.serviceConfig = &config.ServiceConfig{}

	ds := service.RunningService()

	if err := ds.LoadCustomConfig(s.serviceConfig, "SimpleCustom"); err != nil {
		return fmt.Errorf("unable to load 'SimpleCustom' custom configuration: %s", err.Error())
	}

	lc.Infof("Custom config is: %v", s.serviceConfig.SimpleCustom)

	if err := s.serviceConfig.SimpleCustom.Validate(); err != nil {
		return fmt.Errorf("'SimpleCustom' custom configuration validation failed: %s", err.Error())
	}

	if err := ds.ListenForCustomConfigChanges(
		&s.serviceConfig.SimpleCustom.Writable,
		"SimpleCustom/Writable", s.ProcessCustomConfigChanges); err != nil {
		return fmt.Errorf("unable to listen for changes for 'SimpleCustom.Writable' custom configuration: %s", err.Error())
	}

	// Launch go routine triggered by channel events
	go s.handleReadAndSend(command)

	registered := interfaces.DeviceServiceSDK.Devices(interfaces.Service())
	for _, device := range registered {
		log.Printf("Device: %v", device)
	}

	return nil
}

// ProcessCustomConfigChanges ...
func (s *SimpleDriver) ProcessCustomConfigChanges(rawWritableConfig interface{}) {
	updated, ok := rawWritableConfig.(*config.SimpleWritable)
	if !ok {
		s.lc.Error("unable to process custom config updates: Can not cast raw config to type 'SimpleWritable'")
		return
	}

	s.lc.Info("Received configuration updates for 'SimpleCustom.Writable' section")

	previous := s.serviceConfig.SimpleCustom.Writable
	s.serviceConfig.SimpleCustom.Writable = *updated

	if reflect.DeepEqual(previous, *updated) {
		s.lc.Info("No changes detected")
		return
	}

	// Now check to determine what changed.
	// In this example we only have the one writable setting,
	// so the check is not really need but left here as an example.
	// Since this setting is pulled from configuration each time it is need, no extra processing is required.
	// This may not be true for all settings, such as external host connection info, which
	// may require re-establishing the connection to the external host for example.
	if previous.DiscoverSleepDurationSecs != updated.DiscoverSleepDurationSecs {
		s.lc.Infof("DiscoverSleepDurationSecs changed to: %d", updated.DiscoverSleepDurationSecs)
	}
}

// HandleReadCommands triggers a protocol Read operation for the specified device.
func (s *SimpleDriver) HandleReadCommands(deviceName string, protocols map[string]models.ProtocolProperties, reqs []sdkModels.CommandRequest) (res []*sdkModels.CommandValue, err error) {
	s.lc.Debugf("SimpleDriver.HandleReadCommands: protocols: %v resource: %v attributes: %v", protocols, reqs[0].DeviceResourceName, reqs[0].Attributes)

	return nil, fmt.Errorf("RestDriver.HandleReadCommands; read commands not supported")
}

// HandleWriteCommands passes a slice of CommandRequest struct each representing
// a ResourceOperation for a specific device resource.
// Since the commands are actuation commands, params provide parameters for the individual
// command.
func (s *SimpleDriver) HandleWriteCommands(deviceName string, protocols map[string]models.ProtocolProperties, reqs []sdkModels.CommandRequest,
	params []*sdkModels.CommandValue) error {

	for i, req := range reqs {
		s.lc.Debugf("SimpleDriver.HandleWriteCommands: protocols: %v, resource: %v, parameters: %v, attributes: %v", protocols, reqs[i].DeviceResourceName, params[i], reqs[i].Attributes)
		select {
		case command <- channelMsg{
			deviceName: deviceName,
			req:        req,
			command:    params[0].Value.(bool),
		}:
			log.Printf("Pushed event to message channel")
		default:
			log.Printf("WARN: Channel is full!")
		}
	}

	return nil

}

func (s *SimpleDriver) handleReadAndSend(command chan channelMsg) {

	var cv *sdkModels.CommandValue
	var read channelMsg
	var trueEvent bool

	// Wait for device service to be available
	attempt := 0
	startPipeline := false
	for !startPipeline {
		//log.Printf("DEVICES: %v", interfaces.Service().Devices())
		//for _, device := range interfaces.Service().Devices() {
		//	err := interfaces.Service().UpdateDevice(device)
		//	if err != nil {
		//		log.Printf("Cannot update device %s in core MetaData and Cache. Error: %s", device.Name, err)
		//		time.Sleep(5 * time.Second)
		//		continue
		//	}
		//}
		ds := interfaces.Service()
		_, errModbus := ds.GetDeviceByName(ds.Name())
		if errModbus != nil {
			attempt++
			log.Printf("Attempt: %d. Device '%s' not available", attempt, ds.Name())
			if attempt > MAX_RETRY {
				os.Exit(0)
			}
			time.Sleep(3 * time.Second)
			continue
		}
		startPipeline = true
	}

readChannel: // Add label to continue for loop till first event is triggered
	for {
		select {
		case read = <-command:
			log.Printf("New incoming event. Updating data...")
			if read.command {
				trueEvent = true
			} else {
				trueEvent = false
			}
		default:
			if trueEvent {
				log.Printf("No new event. Continue reading data...")
			} else {
				time.Sleep(time.Second)
				continue readChannel
			}
		}
		switch read.req.DeviceResourceName {
		case "GPIO":
			if read.command {
				s.modbus = HandleSensors(s.ModbusConfig, s.SerialResource, s.Verbose, s.ReadGap)
				modbus, err := json.Marshal(s.modbus)
				if err != nil {
					log.Printf("Cannot parse MODBUS data to JSON. Error: %s", err)
					cv, _ = sdkModels.NewCommandValue(read.req.DeviceResourceName, common.ValueTypeString, err)
				} else {
					cv, _ = sdkModels.NewCommandValue(read.req.DeviceResourceName, common.ValueTypeString, string(modbus))
				}
				s.lc.Info(fmt.Sprintf("Data sent to core data: %s", string(modbus)))
			} else {
				log.Println("Stop readings, pump timeout reached...")
				cv, _ = sdkModels.NewCommandValue(read.req.DeviceResourceName, common.ValueTypeString, "Stop readings, pump timeout reached...")
			}
		default:
			log.Printf("Unknown Device Command %s.", read.req.DeviceResourceName)
			cv, _ = sdkModels.NewCommandValue(read.req.DeviceResourceName, common.ValueTypeString, "Unknown Device Command")
		}
		s.handleAsyncCommunication(read.deviceName, cv)
		time.Sleep(time.Second)
	}
}

func (s *SimpleDriver) handleAsyncCommunication(deviceName string, cv *sdkModels.CommandValue) {
	res := make([]*sdkModels.CommandValue, 1)
	log.Println("Pushing modbus readings to EdgeX Core Data")
	res[0] = cv
	asyncValues := &sdkModels.AsyncValues{
		DeviceName:    deviceName,
		CommandValues: res,
	}
	s.asyncCh <- asyncValues
}

// Stop the protocol-specific DS code to shutdown gracefully, or
// if the force parameter is 'true', immediately. The driver is responsible
// for closing any in-use channels, including the channel used to send async
// readings (if supported).
func (s *SimpleDriver) Stop(force bool) error {
	// Then Logging Client might not be initialized
	if s.lc != nil {
		s.lc.Debugf("SimpleDriver.Stop called: force=%v", force)
	}
	return nil
}

// AddDevice is a callback function that is invoked
// when a new Device associated with this Device Service is added
func (s *SimpleDriver) AddDevice(deviceName string, protocols map[string]models.ProtocolProperties, adminState models.AdminState) error {
	s.lc.Debugf("a new Device is added: %s", deviceName)
	return nil
}

// UpdateDevice is a callback function that is invoked
// when a Device associated with this Device Service is updated
func (s *SimpleDriver) UpdateDevice(deviceName string, protocols map[string]models.ProtocolProperties, adminState models.AdminState) error {
	s.lc.Debugf("Device %s is updated", deviceName)
	return nil
}

// RemoveDevice is a callback function that is invoked
// when a Device associated with this Device Service is removed
func (s *SimpleDriver) RemoveDevice(deviceName string, protocols map[string]models.ProtocolProperties) error {
	s.lc.Debugf("Device %s is removed", deviceName)
	return nil
}

// Discover triggers protocol specific device discovery, which is an asynchronous operation.
// Devices found as part of this discovery operation are written to the channel devices.
func (s *SimpleDriver) Discover() {
	proto := make(map[string]models.ProtocolProperties)
	proto["other"] = map[string]string{"Address": "simple02", "Port": "301"}

	device2 := sdkModels.DiscoveredDevice{
		Name:        "Simple-Device02",
		Protocols:   proto,
		Description: "found by discovery",
		Labels:      []string{"auto-discovery"},
	}

	proto = make(map[string]models.ProtocolProperties)
	proto["other"] = map[string]string{"Address": "simple03", "Port": "399"}

	device3 := sdkModels.DiscoveredDevice{
		Name:        "Simple-Device03",
		Protocols:   proto,
		Description: "found by discovery",
		Labels:      []string{"auto-discovery"},
	}

	res := []sdkModels.DiscoveredDevice{device2, device3}

	time.Sleep(time.Duration(s.serviceConfig.SimpleCustom.Writable.DiscoverSleepDurationSecs) * time.Second)
	s.deviceCh <- res
}
