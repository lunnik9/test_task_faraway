package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"os/signal"
	"syscall"
	"test_task_faraway/config"
	"test_task_faraway/repository"
	"test_task_faraway/server"
	"test_task_faraway/services"
	"time"
)

const serverShutdownTimeout = 5

func main() {
	appConfig := config.LoadConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := log.Default()
	loader := services.NewLoaderStruct()

	challenger := services.NewChallengerStruct(loader, logger)

	cache := repository.NewCacheStruct(redis.NewClient(&redis.Options{
		Addr:     appConfig.RedisConfig.Address,
		Password: appConfig.RedisConfig.Password,
		DB:       appConfig.RedisConfig.DB,
	}))

	mongoConn, err := mongo.Connect(ctx, &options.ClientOptions{
		Hosts: []string{appConfig.MongoConfig.URI},
	})
	if err != nil {
		logger.Fatalf("cannot connect to mongo: %s", err)
	}

	statistic := repository.NewStatModuleStruct(mongoConn)

	s := server.NewServer(challenger, cache, statistic, repository.NewWordsOfWisdom(), logger)

	go func() {
		err = s.Listen(ctx, appConfig.ServerPort)
		if err != nil {
			logger.Fatalf(err.Error())
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	<-stopChan
	logger.Println("Shutting down server...")

	cancel()

	timeoutCtx, cancelTimeout := context.WithTimeout(context.Background(), serverShutdownTimeout*time.Second)
	defer cancelTimeout()

	logger.Println("Closed listener")

	<-timeoutCtx.Done()
}
