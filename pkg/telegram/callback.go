package telegram

import (
	"strconv"
	"github.com/russellcxl/google-trends/pkg/types"
	"strings"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (t *teleBot) handleCallbackQuery(update tgbotapi.Update) {
	if update.CallbackQuery != nil {
		fmt.Println(update.CallbackQuery.Data)
		arr := strings.Split(update.CallbackQuery.Data, "_")
		prefix := arr[0]
		output := "test"
		
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
