package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/singhpranshu/splitbill/repository/model"
)

func (db *DB) GetGroup(ctx context.Context, id int) (*model.Group, error) {
	// Simulate a database call to get user
	var group *model.Group
	result := db.db.Where("id = ?", id).Find(&group)
	if result.Error != nil {
		return nil, result.Error
	}
	return group, nil
}

func (db *DB) CreateGroup(ctx context.Context, group *model.Group) (*model.Group, error) {
	// Simulate a database call to create user

	result := db.db.Omit("created_at", "updated_at", "id").Create(group)
	if result.Error != nil {
		fmt.Println("Error creating group:", result.Error)
		if errors.As(result.Error, &pgErr) {
			// Handle specific PostgreSQL errors
			if pgErr.Code == "23505" { // Unique constraint violation
				return nil, newSqlDuplicateEntryError("duplicate entry")
			}
		}
		return nil, result.Error
	}
	return group, nil
}

func (db *DB) GetUsersByGroupIdAndUserId(ctx context.Context, groupId int, userId []int) ([]model.UserGroupMap, error) {
	// Simulate a database call to get user
	var user []model.UserGroupMap
	result := db.db.Where("group_id = ? and member_id in ?", groupId, userId).Find(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}
