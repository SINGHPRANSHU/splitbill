package db

import (
	"context"
	"errors"
	"log"

	"github.com/singhpranshu/splitbill/repository/model"
	"gorm.io/gorm"
)

func (db *DB) GetUserById(ctx context.Context, id int) (*model.User, error) {
	// Simulate a database call to get user
	var user *model.User
	result := db.db.Where("id = ?", id).First(&user).Select("username", "email", "phone", "created_at", "updated_at", "id")
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		log.Println("Error getting user:", result.Error)
		return nil, errors.New("something went wrong")
	}
	return user, nil
}

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

	result := db.db.Omit("created_at", "updated_at", "id").Create(user)
	if result.Error != nil {
		if errors.As(result.Error, &pgErr) {
			// Handle specific PostgreSQL errors
			if pgErr.Code == "23505" { // Unique constraint violation
				return nil, newSqlDuplicateEntryError("duplicate entry")
			}
		}
		return nil, result.Error
	}
	return user, nil
}
