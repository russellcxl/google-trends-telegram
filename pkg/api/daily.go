package api

import (
	"log"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/russellcxl/google-trends/pkg/utils"

	jsoniter "github.com/json-iterator/go"
	"github.com/russellcxl/google-trends/pkg/types"
)

// https://trends.google.com/trends/api/dailytrends?hl=en-GB&tz=-480&geo=SG&hl=en-GB&ns=15

type GoogleClient struct {
	client            *http.Client
	params            url.Values
	config            utils.Config
	validCountryCodes map[string]bool
}

type DailyOpts struct {
	Country *string
}

func NewGoogleClient(config utils.Config) *GoogleClient {
	defaultParams := url.Values{}
	for _, val := range config.GoogleClient.DefaultParams {
		defaultParams.Set(val[0], val[1])
	}
	return &GoogleClient{
		client:            http.DefaultClient,
		params:            defaultParams,
		config:            config,
		validCountryCodes: getCountryCodes(),
	}
}

func getCountryCodes() map[string]bool {
	m := make(map[string]bool)
	url := "https://assets.api-cdn.com/serpwow/serpwow_google_trends_geos.json"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("failed to get Google country codes: %v", err)
	}
	defer resp.Body.Close()
	var data types.CountryCodes
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Fatalf("failed to decode Google country codes: %v", err)
	}
	for _, c := range data.Children {
		m[c.ID] = true
	}
	return m
}

func (c *GoogleClient) GetDailyTrends(opts *DailyOpts) string {
	params := cloneParams(c.params)
	if opts != nil {
		if opts.Country != nil {
			code := *opts.Country
			if !c.validCountryCodes[code] {
				return "Invalid country code. Try something like SG or MY"
			}
			params.Set("geo", code)
		}
	}
	path := types.BaseURL + types.DailyTrendsURLPrefix
	u, _ := url.Parse(path)
	u.RawQuery = params.Encode()
	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, u.String(), nil)
	if err != nil {
		panic(err)
	}
	r.Header.Add("Accept", "application/json")
	resp, err := c.client.Do(r)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	data := string(b)
	data = strings.Replace(data, ")]}',", "", 1)
	out := new(types.Daily)
	jsoniter.UnmarshalFromString(data, out)

	var searches []types.TrendingSearch
	for _, v := range out.Default.Searches {
		searches = append(searches, v.Searches...)
	}

	var list string
	listCount := c.config.GoogleClient.Daily.ListCount
	for i := 0; i < listCount; i++ {
		s := searches[i]
		list += fmt.Sprintf("*%s*\n      _~%s searches_\n\n", s.Title.Query, s.FormattedTraffic)
	}
	output := fmt.Sprintf("Top %d trending topics in %s today:\n\n%s", listCount, params.Get("geo"), list)
	return output
}

func cloneParams(params url.Values) url.Values {
	m := url.Values{}
	for k, v := range params {
		m[k] = v
	}
	return m
}