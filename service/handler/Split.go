package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/singhpranshu/splitbill/common"
	"github.com/singhpranshu/splitbill/common/dto"
	jwt "github.com/singhpranshu/splitbill/service/middleware"
)

// @Summary get split data by group id
// @Description Retrieve split details
// @Produce json
// @Success 200
// @Router /split/{id} [get]
func (h *Handler) GetSplitData(w http.ResponseWriter, r *http.Request) {
	// Handler logic to get user
	groupId := chi.URLParam(r, "id")
	userClaim, _ := r.Context().Value(jwt.UserClaimKeyName).(*jwt.UserClaims)
	h.Logger.Println("userClaim:", userClaim)
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
		h.Logger.Println("Error getting Expense:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(common.GetHttpErrorResponse(http.StatusNotFound, "something went wrong")))
		return
	}

	//calculate the split
	userSplitData, err := h.DB.GetSplitByUser(r.Context(), intId)
	userBalances := make(map[int]int)
	if err == nil {
		for _, split := range userSplitData {
			userBalances[split.FromUser] -= split.Amount
			userBalances[split.ToUser] += split.Amount
		}
		h.Logger.Println("userSplitData", &userSplitData)
	}
	response := make(map[string]interface{})
	response["user_id"] = userClaim.UserID
	response["balances"] = userBalances[userClaim.UserID]

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(splitData)
}

// @Summary get final split data by group id
// @Description Retrieve split details
// @Produce json
// @Success 200
// @Router /split/{groupId} [get]
func (h *Handler) GetSplitDataByGroupId(w http.ResponseWriter, r *http.Request) {
	// Handler logic to get user
	groupId := chi.URLParam(r, "groupId")
	userClaim, _ := r.Context().Value(jwt.UserClaimKeyName).(*jwt.UserClaims)
	h.Logger.Println("userClaim:", userClaim)

	if groupId == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(common.GetHttpErrorResponse(http.StatusBadRequest, "user_id is required")))
		return
	}
	groupIdInt, err := strconv.Atoi(groupId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(common.GetHttpErrorResponse(http.StatusBadRequest, "invalid user_id")))
		return
	}
	//calculate the split
	userSplitData, err := h.DB.GetSplitByUser(r.Context(), groupIdInt)
	if err != nil {
		http.Error(w, "failed to get split data", http.StatusInternalServerError)
		return
	}
	netMap := make(map[string]int)

	for _, split := range userSplitData {
		// Sort pair so (A, B) and (B, A) always go to the same key
		var key string
		if split.FromUser < split.ToUser {
			key = fmt.Sprintf("%d-%d", split.FromUser, split.ToUser)
			netMap[key] += split.Amount
		} else {
			key = fmt.Sprintf("%d-%d", split.ToUser, split.FromUser)
			netMap[key] -= split.Amount
		}
	}

	// Step 3: Prepare response with merged net balances
	var enriched []dto.NetSplit
	for key, amt := range netMap {
		if amt == 0 {
			continue
		}

		// Split the key into two user IDs
		parts := strings.Split(key, "-")
		u1, _ := strconv.Atoi(parts[0])
		u2, _ := strconv.Atoi(parts[1])

		// Ensure the correct direction for the net amount
		if amt > 0 {
			enriched = append(enriched, dto.NetSplit{FromUser: u1, ToUser: u2, Amount: amt})
		} else {
			enriched = append(enriched, dto.NetSplit{FromUser: u2, ToUser: u1, Amount: -amt})
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(enriched)
}

// @Summary create split data by group id
// @Description Retrieve split details
// @Produce json
// @Success 200
// @Router /split [post]
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
	userClaim, ok := r.Context().Value(jwt.UserClaimKeyName).(*jwt.UserClaims)
	h.Logger.Println("userClaim:", userClaim, "ok:", ok, jwt.UserClaimKeyName, r.Context().Value(jwt.UserClaimKeyName))
	if !ok {
		h.Logger.Println("Error getting user claim from context")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(common.GetHttpErrorResponse(http.StatusUnauthorized, "invalid token")))
		return
	}
	_, err := h.DB.GetUserById(r.Context(), userClaim.UserID)
	if err != nil {
		h.Logger.Println("Error getting user:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(common.GetHttpErrorResponse(http.StatusUnauthorized, err.Error())))
		return
	}

	member, err := h.DB.GetUsersByGroupIdAndUserId(r.Context(), splitData.GroupId, splitData.SplitWith)
	if err != nil {
		h.Logger.Println("Error getting user:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(common.GetHttpErrorResponse(http.StatusNotFound, err.Error())))
		return
	}

	h.Logger.Println("members:", len(member), "split with:", len(splitData.SplitWith))

	splitModel := dto.GetSplitModelFromDto(splitData, userClaim.UserID)

	// if len(member) != len(splitData.SplitWith) || len(splitModel) != len(splitData.SplitWith) {
	// 	h.Logger.Println("split with has wrong data:", err)
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	w.Write([]byte(common.GetHttpErrorResponse(http.StatusNotFound, "split with has wrong data")))
	// 	return
	// }

	splitDataModel, err := h.DB.AddSplit(r.Context(), splitModel)

	if err != nil {
		h.Logger.Println("Error creating split:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(common.GetHttpErrorResponse(http.StatusInternalServerError, err.Error())))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(splitDataModel)
}
