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

type ProgramParams struct {
	Usernames      []string
	ScrapperParams *ScrapperParams
	Limit          int
}
