package api

import (
	"path/filepath"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/russellcxl/google-trends/pkg/utils"

	"github.com/russellcxl/google-trends/config"

	jsoniter "github.com/json-iterator/go"
	"github.com/russellcxl/google-trends/pkg/types"
)

// https://trends.google.com/trends/api/dailytrends?hl=en-GB&tz=-480&geo=SG&hl=en-GB&ns=15

type GoogleClient struct {
	client            *http.Client
	params            url.Values
	config            *config.Config
	validCountryCodes map[string]bool
}

type DailyOpts struct {
	Country *string
}

func NewGoogleClient() *GoogleClient {
	cfg := config.GetConfig()
	defaultParams := url.Values{}
	for _, val := range cfg.GoogleClient.DefaultParams {
		defaultParams.Set(val[0], val[1])
	}
	return &GoogleClient{
		client:            http.DefaultClient,
		params:            defaultParams,
		config:            cfg,
		validCountryCodes: getCountryCodes(),
	}
}

func getCountryCodes() map[string]bool {
	m := make(map[string]bool)
	fileName := filepath.Join(os.Getenv("DATA_PATH"), "country_codes.json")
	b, err := os.ReadFile(fileName)
	
	// if not found locally, get from URL and save it
	var data types.CountryCodes
	if err != nil {

		// get from URL
		log.Println("Country codes not found. Getting from URL")
		url := "https://assets.api-cdn.com/serpwow/serpwow_google_trends_geos.json"
		resp, err := http.Get(url)
		if err != nil {
			log.Fatalf("failed to get Google country codes: %v", err)
		}
		defer resp.Body.Close()
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			log.Fatalf("failed to decode Google country codes: %v", err)
		}

		// write to file
		output, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			log.Fatalf("failed to marshal country codes: %v", err)
		}
		if err = os.WriteFile("country_codes.json", output, 0644); err != nil {
			log.Fatalf("failed to write country codes to json file: %v", err)
		}
	} else {
		if err := json.Unmarshal(b, &data); err != nil {
			log.Fatalf("failed to read country codes from json: %v", err)
		}
	}

	// add to map and return
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
		if !utils.IsToday(v.FormattedDate) {
			break
		}
		searches = append(searches, v.Searches...)
	}

	var list string
	listCount := c.config.GoogleClient.Daily.ListCount
	if listCount > len(searches) {
		listCount = len(searches)
	}
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
