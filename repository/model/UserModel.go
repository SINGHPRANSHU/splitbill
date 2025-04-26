package model

type User struct {
	ID        int     `json:"id"`
	Username  string  `json:"username" validate:"required,min=6,max=20"`
	Email     string  `json:"email" validate:"required,email"`
	Password  string  `json:"password,omitempty" validate:"required,min=8,max=20"`
	Phone     string  `json:"phone" validate:"required"`
	CreatedAt *string `json:"created_at"`
	UpdatedAt *string `json:"updated_at"`
}	
