package handler

import db "github.com/singhpranshu/splitbill/repository"

type Handler struct {
	DB *db.DB
}

func NewHandler(db *db.DB) *Handler {
	return &Handler{
		DB: db,
	}
}