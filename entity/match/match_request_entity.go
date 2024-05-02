package match_entity

type MatchCreateRequest struct {
	MatchCatId string `json:"matchCatId" validate:"required"`
	UserCatId  string `json:"userCatId" validate:"required"`
	Message    string `json:"message" validate:"required,min=5,max=120"`
}

type MatchApproveRequest struct {
	MatchId string `json:"matchCatId" validate:"required"`
}

type MatchDeleteParams struct {
	Id string `param:"id" validate:"required"`
}
