package telegram

import (
	"strconv"
	"strings"

	"github.com/russellcxl/google-trends/pkg/types"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (t *teleBot) handleCallbackQuery(update tgbotapi.Update) {
	if update.CallbackQuery != nil {
		arr := strings.Split(update.CallbackQuery.Data, "_")
		prefix := arr[0]
		output := "I'm afraid don't know how to respond to that"

		switch prefix {
		case types.DailyCallbackPrefix:
			country := arr[1]
			i := arr[2]
			idx, _ := strconv.Atoi(i)
			output = t.gClient.GetDailyTrendsTopic(country, idx)
		}

		// send message to user
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, output)
		msg.ParseMode = tgbotapi.ModeMarkdown
		t.Send(msg)

		// answer the callback query to remove the "working" status from the button
		callbackConfig := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
		t.AnswerCallbackQuery(callbackConfig)

	}
}
