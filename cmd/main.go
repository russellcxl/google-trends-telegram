package main

import (
	"log"

	"github.com/russellcxl/google-trends/pkg/utils"

	"github.com/russellcxl/google-trends/pkg/api"

	"github.com/joho/godotenv"
	"github.com/russellcxl/google-trends/pkg/telegram"
)

func main() {

	// load env variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("failed to load .env: %v", err)
	}

	// initialize google trends client with default params
	config := utils.GetConfig("config.yaml")
	gClient := api.NewGoogleClient(*config)

	// allows program to start receiving messages from the telegram bot
	telegram.Run(gClient)
}
