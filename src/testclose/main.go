package main

import (
	"fmt"
	"hubs/common"
	"os"
	"os/signal"
	"time"

	"github.com/wonderivan/logger"
)

var Version = "development"

func main() {
	logger.Info(Version)
	logger.Info(common.User)
	logger.Info(common.Time)
	killSignal := make(chan os.Signal, 1)
	signal.Notify(killSignal, os.Interrupt)
	go func() {
		for {
			logger.Info("Doing Work")
			time.Sleep(1 * time.Second)
		}
	}()
	<-killSignal
	fmt.Println("Thanks for using Golang!")
}
