package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/russellcxl/google-trends/pkg/session"

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

	// terminate gracefully
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(0)
	}()

	// init clients
	rClient, err := session.New(os.Getenv("REDIS_ADDRESS"), os.Getenv("REDIS_PASSWORD"), 0)
	if err != nil {
		log.Fatalln(err)
	}
	gClient := api.NewGoogleClient(rClient)

	// allows program to start receiving messages from the telegram bot
	teleBot := telegram.New(gClient)
	teleBot.Run(gClient)
}
