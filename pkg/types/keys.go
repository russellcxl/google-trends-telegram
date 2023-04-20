package types

import (
	"fmt"
)

func GetDailyCallbackVal(country string, idx int) string {
	return fmt.Sprintf("%s_%s_%d", DailyCallbackPrefix, country, idx)
}