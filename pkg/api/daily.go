package api

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"github.com/russellcxl/google-trends/pkg/types"

	jsoniter "github.com/json-iterator/go"
)

// https://trends.google.com/trends/api/dailytrends?hl=en-GB&tz=-480&geo=SG&hl=en-GB&ns=15

func Run() {
	path := "https://trends.google.com/trends/api/dailytrends"
	u, _ := url.Parse(path)
	p := url.Values{}
	p.Set("geo", "SG")
	p.Set("hl", "en-GB")
	p.Set("tz", "480")
	p.Set("ns", "15")
	u.RawQuery = p.Encode()
	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, u.String(), nil)
	if err != nil {
		panic(err)
	}
	r.Header.Add("Accept", "application/json")
	client := http.DefaultClient
	resp, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
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
	for _, s := range searches {
		fmt.Println(s.Title.Query)
	}
}
