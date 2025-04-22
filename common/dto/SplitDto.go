package dto

import "github.com/singhpranshu/splitbill/repository/model"

type SplitDto struct {
	Amount         int    `json:"amount" validate:"required"`
	ExpenseAddedBy int    `json:"expense_added_by" validate:"required"`
	GroupId        int    `json:"group_id" validate:"required"`
	SplitWith      []int  `json:"split_with" validate:"required"`
	Type           string `json:"type" validate:"required,oneof=share paid"`
}

func GetSplitModelFromDto(splitDto *SplitDto) []model.Split {
	var splits []model.Split
	for _, userId := range splitDto.SplitWith {
		if splitDto.ExpenseAddedBy == userId {
			continue
		}
		split := model.Split{
			FromUser:  splitDto.ExpenseAddedBy,
			ToUser:    userId,
			Type:      splitDto.Type,
			Amount:    splitDto.Amount,
			GroupId:   splitDto.GroupId,
			CreatedAt: nil,
			UpdatedAt: nil,
		}
		splits = append(splits, split)
	}
	return splits
}
