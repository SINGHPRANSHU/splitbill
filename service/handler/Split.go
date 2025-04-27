package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/singhpranshu/splitbill/common"
	"github.com/singhpranshu/splitbill/common/dto"
	jwt "github.com/singhpranshu/splitbill/service/middleware"
)

func (h *Handler) GetSplitData(w http.ResponseWriter, r *http.Request) {
	// Handler logic to get user
	groupId := chi.URLParam(r, "id")
	intId, err := strconv.Atoi(groupId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(common.GetHttpErrorResponse(http.StatusBadRequest, "invalid user_id")))
		return
	}
	if groupId == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(common.GetHttpErrorResponse(http.StatusBadRequest, "user_id is required")))
		return
	}
	splitData, err := h.DB.GetSplit(r.Context(), intId)
	if err != nil {
		log.Println("Error getting Expense:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(common.GetHttpErrorResponse(http.StatusNotFound, "something went wrong")))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(splitData)
}

func (h *Handler) AddExpense(w http.ResponseWriter, r *http.Request) {
	// Handler logic to create user
	var splitData *dto.SplitDto
	if err := json.NewDecoder(r.Body).Decode(&splitData); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(common.GetHttpErrorResponse(http.StatusBadRequest, "invalid request body")))
		return
	}
	if err := validate.Struct(splitData); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(common.GetHttpErrorResponse(http.StatusBadRequest, err.Error())))
		return
	}
	userClaim, ok := r.Context().Value(jwt.UserClaimKeyName).(jwt.UserClaims)
	if !ok {
		log.Println("Error getting user claim from context")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(common.GetHttpErrorResponse(http.StatusUnauthorized, "invalid token")))
		return
	}
	_, err := h.DB.GetUserById(r.Context(), userClaim.UserID)
	if err != nil {
		log.Println("Error getting user:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(common.GetHttpErrorResponse(http.StatusUnauthorized, err.Error())))
		return
	}

	member, err := h.DB.GetUsersByGroupIdAndUserId(r.Context(), splitData.GroupId, splitData.SplitWith)
	if err != nil {
		log.Println("Error getting user:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(common.GetHttpErrorResponse(http.StatusNotFound, err.Error())))
		return
	}

	log.Println("members:", len(member), "split with:", len(splitData.SplitWith))

	splitModel := dto.GetSplitModelFromDto(splitData, userClaim.UserID)

	if len(member) != len(splitData.SplitWith) || len(splitModel) != len(splitData.SplitWith) {
		log.Println("split with has wrong data:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(common.GetHttpErrorResponse(http.StatusNotFound, "split with has wrong data")))
		return
	}

	splitDataModel, err := h.DB.AddSplit(r.Context(), splitModel)

	if err != nil {
		log.Println("Error creating split:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(common.GetHttpErrorResponse(http.StatusInternalServerError, err.Error())))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(splitDataModel)
}
