package cat_entity

type CatCreateRequest struct {
	Name        string   `json:"name" validate:"required,min=1,max=30"`
	Race        string   `json:"race" validate:"required,catRace"`
	Sex         string   `json:"sex" validate:"required,sex"`
	AgeInMonth  int      `json:"ageInMonth" validate:"required,min=1,max=120082"`
	Description string   `json:"description" validate:"required,min=1,max=200"`
	ImageURLs   []string `json:"imageUrls" validate:"required,min=1,dive,required,url"`
}

type CatEditRequest struct {
	Name        string   `json:"name" validate:"min=1,max=30"`
	Race        string   `json:"race" validate:"catRace"`
	Sex         string   `json:"sex" validate:"sex"`
	AgeInMonth  int      `json:"ageInMonth" validate:"min=1,max=120082"`
	Description string   `json:"description" validate:"min=1,max=200"`
	ImageURLs   []string `json:"imageUrls" validate:"min=1,dive,required,url"`
}
type CatSearchQuery struct {
	Id         string `query:"id"`
	Race       string `query:"race" validate:"omitempty"`
	Sex        string `query:"sex" validate:"omitempty"`
	HasMatched string `query:"hasMatched" validate:"omitempty"`
	AgeInMonth string `query:"ageInMonth" validate:"omitempty,regex=^([<>]?)\\d+$"`
	Owned      string `query:"owned" validate:"omitempty"`
	Search     string `query:"search"`
	Limit      string `query:"limit" validate:"omitempty,number,min=0"`
	Offset     string `query:"offset" validate:"omitempty,number,min=0"`
}
