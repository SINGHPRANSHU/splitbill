package dto

import "github.com/singhpranshu/splitbill/repository/model"

type GroupDto struct {
	ID        int       `json:"id"`
	Name      string    `json:"name" validate:"required,min=6,max=20"`
	CreatedBy int       `json:"created_by"`
	CreatedAt *string   `json:"created_at"`
	UpdatedAt *string   `json:"updated_at"`
	Members   []UserDto `json:"members"`
}

func GetGroupDtoFromModel(group *model.Group, users []model.User) *GroupDto {
	var userDto []UserDto = []UserDto{}
	for _, user := range users {
		userDto = append(userDto, *GetUserDtoFromModel(&user))
	}
	return &GroupDto{
		ID:        group.ID,
		Name:      group.Name,
		CreatedBy: group.CreatedBy,
		CreatedAt: group.CreatedAt,
		UpdatedAt: group.UpdatedAt,
		Members:   userDto,
	}
}
