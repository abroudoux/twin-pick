package domain

type Filters struct {
	Limit    int
	Duration Duration
}

type ScrapperParams struct {
	Genres   []string
	Platform string
}

type Params struct {
	Filters        *Filters
	ScrapperParams *ScrapperParams
}

type PickParams struct {
	Usernames []string
	Params    *Params
}

type SpotParams struct {
	Params *Params
}

func NewScrapperParams(genres []string, platform string) *ScrapperParams {
	return &ScrapperParams{
		Genres:   genres,
		Platform: platform,
	}
}

func NewFilters(limit int, duration Duration) *Filters {
	return &Filters{
		Limit:    limit,
		Duration: duration,
	}
}

func NewParams(filters *Filters, scrapperParams *ScrapperParams) *Params {
	return &Params{
		Filters:        filters,
		ScrapperParams: scrapperParams,
	}
}

func NewPickParams(usernames []string, params *Params) *PickParams {
	return &PickParams{
		Usernames: usernames,
		Params:    params,
	}
}

func NewSpotParams(params *Params) *SpotParams {
	return &SpotParams{
		Params: params,
	}
}
