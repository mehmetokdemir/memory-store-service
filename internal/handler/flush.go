package handler

import (
	// Go imports
	"encoding/json"
	"net/http"
	"os"
	"sync"
	"workout/memory-store-service/constant"

	// Internal imports
	. "workout/memory-store-service/model"
)

// Flush godoc
// @Summary      Flush
// @Description  Delete all store values
// @Tags         Memory
// @Produce      json
// @Success      200  {object}  ApiResponse{data=model.Store} "Success"
// @Router       /memory [delete]
func (h *Handler) Flush(w http.ResponseWriter, r *http.Request) {
	wgDone := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		h.Cache.Flush()
	}()
	go func() {
		defer wg.Done()
		_ = os.RemoveAll(constant.TmpDataFile)
	}()

	go func() {
		wg.Wait()
		close(wgDone)
	}()

	select {
	case <-wgDone:
		break
	}

	if _, err := os.Stat(File); err != nil || h.Cache.Items() == nil {
		if err := json.NewEncoder(w).Encode(GenerateResponse(http.StatusOK, DescriptionEnumSuccess, nil)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
