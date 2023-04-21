package api

import (
	"html"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/russellcxl/google-trends/pkg/utils"

	jsoniter "github.com/json-iterator/go"
	"github.com/russellcxl/google-trends/pkg/types"
)

// GetDailyTrends returns a list of the top trending topics
func (c *GoogleClient) GetDailyTrends(opts *types.DailyOpts) (text string, keyboard [][]string) {

	// validate opts
	if err := c.validateOpts(opts); err != nil {
		text = "Invalid parameters"
		return
	}

	// get daily trending topics from cache / url
	country := "SG"
	if opts != nil {
		if opts.Country != nil {
			country = *opts.Country
		}
	}
	out, err := c.getDaily(country)
	if err != nil {
		text = err.Error()
		return
	}

	// parse the results
	var searches []types.TrendingSearch
	for _, v := range out.Default.Searches {
		if !utils.IsToday(v.FormattedDate) {
			break
		}
		searches = append(searches, v.Searches...)
	}

	// get the number of topics to return
	var list string
	listCount := c.config.GoogleClient.Daily.ListCount
	var fewerThanExpected bool
	if listCount > len(searches) {
		listCount = len(searches)
		fewerThanExpected = true
	}

	// keyboard slice should contain [ [TEXT_SHOWN_TO_USER, CALLBACK_VALUE], ... ]
	for i := 0; i < listCount; i++ {
		s := searches[i]
		list += fmt.Sprintf("%s *%s*\n      _~%s searches_\n\n", utils.DigitUnicodesMap[i+1], s.Title.Query, s.FormattedTraffic)
		callback := types.GetDailyCallbackVal(country, i)
		keyboard = append(keyboard, []string{utils.DigitUnicodesMap[i+1], callback})
	}
	text = fmt.Sprintf("Top 7 trending topics in %s today:\n\n%s\nClick on any of the corresponding numbers below for more details!\n", country, list)
	if fewerThanExpected {
		text = fmt.Sprintf("%s\n\n_Oops! Looks like there are only %d topics right now._", text, listCount)
	}

	return
}


// GetDailyTrendsTopic returns a list of articles for a trending topic
func (c *GoogleClient) GetDailyTrendsTopic(country string, idx int) string {
	// get daily trending topics from cache / url
	out, err := c.getDaily(country)
	if err != nil {
		return err.Error()
	}
	topic := out.Default.Searches[0].Searches[idx]
	var list string
	for i, a := range topic.Articles {
		source := utils.RemoveTLD(a.Source)
		list += fmt.Sprintf("%s [%s](%s)\n -- %s\n\n", utils.DigitUnicodesMap[i+1], html.UnescapeString(a.Snippet), a.URL, source)
	}
	return fmt.Sprintf("*%s* (_%s searches_)\n\n%s", topic.Title.Query, topic.FormattedTraffic, list)
}

func (c *GoogleClient) getDaily(country string) (*types.Daily, error) {

	// set params according to opts
	params := cloneParams(c.params)
	params.Set("geo", country)
	redisKey := getRedisKey(country) // default resp is always for SG

	// check if data exists in Redis
	exists, err := c.redis.KeyExists(redisKey)
	if err != nil {
		return nil, err
	}

	var data string
	if exists != 0 {
		data, err = c.redis.GetValue(redisKey)
		if err != nil {
			return nil, fmt.Errorf("failed to get daily data from redis: %v", err)
		}
	} else {
		// else, call the API to retrieve data
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
		data = string(b)
		data = strings.Replace(data, ")]}',", "", 1)

		// and then store the data in Redis for 30 minutes
		// TODO: include redis lock
		expiry := time.Minute * time.Duration(c.config.GoogleClient.Daily.RefreshIntervalMinutes)
		if err = c.redis.SetValue(redisKey, data, expiry); err != nil {
			return nil, fmt.Errorf("failed to set daily data in redis: %v", err)
		}
	}
	out := new(types.Daily)
	jsoniter.UnmarshalFromString(data, out)
	return out, nil
}

func (c *GoogleClient) validateOpts(opts *types.DailyOpts) error {
	if opts != nil {
		if opts.Country != nil {
			code := *opts.Country
			if !c.validCountryCodes[code] {
				return fmt.Errorf("Invalid country code")
			}
		}
	}
	return nil
}

func cloneParams(params url.Values) url.Values {
	m := url.Values{}
	for k, v := range params {
		m[k] = v
	}
	return m
}

func getRedisKey(country string) string {
	return fmt.Sprintf("%s_%s", types.DailyRedisKey, country)
}
