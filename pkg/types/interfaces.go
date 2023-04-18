package types

type GoogleClient interface {
	GetDailyTrends(*DailyOpts) string
}