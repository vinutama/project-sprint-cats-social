package cat_entity

type CatCreateDataResponse struct {
	Id        string `json:"id"`
	CreatedAt string `json:"createdAt"`
}
type CatCreateResponse struct {
	Message string                 `json:"message"`
	Data    *CatCreateDataResponse `json:"data"`
}
