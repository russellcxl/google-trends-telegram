package types

type Daily struct {
	Default TrendingSearchesDays `json:"default"`
}

type TrendingSearchesDays struct {
	Searches []TrendingSearchDays `json:"trendingSearchesDays"`
}

type TrendingSearchDays struct {
	FormattedDate string            `json:"formattedDate"`
	Searches      []TrendingSearch `json:"trendingSearches"`
}

type TrendingSearch struct {
	Title            SearchTitle     `json:"title"`
	FormattedTraffic string           `json:"formattedTraffic"`
	Image            SearchImage     `json:"image"`
	Articles         []*SearchArticle `json:"articles"`
}

type SearchTitle struct {
	Query string `json:"query"`
}

type SearchImage struct {
	NewsURL  string `json:"newsUrl"`
	Source   string `json:"source"`
	ImageURL string `json:"imageUrl"`
}

type SearchArticle struct {
	Title   string       `json:"title"`
	TimeAgo string       `json:"timeAgo"`
	Source  string       `json:"source"`
	Image   *SearchImage `json:"image"`
	URL     string       `json:"url"`
	Snippet string       `json:"snippet"`
}
