package main

import (
	"context"
	"encoding/json"
	"fmt"
	"hubs/common"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/Azure/azure-amqp-common-go/v3/conn"
	"github.com/Azure/azure-amqp-common-go/v3/sas"
	eventhub "github.com/Azure/azure-event-hubs-go/v3"
	"github.com/Azure/azure-event-hubs-go/v3/eph"
	"github.com/Azure/azure-event-hubs-go/v3/storage"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/vutran/ansi/colors"
	"github.com/wonderivan/logger"
)

func main() {

	// Azure Storage account information
	storageAccountName := common.MustPassEvn("STORAGE_NAME")
	storageAccountKey := common.MustPassEvn("STORAGE_KEY")

	// Azure Storage container to store leases and checkpoints
	storageContainerName := "ephcontainer"

	// Azure Event Hub connection string
	eventHubConnStr := common.MustPassEvn("EVENT_HUBS_STRING")
	parsed, err := conn.ParsedConnectionFromStr(eventHubConnStr)
	if err != nil {
		// handle error
		logger.Error(err)
		os.Exit(1)
	}

	// create a new Azure Storage Leaser / Checkpointer
	cred, err := azblob.NewSharedKeyCredential(storageAccountName, storageAccountKey)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	leaserCheckpointer, err := storage.NewStorageLeaserCheckpointer(cred, storageAccountName, storageContainerName, azure.PublicCloud)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	// SAS token provider for Azure Event Hubs
	provider, err := sas.NewTokenProvider(sas.TokenProviderWithKey(parsed.KeyName, parsed.Key))
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// create a new EPH processor
	processor, err := eph.New(ctx, parsed.Namespace, parsed.HubName, provider, leaserCheckpointer, leaserCheckpointer)
	if err != nil {
		fmt.Println(err)
		return
	}

	// register a message handler -- many can be registered
	handlerID, err := processor.RegisterHandler(ctx,
		func(c context.Context, e *eventhub.Event) error {
			//fmt.Println(string(e.Data))
			var anomaly common.AnomalyEvent
			err = json.Unmarshal(e.Data, &anomaly)
			if err == nil {
				//msg := styles.Bold(colors.Red(anomaly.Value))
				fmt.Println(colors.Green("Device ID:"), anomaly.DeviceID)
				fmt.Println(colors.Green("Type:"), anomaly.EventType)
				fmt.Println(colors.Yellow("Property:"), anomaly.Property)
				//msg := fmt.Printf("%s %f", colors.Yellow("Value:"), anomaly.Value)
				strValue := strconv.FormatFloat(anomaly.Value, 'f', 5, 64)
				fmt.Println(colors.Yellow("Value:"), colors.Red(strValue))
			}
			return nil
		})
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	fmt.Printf("handler id: %q is running\n", handlerID)

	// unregister a handler to stop that handler from receiving events
	// processor.UnregisterHandler(ctx, handleID)

	// start handling messages from all of the partitions balancing across multiple consumers
	err = processor.StartNonBlocking(ctx)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	// Wait for a signal to quit:
	logger.Info("Listening for events. Press CTRL+C to exit.")
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	logger.Debug("Terminating program & closing EPH Processor")
	err = processor.Close(context.Background())
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	//os.Exit(0)

	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Interrupt, os.Kill)
	// <-signalChan

}
