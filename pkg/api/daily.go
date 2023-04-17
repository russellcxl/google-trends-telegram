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
