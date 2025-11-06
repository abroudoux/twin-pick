package domain

type Film struct {
	Title           string   `json:"title"`
	Duration        int      `json:"duration"`
	Directors       []string `json:"directors"`
	DetailsEndpoint string   `json:"-"`
	Year            int      `json:"year"`
}

type Watchlist struct {
	Username string
	Films    []*Film
}

type ScrapperParams struct {
	Genres   []string
	Platform string
}

type PickParams struct {
	Usernames      []string
	ScrapperParams *ScrapperParams
	Limit          int
	Duration       Duration
}

type SpotParams struct {
	ScrapperParams *ScrapperParams
	Limit          int
}
