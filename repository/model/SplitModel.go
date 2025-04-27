package model

const (
	PaidSplitType  = "paid"
	ShareSplitType = "share"
)

type Split struct {
	ID        int     `json:"id"`
	FromUser  int     `json:"from_user" validate:"required"`
	ToUser    int     `json:"to_user" validate:"required"`
	Type      string  `json:"type" validate:"required,oneof=share paid"`
	Amount    int     `json:"amount" validate:"required"`
	GroupId   int     `json:"group_id" validate:"required"`
	CreatedAt *string `json:"created_at"`
	UpdatedAt *string `json:"updated_at"`
}
