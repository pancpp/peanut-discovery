package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pancpp/peanut-discovery/app"
	"github.com/pancpp/peanut-discovery/conf"
	"github.com/pancpp/peanut-discovery/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// config
	if err := conf.Init(); err != nil {
		log.Fatal(err)
	}

	// logger
	if err := logger.Init(); err != nil {
		log.Fatal(err)
	}

	// say hello
	log.Println("Hello, peanut discovery server!")

	// app
	if err := app.Init(ctx); err != nil {
		log.Fatal(err)
	}

	// wait for keyboard interruption
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	<-sigChan

	// say goodbye
	log.Println("Goodbye!")
}
