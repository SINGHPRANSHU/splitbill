package model

type Group struct {
	ID        int     `json:"id"`
	Name      string  `json:"name" validate:"required,min=6,max=20"`
	CreatedBy int     `json:"created_by" validate:"required"`
	CreatedAt *string `json:"created_at"`
	UpdatedAt *string `json:"updated_at"`
}
