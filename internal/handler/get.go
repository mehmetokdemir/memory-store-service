package handler

import (
	// Go imports
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	// Internal imports
	. "workout/memory-store-service/model"
)

// Get godoc
// @Summary      Read Value
// @Description  Read the value of the key
// @Tags         Memory
// @Produce      json
// @Param key query string true "Key"
// @Success      200  {object}  ApiResponse{data=model.Store} "Success"
// @Router       /memory [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	fmt.Println(strings.Split(r.URL.String(), "/"))
	if !r.URL.Query().Has("key") {
		http.Error(w, fmt.Sprintf("key %s not found", r.URL.Query().Get("key")), http.StatusNotFound)
		return
	}

	if r.URL.Query().Get("key") == "" {
		http.Error(w, fmt.Sprintf("key %s not found", r.URL.Query().Get("key")), http.StatusNotFound)
		return
	}

	foo, ok := h.Cache.Get(r.URL.Query().Get("key"))
	if !ok {
		http.Error(w, "sdsdd", http.StatusNotFound)
		return
	}

	strVal, ok := foo.(string)
	if !ok {
		http.Error(w, "val type is not a string", http.StatusForbidden)
		return
	}

	if err := json.NewEncoder(w).Encode(GenerateResponse(http.StatusOK, DescriptionEnumSuccess, Store{
		Key:   r.URL.Query().Get("key"),
		Value: strVal,
	})); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
