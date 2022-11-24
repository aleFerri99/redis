package functions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/interfaces"
)

const (
	maxRetry = 5
)

type SendCommand struct {
}

func NewSendCommand() *SendCommand {
	return &SendCommand{}
}

func (s *SendCommand) SendCommand(funcCtx interfaces.AppFunctionContext, data interface{}) (bool, interface{}) {
	lc := funcCtx.LoggingClient()

	lc.Debug("Sending Command")

	if data == nil {
		return false, errors.New("SendCommand: No data received")
	}

	gpio := data.(*GPIO)
	if verbose {
		prettyPrint, err := json.MarshalIndent(gpio, "", " ")
		if err != nil {
			log.Printf("Cannot pretty print GPIO data struct. Error: %s", err)
			return false, nil
		}
		log.Printf("GPIO:\n%s", string(prettyPrint))
	}

	// All the following values can be retreived by the device-service/cmd/device-service/res/profile/profile.yaml file
	action := os.Getenv("ACTION")              // PUT or GET -> set
	device := os.Getenv("DEVICE_NAME")         // DEVICE NAME -> "Modbus-Device"
	command := os.Getenv("COMMAND_NAME")       // COMMAND NAME -> "Modbus"
	resourceName := os.Getenv("RESOURCE_NAME") // RESOURCE NAME -> "GPIO"

	var response interface{}
	var err error

	switch action {
	case "set":
		lc.Infof("executing %s action", action)
		lc.Infof("Sending command '%s' for device '%s'", command, device)

		settings := make(map[string]string)
		settings[resourceName] = strconv.FormatBool(gpio.State)
		if verbose {
			prettyPrint, err := json.MarshalIndent(settings, "", " ")
			if err != nil {
				log.Printf("Cannot pretty print commands map. Error: %s", err)
				return false, nil
			}
			log.Printf("COMMANDS MAP:\n%s", string(prettyPrint))
		}

		// Add retry logic to handle start up case
		retry := true
		attempt := 0
		for retry {
			response, err = funcCtx.CommandClient().IssueSetCommandByName(context.Background(), device, command, settings)
			if err != nil && attempt <= maxRetry {
				log.Printf("failed to send '%s' set command to '%s' device: %s", command, device, err.Error())
				log.Printf("Retry attempt: %d", attempt)
				time.Sleep(3 * time.Second)
				attempt++
			} else if err != nil && attempt > maxRetry {
				return false, fmt.Errorf("failed to send '%s' set command to '%s' device: %s", command, device, err.Error())
			} else {
				log.Printf("Request worked at attempt %d", attempt)
				retry = false
			}
		}

	case "get":
		lc.Infof("executing %s action", action)
		lc.Infof("Sending command '%s' for device '%s'", command, device)
		response, err = funcCtx.CommandClient().IssueGetCommandByName(context.Background(), device, command, "no", "yes")
		if err != nil {
			return false, fmt.Errorf("failed to send '%s' get command to '%s' device: %s", command, device, err.Error())
		}

	default:
		lc.Errorf("Invalid action requested: %s", action)
		return false, nil
	}

	return true, response
}
