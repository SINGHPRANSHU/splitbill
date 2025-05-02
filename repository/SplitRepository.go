package db

import (
	"context"
	"errors"
	"log"

	"github.com/singhpranshu/splitbill/repository/model"
	"gorm.io/gorm"
)

func (db *DB) GetSplit(ctx context.Context, id int) ([]model.Split, error) {
	// Simulate a database call to get user
	var split []model.Split
	result := db.db.Where("group_id = ?", id).Find(&split)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		log.Println("Error getting user:", result.Error)
		return nil, errors.New("something went wrong")
	}
	return split, nil
}

func (db *DB) AddSplit(ctx context.Context, split []model.Split) ([]model.Split, error) {
	// Simulate a database call to get user
	result := db.db.Omit("created_at", "updated_at", "id").Create(split)
	if result.Error != nil {
		log.Println("Error creating split:", result.Error)
		return nil, errors.New("something went wrong")
	}
	return split, nil
}

func (db *DB) GetSplitByUser(ctx context.Context, loggedInUserId string) ([]model.Split, error) {
	// Simulate a database call to get user
	var split []model.Split
	result := db.db.Where("from_user = ?", loggedInUserId).Find(&split).Group("from_user")
	if result.Error != nil {
		log.Println("Error creating split:", result.Error)
		return nil, errors.New("something went wrong")
	}
	return split, nil
}
