//
// Copyright (c) 2021 Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"app-concentration-filter-mqtt-telegraf/mqttHandler"
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/interfaces"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/transforms"
)

const (
	serviceKey = "app-concentration-filter-mqtt-telegraf"
)

var (
	protocol            = flag.String("Protocol", os.Getenv("PROTOCOL"), "Protocol of the MQTT Broker")
	broker              = flag.String("Broker", os.Getenv("BROKER"), "Endpoint of the MQTT Broker")
	port                = flag.String("Port", os.Getenv("PORT"), "Port of the MQTT Broker")
	clientId            = flag.String("Client ID", os.Getenv("CLIENTID_CONCENTRATION"), "Client ID used for Broker subscription")
	user                = flag.String("User", os.Getenv("USERNAME"), "Username to connect to the MQTT Broker")
	password            = flag.String("Password", os.Getenv("PASSWORD"), "Password to connect to the MQTT Broker")
	qos                 = flag.String("QoS", os.Getenv("QOS"), "Type of message that is sent to the MQTT Broker. 0 -> Send without reception check, 1 -> Send more than once with reception check, 2 -> Send exactly once")
	retained            = flag.String("Retained", os.Getenv("RETAINED"), "Enstablish if messages published by the client are retained")
	autoreconnect       = flag.String("Autoreconnect", os.Getenv("AUTORECONNECT"), "Include automatic reconnect to MQTT Broker")
	cleanSession        = flag.String("Cleansession", os.Getenv("CLEANSESSION"), "If true, make the session between client and broker not persistent")
	store               = flag.String("Store", os.Getenv("STORE"), "Path to which persistent data must be saved")
	connectRety         = flag.String("Connectrety", os.Getenv("CONNECTRETRY"), "Maximum attempts of cennection retry")
	connectRetyInterval = flag.String("Connectretryinterval", os.Getenv("CONNECTRETRYINTERVAL"), "Time interval between consecutive reconnection attempts")
	TLSConnection       = flag.String("TLSConnection", os.Getenv("TLSCONNECTION"), "Variable that allow to include TLS handshake during connection or neglect it")
	clientAuthority     = flag.String("Clientauthority", os.Getenv("CLIENTAUTHORITY"), "Client certificate authority")
	clientCertificate   = flag.String("Clientcertificate", os.Getenv("CLIENTCERTIFICATE"), "Client certificate")
	clientKey           = flag.String("Clientkey", os.Getenv("CLIENTKEY"), "Client key")
	ignore              = flag.Bool("Ignore", false, "Skip the MQTT publishing if connection to broker returns error")
	topic_sensor        = flag.String("topic", os.Getenv("TOPIC_CONCENTRATION"), "MQTT topic on which sensor data will be published")
	topic_prefix        = flag.String("topic_prefix", os.Getenv("TOPIC_PREFIX"), "MQTT topic prefix on which sensor data will be published")
	msgChan             = make(chan mqttHandler.MsgToSend, 200) // GoRoutine for Asyn MQTT publishing. Allow a queue of up to 200 messages
	verbose             = flag.Bool("verbose", false, "Print/Hide logs")
	confdir             = flag.String("confdir", "", "EdgeX configuration directory")
	redisAddress        = flag.String("Redis Address", os.Getenv("REDISADDRESS"), "Redis address")
	redisPass           = flag.String("Redis Password", os.Getenv("REDISPASS"), "Redis password")
)

func main() {
	// turn off secure mode for examples. Not recommended for production
	_ = os.Setenv("EDGEX_SECURITY_SECRET_STORE", "false")
	verbose, err := strconv.ParseBool(os.Getenv("VERBOSE"))
	if err != nil {
		log.Printf("Cannot parse value %v to bool. Taking default...", os.Getenv("VERBOSE"))
		verbose = false
	}
	// Setup redis client
	redisClient := createRedisClientConnection(*redisAddress, *redisPass, 0)

	// Setup MQTT client
	mqttClient := mqttHandler.SetupClient(
		*protocol,
		*broker,
		*port,
		*clientId,
		*user,
		*password,
		*qos,
		*retained,
		*autoreconnect,
		*cleanSession,
		*store,
		*connectRety,
		*connectRetyInterval,
		*TLSConnection,
		*clientAuthority,
		*clientCertificate,
		*clientKey,
	)

	err = mqttHandler.Connect(mqttClient)
	if err != nil {
		log.Println("MQTT publishing will be ignored!")
		*ignore = true
	}

	// Questa va istanziata solo all'inzio dell'Application Service, non ad ogni nuovo pacchetto dal Device Service
	go func() {
		for msgStruct := range msgChan {
			// Send data to server
			mqttHandler.Publish(msgStruct.Topic, mqttClient, msgStruct.Payload)
		}
	}()

	// 1) First thing to do is to create an new instance of an EdgeX Application Service.
	service, ok := pkg.NewAppService(serviceKey)
	if !ok {
		os.Exit(-1)
	}

	// Leverage the built in logging service in EdgeX
	lc := service.LoggingClient()

	// 2) shows how to access the application's specific configuration settings.
	deviceNames, err := service.GetAppSettingStrings("DeviceNames")
	if err != nil {
		lc.Error(err.Error())
		os.Exit(-1)
	}
	lc.Info(fmt.Sprintf("Filtering for devices %v", deviceNames))

	// 3) This is our functions pipeline configuration, the collection of functions to
	// execute every time an event is triggered.
	if err := service.SetFunctionsPipeline(
		transforms.NewFilterFor(deviceNames).FilterByDeviceName,
		transforms.NewConversion().TransformToJSON,
		NewConcentration(&verbose).ParseSensorData,
		sendToServer,
		redisClient.saveToRedis,
	); err != nil {
		lc.Error("SetFunctionsPipeline failed: " + err.Error())
		os.Exit(-1)
	}

	// 4) Lastly, we'll go ahead and tell the SDK to "start" and begin listening for events
	// to trigger the pipeline.
	err = service.MakeItRun()
	if err != nil {
		lc.Error("MakeItRun returned error: ", err.Error())
		os.Exit(-1)
	}

	// Do any required cleanup here
	mqttHandler.Disconnect(mqttClient)
	os.Exit(0)
}

func sendToServer(ctx interfaces.AppFunctionContext, data interface{}) (continuePipeline bool, result interface{}) {

	if *verbose {
		prettyPrint, err := json.MarshalIndent(data.(map[string]interface{}), "", "\t")
		if err != nil {
			log.Printf("Cannot pretty print MQTT message: %s", err)
		} else {
			log.Printf("Sending data to MQTT broker.\nMessage: %s", prettyPrint)
		}
	}
	fmt.Println(data)

	for deviceName, device := range data.(map[string]interface{}) {
		// Create JSON object for the device i
		dataToSend, err := json.Marshal(device)
		if err != nil {
			return false, err
		}
		fmt.Println(deviceName)
		fmt.Println(string(dataToSend))
		serial := getserial()
		fmt.Println(serial)
		msgChan <- mqttHandler.MsgToSend{
			DeviceName: deviceName,
			Payload:    string(dataToSend),
			Topic:      *topic_prefix + serial + *topic_sensor,
		}
	}
	return true, data.(map[string]interface{})
}

func getserial() string {
	//Extract serial from cpuinfo file
	cpuserial := "0000000000000000"
	f, err := os.Open("/proc/cpuinfo")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		info_name := scanner.Text()
		if strings.Contains(info_name, "Serial") {
			cpuserial = strings.TrimSpace(strings.Split(info_name, ":")[1])
			break
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		cpuserial = "ERROR000000000"
	}
	return cpuserial
}
