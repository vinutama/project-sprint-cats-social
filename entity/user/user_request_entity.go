package user_entity

// TODO: validate not work, need to install https://github.com/go-playground/validator
type UserRegisterRequest struct {
	Name     string `json:"name" validate:"required,min=5,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=5,max=15"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=5,max=15"`
}
