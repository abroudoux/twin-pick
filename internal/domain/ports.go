package domain

type WatchlistProvider interface {
	GetWatchlist(username string, scrapperParams *ScrapperParams) (*Watchlist, error)
}

type SuggestionsProvider interface {
	GetSuggestions(scrapperParams *ScrapperParams) ([]*Film, error)
}

type DetailsProvider interface {
	GetFilmDetails(film *Film) (*Film, error)
}
