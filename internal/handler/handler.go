package handler

import (
	// Go imports
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	// External imports
	"github.com/patrickmn/go-cache"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	File = "tmp/TIMESTAMP-data.json"
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

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")

	w.Header().Set("content-type", "application/json")
	switch {
	case r.Method == http.MethodGet:
		h.Get(w, r)
		return
	case r.Method == http.MethodPost:
		h.Set(w, r)
		return
	case r.Method == http.MethodDelete:
		h.Flush(w, r)
		return
	case r.Method == http.MethodGet && strings.Contains(r.URL.Path, "swagger"):
		fmt.Println("girdi")
		h.HandleFunc("/swagger/", httpSwagger.WrapHandler)
		return
	}
}
