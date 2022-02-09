package handler

import (
	// Go imports
	"encoding/json"
	"io/ioutil"
	"net/http"

	// External imports
	"github.com/patrickmn/go-cache"

	// Internal imports
	. "workout/memory-store-service/model"
)

// Set godoc
// @Summary      Create New Store
// @Description  Set a new key with value
// @Tags         Memory
// @Produce      json
// @Param request body model.SetMemory true "Example Request"
// @Success      200  {object}  ApiResponse{data=model.Store} "Success"
// @Router       /memory [post]
func (h *Handler) Set(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		if err := json.NewEncoder(w).Encode(GenerateResponse(http.StatusBadRequest, DescriptionEnumBodyReadError, err.Error())); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var req SetMemory
	if err := json.Unmarshal(body, &req); err != nil {
		if err := json.NewEncoder(w).Encode(GenerateResponse(http.StatusBadRequest, DescriptionEnumBodyDecodeError, err.Error())); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if (req.Key == nil || *req.Key == "") || (req.Value == nil || *req.Value == "") {
		if err := json.NewEncoder(w).Encode(GenerateResponse(http.StatusBadRequest, DescriptionEnumBodyError, nil)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Not use lock or unlock because of "set" function make it.
	h.Cache.Set(*req.Key, *req.Value, cache.NoExpiration)

	if err := json.NewEncoder(w).Encode(GenerateResponse(http.StatusOK, DescriptionEnumSuccess, Store{
		Key:   *req.Key,
		Value: *req.Value,
	})); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
