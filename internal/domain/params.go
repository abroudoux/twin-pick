package domain

type Params struct {
	Filters         *Filters
	ScrapperFilters *ScrapperFilters
}

type PickParams struct {
	Usernames []string
	Params    *Params
}

type SpotParams struct {
	Params *Params
}

func NewParams(filters *Filters, scrapperFilters *ScrapperFilters) *Params {
	return &Params{
		Filters:         filters,
		ScrapperFilters: scrapperFilters,
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
