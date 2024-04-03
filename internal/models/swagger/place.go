package swagger

type Filters struct {
	DistrictID     int    `json:"district_id,omitempty"`
	CityID         int    `json:"city_id,omitempty"`
	TagIDs         []int  `json:"tag_ids,omitempty"`
	PaginationPage int    `json:"pagination_page"`
	Name           string `json:"name,omitempty"`
	Variety        string `json:"variety"`
}
