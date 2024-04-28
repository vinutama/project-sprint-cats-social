package user_entity

type UserRegisterRequest struct {
	Name     string `json:"name" validate:"required,min=5,max=50"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required,min=5,max=15"`
}
