package db

import (
	"context"

	"github.com/singhpranshu/splitbill/repository/model"
)

func (db *DB) GetUser(ctx context.Context, userID string) (*model.User, error) {
	// Simulate a database call to get user
	return &model.User{}, nil
}