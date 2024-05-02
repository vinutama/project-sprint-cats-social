package match_entity

type MatchCreateResponse struct {
	Message string `json:"message"`
}

type CatDataResponse struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Race        string   `json:"race"`
	Sex         string   `json:"sex"`
	AgeInMonth  int      `json:"ageInMonth"`
	ImageURLs   string `json:"imageUrls"`
	Description string   `json:"description"`
	HasMatched  bool     `json:"hasMatched"`
	CreatedAt   string   `json:"createdAt"`
}

type IssuerDataResponse struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
}

type MatchGetDataResponse struct {
	Id             string              `json:"id"`
	IssuedBy       *IssuerDataResponse `json:"issuedBy"`
	MatchCatDetail *CatDataResponse    `json:"matchCatDetail"`
	UserCatDetail  *CatDataResponse    `json:"userCatDetail"`
	Message        string              `json:"message"`
	CreatedAt      string              `json:"createdAt"`
}
type MatchGetResponse struct {
	Message string                  `json:"message"`
	Data     *[]MatchGetDataResponse `json:"data"`
}
