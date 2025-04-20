package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/singhpranshu/splitbill/common"
)

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	// Handler logic to get user
	userid := chi.URLParam(r, "user_id")
	if userid == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(common.GetHttpErrorResponse(http.StatusBadRequest, "user_id is required")))
		return
	}
	user, err := h.DB.GetUser(r.Context(), userid)
	if err != nil{
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(common.GetHttpErrorResponse(http.StatusNotFound, "something went wrong")))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Handler logic to create user
	w.Write([]byte("CreateUser"))
}
