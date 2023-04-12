package telegram

import (
	"strings"
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
		// b, _ := json.MarshalIndent(update, " ", " ")
		// fmt.Println(string(b))
		
		var resp string
		if update.Message.IsCommand() {
			resp = bot.handleCmd(update.Message.Command(), update.Message.CommandArguments())
		}
		message := tgbotapi.NewMessage(userID, resp)
		message.ParseMode = tgbotapi.ModeMarkdown
		_, err := bot.Send(message)
		if err != nil {
			log.Printf("failed to send message back to user (%d): %v", userID, err)
		}
	}
}

func (t teleBot) handleCmd(cmd, args string) string {
	switch cmd {
	case "getdaily":
		var _args []string
		if args != "" {
			_args = strings.Split(args, " ")			
		}
		var opts *api.DailyOpts
		if len(_args) > 1 {
			return "Too many arguments for /getdaily. Should only contain 1: {COUNTRY}"
		}
		if len(_args) > 0 {
			country := _args[0]
			opts = &api.DailyOpts {
				Country: &country,
			}
		}
		return t.gClient.GetDailyTrends(opts)
	}
	return "Whoops! You've entered an invalid command."
}
