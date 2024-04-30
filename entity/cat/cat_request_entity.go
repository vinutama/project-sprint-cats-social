package cat_entity

type CatCreateRequest struct {
	Name        string   `json:"name" validate:"required,min=1,max=30"`
	Race        string   `json:"race" validate:"required,oneof=Persian MaineCoon Siamese Ragdoll Bengal Sphynx BritishShorthair Abyssinian ScottishFold Birman"`
	Sex         string   `json:"sex" validate:"required,oneof=male female"`
	AgeInMonth  int      `json:"ageInMonth" validate:"required,min=1,max=120082"`
	Description string   `json:"description" validate:"required,min=1,max=200"`
	ImageURLs   []string `json:"imageUrls" validate:"required,min=1,dive,required,url"`
}
