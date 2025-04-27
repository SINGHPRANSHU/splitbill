package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/singhpranshu/splitbill/common"
	"github.com/singhpranshu/splitbill/common/dto"
	db "github.com/singhpranshu/splitbill/repository"
	"github.com/singhpranshu/splitbill/repository/model"
)

var validate = validator.New()

// @Summary Get a user by name
// @Description Retrieve user details
// @Param user_name path string true "User Name"
// @Produce json
// @Success 200
// @Router /users/{user_name} [get]
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	// Handler logic to get user
	username := chi.URLParam(r, "user_name")
	if username == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(common.GetHttpErrorResponse(http.StatusBadRequest, "user_id is required")))
		return
	}
	user, err := h.DB.GetUser(r.Context(), username)
	if err != nil {
		log.Println("Error getting user:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(common.GetHttpErrorResponse(http.StatusNotFound, "something went wrong")))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.GetUserDtoFromModel(user))
}

// @Summary create a user by name
// @Description Retrieve user details
// @Produce json
// @Success 200
// @Router /users [post]
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Handler logic to create user
	var user *model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(common.GetHttpErrorResponse(http.StatusBadRequest, "invalid request body")))
		return
	}
	log.Println(user)
	if err := validate.Struct(user); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(common.GetHttpErrorResponse(http.StatusBadRequest, err.Error())))
		return
	}
	user, err := h.DB.CreateUser(r.Context(), user)
	if err != nil {
		log.Println("Error creating user:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		var duplicateEntryError *db.SqlDuplicateEntryError
		if errors.As(err, &duplicateEntryError) {
			w.Write([]byte(common.GetHttpErrorResponse(http.StatusConflict, err.Error())))
			return
		}
		w.Write([]byte(common.GetHttpErrorResponse(http.StatusInternalServerError, "something went wrong")))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.GetUserDtoFromModel(user))
}
