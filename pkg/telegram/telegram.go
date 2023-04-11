package telegram

import (
	"fmt"
	"log"
	"os"

	"github.com/russellcxl/google-trends/pkg/api"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var bot *teleBot

type teleBot struct {
	*tgbotapi.BotAPI
	gClient *api.GoogleClient
}

func Run(gClient *api.GoogleClient) {

	// init bot
	token, found := os.LookupEnv("TELEGRAM_TOKEN")
	if !found {
		log.Fatalf("failed to find telegram token in env")
	}
	b, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("failed to initialize telegram bot: %v", err)
	}
	bot = &teleBot{
		BotAPI:  b,
		gClient: gClient,
	}

	// create channel to receive message from bot
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		log.Fatalf("failed to set update channel for telegram bot: %v", err)
	}

	// process updates
	for update := range updates {

		// ignore non-message updates
		if update.Message == nil {
			continue
		}

		// send message back to user
		userID := update.Message.Chat.ID
		resp := bot.getResp(update.Message.Text)
		message := tgbotapi.NewMessage(userID, resp)
		_, err := bot.Send(message)
		if err != nil {
			log.Printf("failed to send message back to user (%d): %v", userID, err)
		}
	}
}

func (t teleBot) getResp(input string) string {
	switch input {
	case "/getdaily":
		topics := t.gClient.GetDailyTrends()
		var output string
		for _, topic := range topics {
			output += fmt.Sprintf("- %s\n", topic)
		}
		return output
	}
	return "Whoops! You've entered an invalid command."
}
