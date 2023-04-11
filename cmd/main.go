package main

import (
	"log"
	"net/url"

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
	
	config := utils.GetConfig("config.yaml")

	// initialize google trends client
	p := url.Values{}
	for _, val := range config.GoogleClient.DefaultParams {
		p.Set(val[0], val[1])
	}
	gClient := api.NewGoogleClient(p)

	// allows program to start receiving messages from the telegram bot
	telegram.Run(gClient)
}
