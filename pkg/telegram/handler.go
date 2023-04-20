package telegram

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/russellcxl/google-trends/pkg/types"
	"github.com/russellcxl/google-trends/pkg/utils"
)

func (t *teleBot) handleCmd(userID int64, cmd, args string) (string, error) {
	var _args []string
	if args != "" {
		_args = strings.Split(args, " ")
	}

	// check if user is authorized
	if cmd != "start" && !isUserAllowed(userID) {
		return "You're not authorized", nil
	}

	switch cmd {
	case "start":
		var input string
		if len(_args) > 0 {
			input = _args[0]
		}
		token, found := os.LookupEnv("ACCESS_TOKEN")
		if !found {
			return "", fmt.Errorf("failed to find access token in env")
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
				return "", fmt.Errorf("failed to write to allowed users: %v", err)
			}
			return "Welcome to the Google Trends bot!", nil
		}

	case "getdaily":
		var opts *types.DailyOpts
		if len(_args) > 1 {
			return "Too many arguments for /getdaily. Should only contain 1: {COUNTRY}", nil
		}
		if len(_args) > 0 {
			country := _args[0]
			opts = &types.DailyOpts{
				Country: &country,
			}
		}
		return t.gClient.GetDailyTrends(opts), nil
	}
	return "Whoops! You've entered an invalid command.", nil
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
