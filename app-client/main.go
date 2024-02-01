package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	
	"test_task_faraway/app-client/client"
	"test_task_faraway/app-client/config"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	appConfig := config.LoadConfig()

	logger := log.Default()
	c := client.NewClient(logger, appConfig.Host)

	go c.Run(ctx)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	<-stopChan
	logger.Println("Shutting down server...")

	cancel()
}
