package main

import (
	"fmt"
	"github.com/victor-nach/time-tracker/config"
	"github.com/victor-nach/time-tracker/server"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	cfg := config.LoadSecrets()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to start logger: %v", err)
	}

	//mongoStore, _, err := mongo.New(cfg.DBURL, cfg.DBName)
	//if err != nil {
	//	log.Fatalf("failed to open mongodb: %v", err)
	//}

	srv := server.NewServer(nil, cfg, logger)

	// create channel to listen to shutdown signals
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		addr := fmt.Sprintf(":%s", cfg.Port)
		err := srv.Run(addr)
		if err != nil {
			log.Fatal(fmt.Sprintf("failed to start service: %v", err))
		}
	}()

	<-shutdownChan
	log.Println("Closing application")
	// do cleanups before exit
}
