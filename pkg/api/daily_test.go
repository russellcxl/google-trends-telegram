package api

import (
	"net/url"
	"testing"
)

func TestAPI(t *testing.T) {
	p := url.Values{}
	p.Set("geo", "SG")
	p.Set("hl", "en-GB")
	p.Set("tz", "480")
	p.Set("ns", "15")
	gClient := NewGoogleClient(p)
	gClient.GetDailyTrends()
}