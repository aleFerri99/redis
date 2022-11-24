package mqttHandler

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type mqttConfig struct {
	protocol            string
	broker              string
	port                string
	clientId            string
	user                string
	password            string
	qos                 int64
	retained            bool
	autoreconnect       bool
	cleanSession        bool
	store               string
	connectRety         bool
	connectRetyInterval int
	TLSConnection       bool
	clientAuthority     string
	clientCertificate   string
	clientKey           string
}

type MsgToSend struct {
	DeviceName string
	Payload    string
	Topic      string
}

var (
	mqttConfiguration mqttConfig
	err               error
)

/*
	Add generic methods to retrieve message and status about connection form MQTT callbacks
*/

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("TOPIC: %s", msg.Topic())
	log.Printf("MESSAGE: %s", msg.Payload())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("Connection lost: %v", err)
}

/*
	Define genic method to Connect, Subscribe, Publish or Disconnect form/to MQTT Broker
*/

func Connect(client mqtt.Client) error {
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Printf("Cannot connect to specified broker. \n\nError: %s", token.Error())
		return token.Error()
	}
	return nil
}

func Subscribe(topic string, client mqtt.Client) {
	if client.IsConnected() {
		token := client.Subscribe(topic, byte(mqttConfiguration.qos), nil) // No callback onSubscribe has been set
		token.Wait()
		log.Printf("Subscribed to topic %s.", topic)
	}
}

func Unsubscribe(topic string, client mqtt.Client) {
	if client.IsConnected() {
		if token := client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
			log.Printf("Error unsubcribing from topic %s:\n%s", topic, token.Error())
		}
	}
}

func Publish(topic string, client mqtt.Client, message interface{}) {
	if client.IsConnected() {

		token := client.Publish(topic, byte(mqttConfiguration.qos), mqttConfiguration.retained, message)
		token.Wait()
		time.Sleep(time.Second)
	}
}

func Disconnect(client mqtt.Client) {
	client.Disconnect(250)
}

/*
	Define method to setup a new MQTT Client object
*/

func SetupClient(
	protocol string,
	broker string,
	port string,
	clientId string,
	user string,
	password string,
	qos string,
	retained string,
	autoreconnect string,
	cleanSession string,
	store string,
	connectRety string,
	connectRetyInterval string,
	TLSConnection string,
	clientAuthority string,
	clientCertificate string,
	clientKey string,
) mqtt.Client {
	// Populate struct for the mqtt handler with ENV variables from main
	mqttConfiguration.protocol = protocol
	mqttConfiguration.broker = broker
	mqttConfiguration.port = port
	mqttConfiguration.clientId = clientId
	mqttConfiguration.user = user
	mqttConfiguration.password = password
	mqttConfiguration.qos, err = strconv.ParseInt(qos, 10, 64)
	if err != nil {
		log.Println("Cannot parse QoS to its relative integer value. Using default value.")
		mqttConfiguration.qos = 0
	}

	if strings.EqualFold(retained, "true") {
		mqttConfiguration.retained = true
	} else {
		// If value is omitted or is not true neither false, set by default its value to false
		mqttConfiguration.retained = false
	}

	if strings.EqualFold(autoreconnect, "true") {
		mqttConfiguration.autoreconnect = true
	} else {
		// If value is omitted or is not true neither false, set by default its value to false
		mqttConfiguration.autoreconnect = false
	}

	mqttConfiguration.store = store
	if strings.EqualFold(cleanSession, "true") {
		mqttConfiguration.cleanSession = true
		if strings.TrimSpace(store) != "" {
			mqttConfiguration.store = store
		}
	} else {
		// If value is omitted or is not true neither false, set by default its value to false
		mqttConfiguration.cleanSession = false
	}

	if strings.EqualFold(connectRety, "true") {
		mqttConfiguration.connectRety = true
	} else {
		// If value is omitted or is not true neither false, set by default its value to false
		mqttConfiguration.connectRety = false
	}

	mqttConfiguration.connectRetyInterval, err = strconv.Atoi(connectRetyInterval)
	if err != nil {
		log.Printf("Cannot convert connectionRetry value [%s] to its relative integer.", connectRetyInterval)
		mqttConfiguration.connectRetyInterval = 1 // As default value
	}

	if strings.EqualFold(TLSConnection, "true") {
		mqttConfiguration.TLSConnection = true
	} else {
		mqttConfiguration.TLSConnection = false
	}

	mqttConfiguration.clientAuthority = clientAuthority
	mqttConfiguration.clientCertificate = clientCertificate
	mqttConfiguration.clientKey = clientKey

	// Create new Client Options object
	clientOptions := mqtt.NewClientOptions()

	// Populate client object
	clientOptions.AddBroker(mqttConfiguration.protocol + mqttConfiguration.broker + ":" + mqttConfiguration.port)
	clientOptions.SetClientID(mqttConfiguration.clientId)
	clientOptions.SetAutoReconnect(mqttConfiguration.autoreconnect)
	clientOptions.SetCleanSession(mqttConfiguration.cleanSession)
	clientOptions.SetUsername(mqttConfiguration.user)
	clientOptions.SetPassword(mqttConfiguration.password)
	clientOptions.SetConnectRetry(mqttConfiguration.connectRety)
	clientOptions.SetConnectRetryInterval(time.Duration(mqttConfiguration.connectRetyInterval) * time.Millisecond)
	if mqttConfiguration.TLSConnection {
		clientOptions.SetTLSConfig(newTlsConfig())
	}

	// Set callbacks
	clientOptions.SetDefaultPublishHandler(messagePubHandler)
	clientOptions.OnConnect = connectHandler
	clientOptions.OnConnectionLost = connectLostHandler

	logMqttOpts(clientOptions)

	return mqtt.NewClient(clientOptions)
}

/*
	Define a private method to log MQTT configuration to end user
*/

func logMqttOpts(options *mqtt.ClientOptions) {
	log.Println("MQTT client configuration:")
	log.Printf("Client ID: %s", options.ClientID)
	log.Printf("Broker: %s", mqttConfiguration.protocol+mqttConfiguration.broker+":"+mqttConfiguration.port)
	log.Printf("QoS: %d", mqttConfiguration.qos)
	log.Printf("Retained messages: %t", mqttConfiguration.retained)
	log.Printf("Autoreconnect: %t", options.AutoReconnect)
	log.Printf("Clean Session: %t", options.CleanSession)
	log.Printf("Persistency path: %s", mqttConfiguration.store)
	log.Printf("Username: %s", options.Username)
	// Hide password while showing MQTT configuration to user
	log.Printf("Password: %s", strings.Repeat("*", utf8.RuneCountInString(options.Password)))
	log.Printf("Connect retry: %t", options.ConnectRetry)
	log.Printf("Connect retry interval: %v", options.ConnectRetryInterval)
	if mqttConfiguration.TLSConnection {
		log.Printf("Certificate Authority: %s", mqttConfiguration.clientAuthority)
		log.Printf("Client Certificate: %s", mqttConfiguration.clientCertificate)
		log.Printf("Client Key: %s", mqttConfiguration.clientKey)
	}
}

/*
	Define a private method to add TLS encryption during MQTT client setup
*/

func newTlsConfig() *tls.Config {
	certpool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(mqttConfiguration.clientAuthority)
	if err != nil {
		log.Fatalln(err.Error())
	}
	certpool.AppendCertsFromPEM(ca)
	// Import client certificate/key pair if set
	if strings.TrimSpace(mqttConfiguration.clientCertificate) != "" {
		clientKeyPair, err := tls.LoadX509KeyPair(mqttConfiguration.clientCertificate, mqttConfiguration.clientKey)
		if err != nil {
			panic(err)
		}
		return &tls.Config{
			RootCAs:            certpool,
			ClientAuth:         tls.NoClientCert,
			ClientCAs:          nil,
			InsecureSkipVerify: true,
			Certificates:       []tls.Certificate{clientKeyPair},
		}
	} else {
		// If client certificate is not set
		return &tls.Config{
			RootCAs: certpool,
		}
	}
}
