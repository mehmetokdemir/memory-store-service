package handler

import (
	// Go imports
	"net/http"
	"sync"
	"time"

	// External imports
	"github.com/patrickmn/go-cache"
)

type Handler struct {
	sync.Mutex
	*http.ServeMux

	Cache *cache.Cache
}

func Service() *Handler {
	h := &Handler{
		ServeMux: http.NewServeMux(),
		Cache:    cache.New(5*time.Minute, 10*time.Minute),
	}
	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("content-type", "application/json")

	switch {
	case r.Method == http.MethodGet:
		h.Get(w, r)
	case r.Method == http.MethodPost:
		h.Set(w, r)
	case r.Method == http.MethodDelete:
		h.Flush(w, r)
	default:
		http.Error(w, "method not found", http.StatusNotFound)
	}
}
