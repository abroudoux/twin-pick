package domain

type WatchlistProvider interface {
	GetWatchlist(username string, scrapperFilters *ScrapperFilters) (*Watchlist, error)
}

type SuggestionsProvider interface {
	GetSuggestions(scrapperFilters *ScrapperFilters) ([]*Film, error)
}

type DetailsProvider interface {
	GetFilmDetails(film *Film) (*Film, error)
}
