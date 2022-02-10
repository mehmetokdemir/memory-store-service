package handler

import (
	// Go imports
	"log"
	"net/http"
	"os"
	"time"
)

// LogRequestMiddleware middleware
func LogRequestMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		duration := time.Since(start)
		log.Printf("remote:%s method:%s url:%s duration:%s\n", r.RemoteAddr, r.Method, r.URL, duration)
	})
}

func OpenLogFile(logfile string) {
	if logfile != "" {
		lf, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)
		if err != nil {
			log.Fatal("OpenLogfile: os.OpenFile:", err)
		}
		log.SetOutput(lf)
	}
}
