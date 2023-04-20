package telegram

import (
	"log"
	"os"

	"github.com/russellcxl/google-trends/pkg/types"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type teleBot struct {
	*tgbotapi.BotAPI
	gClient types.GoogleClient
	redis   types.RedisClient
}

func New(gClient types.GoogleClient, redisClient types.RedisClient) *teleBot {
	// init bot
	token, found := os.LookupEnv("TELEGRAM_TOKEN")
	if !found {
		log.Fatalf("failed to find telegram token in env")
	}
	b, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("failed to initialize telegram bot: %v", err)
	}
	return &teleBot{
		BotAPI:  b,
		gClient: gClient,
		redis:   redisClient,
	}

}

func (t *teleBot) Run(gClient types.GoogleClient) {

	// create channel to receive message from bot
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates, err := t.GetUpdatesChan(updateConfig)
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
		var resp string
		if update.Message.IsCommand() {
			resp, err = t.handleCmd(userID, update.Message.Command(), update.Message.CommandArguments())
			if err != nil {
				log.Println(err)
				resp = "Something went wrong with the bot :("
			}
		}
		message := tgbotapi.NewMessage(userID, resp)
		message.ParseMode = tgbotapi.ModeMarkdown
		if _, err = t.Send(message); err != nil {
			log.Printf("failed to send message back to user (%d): %v", userID, err)
		}
	}
}
