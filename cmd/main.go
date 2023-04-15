package main

import (
	"syscall"
	"os/signal"
	"os"
	"log"

	"github.com/russellcxl/google-trends/pkg/api"

	"github.com/joho/godotenv"
	"github.com/russellcxl/google-trends/pkg/telegram"
)

func main() {

	// load env variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("failed to load local .env file: %v", err)
	}

	// initialize google trends client with default params
	gClient := api.NewGoogleClient()

	// terminate gracefully
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(0)
	}()

	// allows program to start receiving messages from the telegram bot
	telegram.Run(gClient)
}