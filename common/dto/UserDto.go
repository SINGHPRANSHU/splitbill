package dto

import "github.com/singhpranshu/splitbill/repository/model"

type UserDto struct {
	ID        int    `json:"id"`
	Username  string `json:"username" validate:"required,min=6,max=20"`
	Email     string `json:"email" validate:"required,email"`
	Phone     string `json:"phone" validate:"required"`
	CreatedAt *string `json:"created_at"`
	UpdatedAt *string `json:"updated_at"`
}


func GetUserDtoFromModel(user *model.User) *UserDto {
	return &UserDto{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}