package domain

type Film struct {
	Name string
}

type Watchlist struct {
	Username string
	Films    []Film
}

type ScrapperParams struct {
	Genres   []string
	Platform string
}
