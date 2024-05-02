package cat_entity

type CatCreateDataResponse struct {
	Id        string `json:"id"`
	CreatedAt string `json:"createdAt"`
}
type CatCreateResponse struct {
	Message string                 `json:"message"`
	Data    *CatCreateDataResponse `json:"data"`
}

type CatSearchDataResponse struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Race        string   `json:"race"`
	Sex         string   `json:"sex"`
	AgeInMonth  int      `json:"ageInMonth"`
	ImageURLs   []string `json:"imageUrls"`
	Description string   `json:"description"`
	HasMatched  bool     `json:"hasMatched"`
	CreatedAt   string   `json:"createdAt"`
}

type CatSearchResponse struct {
	Messagge string                   `json:"message"`
	Data     *[]CatSearchDataResponse `json:"data"`
}
type CatDeleteDataResponse struct {
	Id string `json:"id"`
}

type CatDeleteResponse struct {
	Message string                 `json:"message"`
	Data    *CatDeleteDataResponse `json:"data"`
}
