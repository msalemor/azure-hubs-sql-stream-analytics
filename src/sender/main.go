package main

import (
	"context"
	"fmt"
	"hubs/common"
	"os"
	"strconv"
	"time"

	eventhub "github.com/Azure/azure-event-hubs-go/v3"
	"github.com/wonderivan/logger"
)

const (
	RaiseAnomaly    = 3
	AC_VENT         = 1
	GENERATOR_EVENT = 2
	MOTOR_EVENT     = 3
)

var (
	ConnectionStringFromEnv string
	verbose                 = false
	inBatch                 = false
	eventCount              = 10
	delayInMilliseconds     = 250
	timeoutSeconds          = 30
	hub                     *eventhub.Hub
	err                     error
	ctx                     context.Context
	cancel                  context.CancelFunc
)

func mustPassEvn(env string) string {
	v := os.Getenv(env)
	if v == "" {
		logger.Error("Must pass the enviromnet variable:", v)
		os.Exit(1)
	}
	return v
}

func getEnvironmentVariables() {
	if verbose {
		logger.Info("Getting environment variables")
	}
	ConnectionStringFromEnv = mustPassEvn("EVENT_HUBS_STRING")

	eventCountEnv := os.Getenv("EVENT_COUNT")
	if eventCountEnv != "" {
		eventCount, _ = strconv.Atoi(eventCountEnv)
	}
	delayInMillisecondsEnv := os.Getenv("EVENT_DELAY")
	if delayInMillisecondsEnv != "" {
		delayInMilliseconds, _ = strconv.Atoi(delayInMillisecondsEnv)
	}
	timeoutInSecondsEnv := os.Getenv("EVENT_TIMEOUT")
	if timeoutInSecondsEnv != "" {
		timeoutSeconds, _ = strconv.Atoi(timeoutInSecondsEnv)
	}
	verboseEnv := os.Getenv("VERBOSE")
	if verboseEnv != "" {
		verbose = true
	}
	batchEnv := os.Getenv("BATCH")
	if batchEnv != "" {
		inBatch = true
	}
}

func sendIndividual() {
	if verbose {
		logger.Debug(fmt.Sprintf("Sending %d messages every %d milliseconds", eventCount, delayInMilliseconds))
	}
	for i := 1; i <= eventCount; i++ {
		evt := common.GetRandomEvent()
		if verbose {
			logger.Debug(fmt.Sprintf("Event %d: %s", i, evt))
		}
		err = hub.Send(ctx, eventhub.NewEventFromString(evt))

		if err != nil {
			logger.Error("Unable to send message to Event Hubs.", err)
			os.Exit(1)
		}
		time.Sleep(time.Duration(delayInMilliseconds) * time.Millisecond)
	}
}

func sendInBatch() {
	if verbose {
		logger.Debug(fmt.Sprintf("Sending %d messages in batch", eventCount))
	}
	var events []*eventhub.Event
	for i := 1; i <= eventCount; i++ {
		evt := common.GetRandomEvent()
		events = append(events, eventhub.NewEventFromString(evt))
	}
	hub.SendBatch(ctx, eventhub.NewEventBatchIterator(events...))
	if verbose {
		logger.Debug(fmt.Sprintf("Sent: %d messages in batch", eventCount))
	}
}

func main() {

	if verbose {
		logger.Info("Connecting to Event hubs")
	}

	getEnvironmentVariables()

	// Connect to Event Hubs before a given context timeout

	hub, err = eventhub.NewHubFromConnectionString(ConnectionStringFromEnv)
	if err != nil {
		logger.Error("Unable to connect to Event Hubs:", err)
		os.Exit(1)
	}
	ctx, cancel = context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	// send a single message into a random partition
	if verbose {
		logger.Info(fmt.Sprintf("Sending: %d events with a delay of: %d milliseconds", eventCount, delayInMilliseconds))
	}

	if inBatch {
		sendInBatch()
	} else {
		sendIndividual()
	}

	// Wait for a signal to quit:
	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Interrupt, os.Kill)
	// <-signalChan

	if verbose {
		logger.Info("Closing connection to Event Hubs")
	}

	err = hub.Close(context.Background())
	if err != nil {
		logger.Error(err)
	}
}
