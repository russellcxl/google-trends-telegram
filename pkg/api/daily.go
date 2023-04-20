package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/russellcxl/google-trends/pkg/utils"

	jsoniter "github.com/json-iterator/go"
	"github.com/russellcxl/google-trends/pkg/types"
)

func (c *GoogleClient) GetDailyTrends(opts *types.DailyOpts) string {
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
	var fewerThanExpected bool
	if listCount > len(searches) {
		listCount = len(searches)
		fewerThanExpected = true
	}
	for i := 0; i < listCount; i++ {
		s := searches[i]
		list += fmt.Sprintf("%s *%s*\n      _~%s searches_\n\n", intToDigitUnicode(i+1), s.Title.Query, s.FormattedTraffic)
	}
	output := fmt.Sprintf("Top 7 trending topics in %s today:\n\n%s\n", params.Get("geo"), list)
	if fewerThanExpected {
		output = fmt.Sprintf("%s\n\n_Oops! Looks like there are only %d topics right now._", output, listCount)
	}
	return output
}

func cloneParams(params url.Values) url.Values {
	m := url.Values{}
	for k, v := range params {
		m[k] = v
	}
	return m
}

func intToDigitUnicode(n int) string {
	return utils.DigitUnicodesMap[n]
}
