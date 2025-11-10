package domain

type Filters struct {
	Limit    int
	Duration Duration
}

type OrderFilter string

const (
	OrderFilterPopular  OrderFilter = "popular"
	OrderFilterHighest  OrderFilter = "highest-rated"
	OrderFilterNewest   OrderFilter = "newest"
	OrderFilterShortest OrderFilter = "shortest"
)

type ScrapperFilters struct {
	Genres   []string
	Platform string
	Order    OrderFilter
}

func NewScrapperFilters(genres []string, platform string, order OrderFilter) *ScrapperFilters {
	return &ScrapperFilters{
		Genres:   genres,
		Platform: platform,
		Order:    order,
	}
}

func NewFilters(limit int, duration Duration) *Filters {
	return &Filters{
		Limit:    limit,
		Duration: duration,
	}
}
