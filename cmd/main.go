package main

import (
	"fmt"
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
	godotenv.Load(".env")

	// terminate gracefully
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(0)
	}()

	// init clients
	redisAddr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	redisClient, err := session.New(redisAddr, os.Getenv("REDIS_PASSWORD"), 0)
	if err != nil {
		log.Fatalln(err)
	}
	googleClient := api.NewGoogleClient(redisClient)

	// allows program to start receiving messages from the telegram bot
	teleBot := telegram.New(googleClient)
	teleBot.Run()
}
