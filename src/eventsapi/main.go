package main

import (
	"context"
	"encoding/json"
	"fmt"
	"hubs/common"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	eventhub "github.com/Azure/azure-event-hubs-go/v3"
	"github.com/gorilla/mux"
	"github.com/wonderivan/logger"
)

var (
	err              error
	wg               sync.WaitGroup
	connectionString string
	verbose          bool
)

func loggingMiddleware(w http.ResponseWriter, r *http.Request) {
	logger.Info(r.Method, r.RequestURI)
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	loggingMiddleware(w, r)
}

func processIndividually(eventRequest common.EventsRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(30)*time.Second)
	defer cancel()
	hub, err := eventhub.NewHubFromConnectionString(connectionString)
	if err != nil {
		logger.Error(err)
		return err
	}
	if verbose {
		logger.Debug(fmt.Sprintf("Sending %d messages every %d milliseconds", eventRequest.Count, eventRequest.Delay))
	}
	for i := 1; i <= eventRequest.Count; i++ {
		evt := common.GetRandomEvent()
		if verbose {
			logger.Debug(fmt.Sprintf("Event %d: %s", i, evt))
		}
		err = hub.Send(ctx, eventhub.NewEventFromString(evt))

		if err != nil {
			logger.Error("Unable to send message to Event Hubs.", err)
			return err
		}
		time.Sleep(time.Duration(eventRequest.Delay) * time.Millisecond)
	}
	wg.Done()
	return nil
}

func processInBatch(eventRequest common.EventsRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(30)*time.Second)
	defer cancel()

	hub, err := eventhub.NewHubFromConnectionString(connectionString)
	if err != nil {
		logger.Error(err)
		return err
	}
	if verbose {
		logger.Debug(fmt.Sprintf("Sending %d messages in batch", eventRequest.Count))
	}
	var events []*eventhub.Event
	for i := 1; i <= eventRequest.Count; i++ {
		evt := common.GetRandomEvent()
		events = append(events, eventhub.NewEventFromString(evt))
		time.Sleep(1 * time.Millisecond)
	}
	hub.SendBatch(ctx, eventhub.NewEventBatchIterator(events...))
	if verbose {
		logger.Debug(fmt.Sprintf("Sent: %d messages in batch", eventRequest.Count))
	}
	wg.Done()
	return nil
}

func postEvents(w http.ResponseWriter, r *http.Request) {
	loggingMiddleware(w, r)

	eventsRequest := common.EventsRequest{}

	if err = json.NewDecoder(r.Body).Decode(&eventsRequest); err != nil {
		logger.Error(err)
		http.Error(w, "Error decoidng response object", http.StatusBadRequest)
		return
	}

	if eventsRequest.Batch {
		wg.Add(1)
		go processInBatch(eventsRequest)
		wg.Wait()
	} else {
		wg.Add(1)
		go processIndividually(eventsRequest)
		wg.Wait()
	}

	response, err := json.Marshal(common.EventsResponse{
		Message: "Messages written to event hubs",
		Count:   eventsRequest.Count,
		Delay:   eventsRequest.Delay,
		Batch:   eventsRequest.Batch,
	})
	if err != nil {
		logger.Error(err)
		http.Error(w, "Error encoding response object", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func mustPassEvn(env string) string {
	v := os.Getenv(env)
	if v == "" {
		logger.Error("Must pass the enviromnet variable:", v)
		os.Exit(1)
	}
	return v
}

func main() {
	r := mux.NewRouter()
	connectionString = mustPassEvn("EVENT_HUBS_STRING")
	verboseEnv := os.Getenv("VERBOSE")
	if (verboseEnv) != "" {
		verbose = true
	}

	r.HandleFunc("/", infoHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/events", postEvents).Methods(http.MethodPost)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Info("Starting server at: http://localhost:8000")
	log.Fatal(srv.ListenAndServe())
}
