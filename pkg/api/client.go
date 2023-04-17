package api

import (
	"github.com/russellcxl/google-trends/pkg/utils"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/russellcxl/google-trends/config"

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
	path := filepath.Join(os.Getenv("DATA_PATH"), "country_codes.json")
	var data types.CountryCodes

	// if not found locally, get from URL and save it
	if err := utils.ReadJSONFile(path, &data); err != nil {
		
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
		if err := utils.WriteJSONFile(path, data); err != nil {
			log.Fatalf("failed to write country codes to json file: %v", err)
		}
	}

	// add to map and return
	for _, c := range data.Children {
		m[c.ID] = true
	}
	return m
}
