package model

type UserGroupMap struct {
	ID        int     `json:"id"`
	GroupId   int     `json:"group_id" validate:"required"`
	MemberId  int     `json:"member_id" validate:"required"`
	CreatedAt *string `json:"created_at"`
	UpdatedAt *string `json:"updated_at"`
}
