package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/russellcxl/google-trends/pkg/types"
)

func (t *teleBot) handleDaily(userID int64, cmd string, args []string) (string, *tgbotapi.InlineKeyboardMarkup, error) {
	var opts *types.DailyOpts
	if len(args) > 1 {
		return "Too many arguments for /getdaily. Should only contain 1: {COUNTRY}", nil, nil
	}
	if len(args) > 0 {
		country := args[0]
		opts = &types.DailyOpts{
			Country: &country,
		}
	}
	text, arr := t.gClient.GetDailyTrends(opts)
	var keyboard *tgbotapi.InlineKeyboardMarkup
	if len(arr) > 0 {
		var rows [][]tgbotapi.InlineKeyboardButton
		for i := 0; i < len(arr); i++ {
			desc := arr[i][0]
			callback := arr[i][1]
			if i%2 == 0 {
				rows = append(rows, []tgbotapi.InlineKeyboardButton{})
			}
			rows[len(rows)-1] = append(rows[len(rows)-1], tgbotapi.NewInlineKeyboardButtonData(desc, callback))
		}
		keyboard = &tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: rows,
		}
	}
	return text, keyboard, nil
}
