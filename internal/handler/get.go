package handler

import (
	// Go imports
	"encoding/json"
	"net/http"
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
	if r.URL.Query().Get("key") == "" || !r.URL.Query().Has("key") {
		if err := json.NewEncoder(w).Encode(GenerateResponse(http.StatusBadRequest, DescriptionEnumInvalidKeyError, nil)); err != nil {
			http.Error(w, DescriptionEnumServerError.String(), http.StatusInternalServerError)
		}
		return
	}

	value, ok := h.Cache.Get(r.URL.Query().Get("key"))
	if !ok {
		if err := json.NewEncoder(w).Encode(GenerateResponse(http.StatusNotFound, DescriptionEnumKeyNotFoundError, nil)); err != nil {
			http.Error(w, DescriptionEnumServerError.String(), http.StatusInternalServerError)
		}
		return
	}

	strValue, ok := value.(string)
	if !ok {
		if err := json.NewEncoder(w).Encode(GenerateResponse(http.StatusBadRequest, DescriptionEnumValueTypeError, nil)); err != nil {
			http.Error(w, DescriptionEnumServerError.String(), http.StatusInternalServerError)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(GenerateResponse(http.StatusOK, DescriptionEnumSuccess, Store{
		Key:   r.URL.Query().Get("key"),
		Value: strValue,
	})); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
