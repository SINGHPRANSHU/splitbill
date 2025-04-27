package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/singhpranshu/splitbill/common"
	"github.com/singhpranshu/splitbill/common/dto"
	db "github.com/singhpranshu/splitbill/repository"
	"github.com/singhpranshu/splitbill/repository/model"
	jwt "github.com/singhpranshu/splitbill/service/middleware"
)

var validate = validator.New()

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
	hashedPassword, err := jwt.HashPassword(user.Password)
	if err != nil {
		log.Println("Error hashing password:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(common.GetHttpErrorResponse(http.StatusInternalServerError, "something went wrong")))
		return
	}
	user.Password = hashedPassword
	user, err = h.DB.CreateUser(r.Context(), user)
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

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	// Handler logic to login user
	var user *dto.UserLoginDto
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
	userModel, err := h.DB.GetUser(r.Context(), user.Username)
	if err != nil {
		log.Println("Invalid User", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(common.GetHttpErrorResponse(http.StatusUnauthorized, "something went wrong")))
		return
	}

	if !jwt.VerifyPassword(user.Password, userModel.Password) {
		log.Println("Invalid Password")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(common.GetHttpErrorResponse(http.StatusUnauthorized, "invalid password")))
		return
	}
	tokenString, err := jwt.CreateToken(*userModel)
	if err != nil {
		log.Println("Error generating token:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(common.GetHttpErrorResponse(http.StatusInternalServerError, "something went wrong")))
		return
	}
	response := dto.UserLoginResponseDto{
		AccessToken: tokenString,
		Username:    userModel.Username,
	}
	setTokenCookie(w, tokenString)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func setTokenCookie(w http.ResponseWriter, token string) {
	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
}
