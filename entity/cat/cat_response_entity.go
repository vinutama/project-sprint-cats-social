package cat_entity

type CatDataPostResponse struct {
	Id        string `json:"id"`
	CreatedAt string `json:"createdAt"`
}
type CatPostResponse struct {
	Message string               `json:"message"`
	Data    *CatDataPostResponse `json:"data"`
}
