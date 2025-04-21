package db

import (
	"context"
	"errors"
	"log"

	"github.com/singhpranshu/splitbill/repository/model"
	"gorm.io/gorm/clause"
)

func (db *DB) GetmembersFromGroupId(ctx context.Context, id int) (*[]model.User, error) {
	// Simulate a database call to get user
	var user *[]model.User = &[]model.User{}
	result := db.db.Where("user_group_maps.group_id = ?", id).Joins("Join user_group_maps on user_group_maps.member_id = users.id").Find(&user).Select("username", "email", "phone", "users.created_at", "users.updated_at", "users.id")
	if result.Error != nil {
		return nil, errors.New("something went wrong")
	}
	log.Println("Error getting user:", user)
	return user, nil
}

func (db *DB) AddMember(ctx context.Context, user *model.UserGroupMap) (*model.UserGroupMap, error) {
	// Simulate a database call to create user

	result := db.db.Clauses(clause.OnConflict{
        Columns:   []clause.Column{{Name: "group_id"}, {Name: "member_id"}}, // Columns to check for conflict
        DoNothing: true, // Do nothing if a conflict is detected
    }).Omit("created_at", "updated_at", "id").Where("not exist (select * from user_group_maps where group_id = ? and member_id = ? )", user.GroupId, user.MemberId).Create(user)
	if result.Error != nil {
		if errors.As(result.Error, &pgErr) {
			// Handle specific PostgreSQL errors
			if pgErr.Code == "23505" { // Unique constraint violation
				return nil, newSqlDuplicateEntryError("duplicate entry")
			}
		}
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("user already exists in group")
	}
	return user, nil
}
