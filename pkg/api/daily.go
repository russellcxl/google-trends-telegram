package api

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/russellcxl/google-trends/pkg/types"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// https://trends.google.com/trends/api/dailytrends?hl=en-GB&tz=-480&geo=SG&hl=en-GB&ns=15

type GoogleClient struct {
	client        *http.Client
	defaultParams url.Values
}

func NewGoogleClient(defaultParams url.Values) *GoogleClient {
	return &GoogleClient{
		client:        http.DefaultClient,
		defaultParams: defaultParams,
	}
}

func (c *GoogleClient) GetDailyTrends() []string {
	path := "https://trends.google.com/trends/api/dailytrends"
	u, _ := url.Parse(path)
	u.RawQuery = c.defaultParams.Encode()
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
	searches := make([]types.TrendingSearch, 0)
	for _, v := range out.Default.Searches {
		searches = append(searches, v.Searches...)
	}
	fmt.Printf("%+v", searches[0])
	var output []string
	for _, s := range searches {
		output = append(output, s.Title.Query)
	}
	return output
}
