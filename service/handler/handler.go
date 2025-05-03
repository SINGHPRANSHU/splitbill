package handler

import (
	db "github.com/singhpranshu/splitbill/repository"
)

type Handler struct {
	DB     *db.DB
	Logger LoggerInterface
}

func NewHandler(db *db.DB, logger LoggerInterface) *Handler {
	return &Handler{
		DB:     db,
		Logger: logger,
	}
}
