package db

import (
	"context"

	"github.com/singhpranshu/splitbill/repository/model"
)

func (db *DB) GetUser(ctx context.Context, username string) (*model.User, error) {
	// Simulate a database call to get user
	var user *model.User
	result := db.db.Where("username = ?", username).Find(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (db *DB) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	// Simulate a database call to create user

	result := db.db.Omit("created_at", "updated_at").Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}
