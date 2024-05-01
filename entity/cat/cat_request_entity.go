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
