package dto

type UserLoginDto struct {
	Username string `json:"username" validate:"required,min=6,max=20"`
	Password string `json:"password" validate:"required"`
}

type UserLoginResponseDto struct {
	AccessToken string `json:"access_token"`
	Username    string `json:"username"`
}
