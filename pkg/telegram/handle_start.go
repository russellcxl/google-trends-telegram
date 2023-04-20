package telegram

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/russellcxl/google-trends/pkg/types"
	"github.com/russellcxl/google-trends/pkg/utils"
)

func (t *teleBot) handleStart(userID int64, cmd string, args []string) (string, *tgbotapi.InlineKeyboardMarkup, error) {
	var input string
	if len(args) > 0 {
		input = args[0]
	}
	token, found := os.LookupEnv("ACCESS_TOKEN")
	if !found {
		return "", nil, fmt.Errorf("failed to find access token in env")
	}
	inputBytes := md5.Sum([]byte(input))
	if token == hex.EncodeToString(inputBytes[:]) {
		path := filepath.Join(os.Getenv("DATA_PATH"), "allowed_users.json")
		ids := new(types.UserIDs)
		if err := utils.ReadJSONFile(path, ids); err != nil {
			ids.UserIDs = []int64{userID}
		} else {
			ids.UserIDs = append(ids.UserIDs, userID)
		}
		if err := utils.WriteJSONFile(path, ids); err != nil {
			return "", nil, fmt.Errorf("failed to write to allowed users: %v", err)
		}
		return "Welcome to the Google Trends bot!", nil, nil
	}
	return "You're not authorized!", nil, nil
}

func isUserAllowed(userID int64) bool {
	path := filepath.Join(os.Getenv("DATA_PATH"), "allowed_users.json")
	var ids types.UserIDs
	if err := utils.ReadJSONFile(path, &ids); err != nil {
		return false
	}
	for _, id := range ids.UserIDs {
		if userID == id {
			return true
		}
	}
	return false
}
