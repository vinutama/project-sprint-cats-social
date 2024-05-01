package match_entity

type MatchCreateDataResponse struct {
	Id        string `json:"id"`
	CreatedAt string `json:"createdAt"`
}

type MatchCreateResponse struct {
	Message string                   `json:"message"`
	Data    *MatchCreateDataResponse `json:"data"`
}

type MatchApproveResponse struct {
	Message string `json:"message"`
}
