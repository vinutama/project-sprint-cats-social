package user_entity

type UserRegisterResponse struct {
	Message string    `json:"message"`
	Data    *UserData `json:"data"`
}

type UserLoginResponse struct {
	Message string    `json:"message"`
	Data    *UserData `json:"data,omitempty"`
	Status  int       `json:"-"`
}

type UserData struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	AccessToken string `json:"accessToken"`
}
