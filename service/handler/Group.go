package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/singhpranshu/splitbill/common"
	"github.com/singhpranshu/splitbill/common/dto"
	db "github.com/singhpranshu/splitbill/repository"
	"github.com/singhpranshu/splitbill/repository/model"
)

func (h *Handler) GetGroup(w http.ResponseWriter, r *http.Request) {
	// Handler logic to get user
	username := chi.URLParam(r, "id")
	intId, err := strconv.Atoi(username)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(common.GetHttpErrorResponse(http.StatusBadRequest, "invalid user_id")))
		return
	}
	if username == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(common.GetHttpErrorResponse(http.StatusBadRequest, "user_id is required")))
		return
	}
	group, err := h.DB.GetGroup(r.Context(), intId)
	if err != nil {
		log.Println("Error getting user:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(common.GetHttpErrorResponse(http.StatusNotFound, "something went wrong")))
		return
	}
	member, err := h.DB.GetmembersFromGroupId(r.Context(), group.ID)
	if err != nil {	
		log.Println("Error getting user:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(common.GetHttpErrorResponse(http.StatusNotFound, err.Error())))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.GetGroupDtoFromModel(group, *member))
}

func (h *Handler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	// Handler logic to create user
	var group *model.Group
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(common.GetHttpErrorResponse(http.StatusBadRequest, "invalid request body")))
		return
	}
	if err := validate.Struct(group); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(common.GetHttpErrorResponse(http.StatusBadRequest, err.Error())))
		return
	}
	user, err := h.DB.GetUserById(r.Context(), group.CreatedBy)
	if err != nil {
		log.Println("Error getting user:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(common.GetHttpErrorResponse(http.StatusUnauthorized, err.Error())))
		return
	}
	group.CreatedBy = user.ID
	group, err = h.DB.CreateGroup(r.Context(), group)
	if err != nil {
		log.Println("Error creating group:", err)
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
	_, err = h.DB.AddMember(r.Context(), &model.UserGroupMap{
		GroupId: group.ID,
		MemberId: user.ID,
	})
	if err != nil {
		log.Println("Error adding member:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(common.GetHttpErrorResponse(http.StatusInternalServerError, err.Error())))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.GetGroupDtoFromModel(group, []model.User{}))
}

func (h *Handler) Addmember(w http.ResponseWriter, r *http.Request) {
	var userGroupMap *model.UserGroupMap
	if err := json.NewDecoder(r.Body).Decode(&userGroupMap); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(common.GetHttpErrorResponse(http.StatusBadRequest, "invalid request body")))
		return
	}
	if err := validate.Struct(userGroupMap); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(common.GetHttpErrorResponse(http.StatusBadRequest, err.Error())))
		return
	}
	_, err := h.DB.GetUserById(r.Context(), userGroupMap.MemberId)
	if err != nil {
		log.Println("Error getting user:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(common.GetHttpErrorResponse(http.StatusUnauthorized, err.Error())))
		return
	}
	_, err = h.DB.GetGroup(r.Context(), userGroupMap.GroupId)
	if err != nil {
		log.Println("Error getting group:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(common.GetHttpErrorResponse(http.StatusUnauthorized, err.Error())))
		return
	}
	userGroupMap, err = h.DB.AddMember(r.Context(), userGroupMap)
	if err != nil {
		log.Println("Error adding member:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(common.GetHttpErrorResponse(http.StatusUnauthorized, err.Error())))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userGroupMap)

}
