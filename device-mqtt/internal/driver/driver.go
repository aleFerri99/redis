// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2019-2021 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.

//UNCOMMENT FOR USING COMMAND

package driver

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	sdkModel "github.com/edgexfoundry/device-sdk-go/v2/pkg/models"
	"github.com/edgexfoundry/device-sdk-go/v2/pkg/service"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/errors"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/models"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/cast"
)

var once sync.Once
var driver *Driver

var (
	scheme        = flag.String("Scheme", os.Getenv("SCHEMEMQTTDEVICE"), "Protocol of the MQTT Broker")
	brokerUrl     = flag.String("BrokerUrl", os.Getenv("MQTTDEVICE_BROKER"), "Endpoint of the MQTT Broker")
	port          = flag.String("Port", os.Getenv("MQTTDEVICE_BROKER_PORT"), "Port of the MQTT Broker")
	mqttClientId  = flag.String("MqttClientId", os.Getenv("CLIENTID_MQTTDEVICE"), "Client ID used for Broker subscription")
	keepAlive     = flag.String("KeepAlive", os.Getenv("KEEPALIVE"), "keepalive of the MQTT Broker")
	qos           = flag.String("QoS", os.Getenv("QOSMQTTDEVICE"), "Type of message that is sent to the MQTT Broker. 0 -> Send without reception check, 1 -> Send more than once with reception check, 2 -> Send exactly once")
	incomingTopic = flag.String("IncomingTopic", os.Getenv("INCOMINGTOPIC"), "MQTT topic on which sensor data will be received")
	//responseFetchInterval = 500
	//credentialsRetryTime  = 120
	//credentialsRetryWait  = 1
	//connRetryWaitTime     = 5
	//useTopicLevels        = false
)

type Driver struct {
	Logger           logger.LoggingClient
	AsyncCh          chan<- *sdkModel.AsyncValues
	CommandResponses sync.Map
	mqttClient       mqtt.Client
	Verbose          bool
	SerialResource   *string
}

func NewProtocolDriver() sdkModel.ProtocolDriver {
	once.Do(func() {
		driver = new(Driver)
	})
	return driver
}

func (d *Driver) Initialize(lc logger.LoggingClient, asyncCh chan<- *sdkModel.AsyncValues, deviceCh chan<- []sdkModel.DiscoveredDevice) error {
	d.Logger = lc
	d.AsyncCh = asyncCh

	service.RunningService()

	client, err := d.createMqttClient()
	if err != nil {
		return errors.NewCommonEdgeX(errors.Kind(err), "unable to initial the MQTT client", err)
	}
	d.mqttClient = client

	return nil
}

func (d *Driver) DisconnectDevice(deviceName string, protocols map[string]models.ProtocolProperties) error {
	d.Logger.Warn("Driver's DisconnectDevice function didn't implement")
	return nil
}

func (d *Driver) HandleReadCommands(deviceName string, protocols map[string]models.ProtocolProperties, reqs []sdkModel.CommandRequest) ([]*sdkModel.CommandValue, error) {
	var responses = make([]*sdkModel.CommandValue, len(reqs))
	//commandTopic, err := fetchCommandTopic(protocols)
	// if err != nil {
	// 	return responses, err
	// }

	// for i, req := range reqs {
	// 	res, err := d.handleReadCommandRequest(req, commandTopic)
	// 	if err != nil {
	// 		driver.Logger.Infof("Handle read commands failed: %v", err)
	// 		return responses, err
	// 	}

	// 	responses[i] = res
	// }

	return responses, nil
}

// func (d *Driver) handleReadCommandRequest(req sdkModel.CommandRequest, topic string) (*sdkModel.CommandValue, error) {
// 	var result = &sdkModel.CommandValue{}
// 	var err error
// 	var qos = byte(0)
// 	var retained = false

// 	var method = "get"
// 	var cmdUuid = uuid.NewString()

// 	var cmd = req.DeviceResourceName
// 	var payload []byte

// 	if d.serviceConfig.MQTTBrokerInfo.useTopicLevels {
// 		topic = fmt.Sprintf("%s/%s/%s/%s", topic, cmd, method, cmdUuid)
// 		// will publish empty payload
// 	} else {
// 		data := make(map[string]interface{})
// 		data["uuid"] = cmdUuid
// 		data["method"] = method
// 		data["cmd"] = cmd

// 		payload, err = json.Marshal(data)
// 		if err != nil {
// 			return result, err
// 		}
// 	}

// 	driver.mqttClient.Publish(topic, qos, retained, payload)

// 	driver.Logger.Debugf("Publish command: %v", string(payload))

// 	// fetch response from MQTT broker after publish command successful
// 	cmdResponse, ok := d.fetchCommandResponse(cmdUuid)
// 	if !ok {
// 		return nil, errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("can not fetch command response: method=%v cmd=%v", method, cmd), nil)
// 	}

// 	driver.Logger.Debugf("Parse command response: %v", cmdResponse)

// 	var response map[string]interface{}
// 	err = json.Unmarshal([]byte(cmdResponse), &response)
// 	if err != nil {
// 		driver.Logger.Errorf("Error unmarshaling response: %s", err)
// 		return nil, errors.NewCommonEdgeX(errors.KindIOError, "Error umarshaling the response", err)
// 	}

// 	reading, ok := response[cmd]
// 	if !ok {
// 		return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("'%s' field not found in the response %s", cmd, cmdResponse), nil)
// 	}

// 	result, err = newResult(req, reading)
// 	if err != nil {
// 		return nil, errors.NewCommonEdgeXWrapper(err)
// 	}
// 	driver.Logger.Debugf("Get command finished: %v", result)

// 	return result, nil
// }

func (d *Driver) HandleWriteCommands(deviceName string, protocols map[string]models.ProtocolProperties, reqs []sdkModel.CommandRequest, params []*sdkModel.CommandValue) error {
	// commandTopic, err := fetchCommandTopic(protocols)
	// if err != nil {
	// 	return errors.NewCommonEdgeXWrapper(err)
	// }

	// for i, req := range reqs {
	// 	err = d.handleWriteCommandRequest(req, commandTopic, params[i])
	// 	if err != nil {
	// 		driver.Logger.Errorf("Handle write commands failed: %v", err)
	// 		return err
	// 	}
	// }

	return nil
}

func (d *Driver) handleWriteCommandRequest(req sdkModel.CommandRequest, topic string, param *sdkModel.CommandValue) errors.EdgeX {
	// var qos = byte(0)
	// var retained = false

	// var method = "set"
	// var cmdUuid = uuid.NewString()
	// var cmd = req.DeviceResourceName
	// var payload []byte
	// data := make(map[string]interface{})

	// commandValue, err := newCommandValue(req.Type, param)
	// if err != nil {
	// 	return errors.NewCommonEdgeXWrapper(err)
	// }
	// if d.serviceConfig.MQTTBrokerInfo.useTopicLevels {
	// 	topic = fmt.Sprintf("%s/%s/%s/%s", topic, cmd, method, cmdUuid)
	// 	data[cmd] = commandValue
	// } else {
	// 	data["uuid"] = cmdUuid
	// 	data["method"] = method
	// 	data["cmd"] = cmd
	// 	data[cmd] = commandValue
	// }

	// payload, err = json.Marshal(data)
	// if err != nil {
	// 	return errors.NewCommonEdgeXWrapper(err)
	// }
	// driver.mqttClient.Publish(topic, qos, retained, payload)

	// driver.Logger.Debugf("Publish command: %v", string(payload))

	// //wait and fetch response from CommandResponses map
	// cmdResponse, ok := d.fetchCommandResponse(cmdUuid)
	// if !ok {
	// 	return errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("can not fetch command response: method=%v cmd=%v", method, cmd), nil)
	// }

	// driver.Logger.Debugf("Put command finished: %v", cmdResponse)

	return nil
}

func (d *Driver) Stop(force bool) error {
	d.Logger.Info("driver is stopping, disconnect the MQTT conn")
	if d.mqttClient.IsConnected() {
		d.mqttClient.Disconnect(5000)
	}
	return nil
}

func newResult(req sdkModel.CommandRequest, reading interface{}) (*sdkModel.CommandValue, error) {
	var err error
	var result = &sdkModel.CommandValue{}
	castError := "fail to parse %v reading, %v"

	// if !checkValueInRange(req.Type, reading) {
	// 	err = fmt.Errorf("parse reading fail. Reading %v is out of the value type(%v)'s range", reading, req.Type)
	// 	driver.Logger.Error(err.Error())
	// 	return result, err
	// }

	var val interface{}
	switch req.Type {
	case common.ValueTypeBool:
		val, err = cast.ToBoolE(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, req.DeviceResourceName, err)
		}
	case common.ValueTypeString:
		val, err = cast.ToStringE(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, req.DeviceResourceName, err)
		}
	case common.ValueTypeUint8:
		val, err = cast.ToUint8E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, req.DeviceResourceName, err)
		}
	case common.ValueTypeUint16:
		val, err = cast.ToUint16E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, req.DeviceResourceName, err)
		}
	case common.ValueTypeUint32:
		val, err = cast.ToUint32E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, req.DeviceResourceName, err)
		}
	case common.ValueTypeUint64:
		val, err = cast.ToUint64E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, req.DeviceResourceName, err)
		}
	case common.ValueTypeInt8:
		val, err = cast.ToInt8E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, req.DeviceResourceName, err)
		}
	case common.ValueTypeInt16:
		val, err = cast.ToInt16E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, req.DeviceResourceName, err)
		}
	case common.ValueTypeInt32:
		val, err = cast.ToInt32E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, req.DeviceResourceName, err)
		}
	case common.ValueTypeInt64:
		val, err = cast.ToInt64E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, req.DeviceResourceName, err)
		}
	case common.ValueTypeFloat32:
		val, err = cast.ToFloat32E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, req.DeviceResourceName, err)
		}
	case common.ValueTypeFloat64:
		val, err = cast.ToFloat64E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, req.DeviceResourceName, err)
		}
	case common.ValueTypeObject:
		val = reading
	default:
		val = reading

	}

	result, err = sdkModel.NewCommandValue(req.DeviceResourceName, req.Type, val)
	if err != nil {
		return nil, err
	}
	result.Origin = time.Now().UnixNano()

	return result, nil
}

func newCommandValue(valueType string, param *sdkModel.CommandValue) (interface{}, error) {
	var commandValue interface{}
	var err error
	switch valueType {
	case common.ValueTypeBool:
		commandValue, err = param.BoolValue()
	case common.ValueTypeString:
		commandValue, err = param.StringValue()
	case common.ValueTypeUint8:
		commandValue, err = param.Uint8Value()
	case common.ValueTypeUint16:
		commandValue, err = param.Uint16Value()
	case common.ValueTypeUint32:
		commandValue, err = param.Uint32Value()
	case common.ValueTypeUint64:
		commandValue, err = param.Uint64Value()
	case common.ValueTypeInt8:
		commandValue, err = param.Int8Value()
	case common.ValueTypeInt16:
		commandValue, err = param.Int16Value()
	case common.ValueTypeInt32:
		commandValue, err = param.Int32Value()
	case common.ValueTypeInt64:
		commandValue, err = param.Int64Value()
	case common.ValueTypeFloat32:
		commandValue, err = param.Float32Value()
	case common.ValueTypeFloat64:
		commandValue, err = param.Float64Value()
	case common.ValueTypeObject:
		commandValue, err = param.ObjectValue()
	default:
		err = fmt.Errorf("fail to convert param, none supported value type: %v", valueType)
	}

	return commandValue, err
}

// // fetchCommandResponse use to wait and fetch response from CommandResponses map
// func (d *Driver) fetchCommandResponse(cmdUuid string) (string, bool) {
// 	var cmdResponse interface{}
// 	var ok bool
// 	for i := 0; i < 5; i++ {
// 		cmdResponse, ok = d.CommandResponses.Load(cmdUuid)
// 		if ok {
// 			d.CommandResponses.Delete(cmdUuid)
// 			break
// 		} else {
// 			time.Sleep(time.Millisecond * time.Duration(d.serviceConfig.MQTTBrokerInfo.Writable.responseFetchInterval))
// 		}
// 	}

// 	return fmt.Sprintf("%v", cmdResponse), ok
// }

func (d *Driver) AddDevice(deviceName string, protocols map[string]models.ProtocolProperties, adminState models.AdminState) error {
	d.Logger.Debugf("Device %s is added", deviceName)
	return nil
}

func (d *Driver) UpdateDevice(deviceName string, protocols map[string]models.ProtocolProperties, adminState models.AdminState) error {
	d.Logger.Debugf("Device %s is updated", deviceName)
	return nil
}

func (d *Driver) RemoveDevice(deviceName string, protocols map[string]models.ProtocolProperties) error {
	d.Logger.Debugf("Device %s is removed", deviceName)
	return nil
}

func (d *Driver) createMqttClient() (mqtt.Client, errors.EdgeX) {

	uri := &url.URL{
		Scheme: strings.ToLower(*scheme),
		Host:   fmt.Sprintf("%s:%s", *brokerUrl, *port),
	}

	var client mqtt.Client
	var err error
	ka, _ := strconv.Atoi(*keepAlive)
	for {
		client, err = d.getMqttClient(*mqttClientId, uri, ka)
		if err != nil {
			driver.Logger.Warnf("Unable to connect to MQTT broker, %s, retrying", err)
			time.Sleep(5 * time.Second)
			continue
		}
		return client, nil
	}
}

func (d *Driver) getMqttClient(clientID string, uri *url.URL, keepAlive int) (mqtt.Client, error) {
	driver.Logger.Infof("Create MQTT client and connection: uri=%v clientID=%v ", uri.String(), clientID)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("%s://%s", uri.Scheme, uri.Host))
	opts.SetClientID(clientID)
	opts.SetUsername(uri.User.Username())
	password, _ := uri.User.Password()
	opts.SetPassword(password)
	opts.SetKeepAlive(time.Second * time.Duration(keepAlive))
	opts.SetAutoReconnect(true)
	opts.OnConnect = d.onConnectHandler

	client := mqtt.NewClient(opts)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		return client, token.Error()
	}

	return client, nil
}

func (d *Driver) onConnectHandler(client mqtt.Client) {
	qs, _ := strconv.Atoi(*qos)
	token := client.Subscribe(*incomingTopic, byte(qs), d.onIncomingDataReceived)
	if token.Wait() && token.Error() != nil {
		client.Disconnect(0)
		driver.Logger.Errorf("could not subscribe to topic '%s': %s",
			incomingTopic, token.Error().Error())
		return
	}
	driver.Logger.Infof("Subscribed to topic '%s' for receiving the async reading", incomingTopic)

	// token = client.Subscribe(responseTopic, qos, d.onCommandResponseReceived)
	// if token.Wait() && token.Error() != nil {
	// 	client.Disconnect(0)
	// 	driver.Logger.Errorf("could not subscribe to topic '%s': %s",
	// 		responseTopic, token.Error().Error())
	// 	return
	// }
	// driver.Logger.Infof("Subscribed to topic '%s' for receiving the request response", responseTopic)

}
