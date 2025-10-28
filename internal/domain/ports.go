package domain

type WatchlistProvider interface {
	GetWatchlist(username string, scrapperParams ScrapperParams) (Watchlist, error)
}
