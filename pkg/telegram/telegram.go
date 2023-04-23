package telegram

import (
	"log"
	"os"
	"strings"

	"github.com/russellcxl/google-trends/pkg/types"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type teleBot struct {
	*tgbotapi.BotAPI
	gClient types.GoogleClient
}

func New(gClient types.GoogleClient) *teleBot {
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
	}
}

func (t *teleBot) Run() {

	// create channel to receive message from bot
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates, err := t.GetUpdatesChan(updateConfig)
	if err != nil {
		log.Fatalf("failed to set update channel for telegram bot: %v", err)
	}

	// process updates
	for update := range updates {

		// handle callbacks
		if update.CallbackQuery != nil {
			t.handleCallbackQuery(update)
		}

		// ignore non-message updates
		if update.Message == nil {
			continue
		}

		// send message back to user
		userID := update.Message.Chat.ID
		var resp string
		var keyboard *tgbotapi.InlineKeyboardMarkup
		if update.Message.IsCommand() {
			resp, keyboard, err = t.handleCmd(userID, update.Message.Command(), update.Message.CommandArguments())
			if err != nil {
				log.Println(err)
				resp = "Something went wrong with the bot :("
			}
		}
		message := tgbotapi.NewMessage(userID, resp)
		message.ParseMode = tgbotapi.ModeMarkdown
		if keyboard != nil {
			message.ReplyMarkup = *keyboard
		}
		if _, err = t.Send(message); err != nil {
			log.Printf("failed to send message back to user (%d): %v", userID, err)
		}

	}
}

// routes the cmd to the correct handler
func (t *teleBot) handleCmd(userID int64, cmd, args string) (string, *tgbotapi.InlineKeyboardMarkup, error) {
	var _args []string
	if args != "" {
		_args = strings.Split(args, " ")
	}

	// check if user is authorized
	if cmd != "start" && !isUserAllowed(userID) {
		return "You're not authorized", nil, nil
	}

	switch cmd {
	case "start":
		return t.handleStart(userID, cmd, _args)
	case "getdaily":
		return t.handleDaily(userID, cmd, _args)
	}
	return "Whoops! You've entered an invalid command.", nil, nil
}
