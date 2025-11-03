package domain

type Film struct {
	Title string
}

type Watchlist struct {
	Username string
	Films    []Film
}

type ScrapperParams struct {
	Genres   []string
	Platform string
}

type PickParams struct {
	Usernames      []string
	ScrapperParams *ScrapperParams
	Limit          int
}

type SpotParams struct {
	ScrapperParams *ScrapperParams
	Limit          int
}
