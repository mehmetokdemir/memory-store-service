package main

import (
	// Go imports
	"encoding/json"
	httpSwagger "github.com/swaggo/http-swagger"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/patrickmn/go-cache"
	// Internal imports
	_ "workout/memory-store-service/docs"
	"workout/memory-store-service/internal/handler"
	"workout/memory-store-service/model"
)

func RequestLogger(targetMux http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		targetMux.ServeHTTP(w, r)

		// log request by who(IP address)
		requesterIP := r.RemoteAddr

		log.Printf(
			"%s\t\t%s\t\t%s\t\t%v",
			r.Method,
			r.RequestURI,
			requesterIP,
			time.Since(start),
		)
	})
}

var srvHandler *handler.Handler

var task = func() {
	// todo : check here
	var result model.Result
	if result == nil {
		result = make(map[string]string)
	}
	cachedItems := srvHandler.Cache.Items()

	srvHandler.Lock()
	for key, val := range cachedItems {
		if valStr, ok := val.Object.(string); ok {
			result[key] = valStr
		}
	}
	srvHandler.Unlock()

	b, err := json.Marshal(result)
	if err != nil {
		log.Println("can not decode data", err.Error())
		return
	}
	if err := os.WriteFile(handler.File, b, os.ModePerm); err != nil {
		log.Println("can not write files", err.Error())
		return
	}
	log.Println("worker result", result)
}

// @title Key Value Store Restful API
// @version 1.0
// @description Key value store restful api
// @host localhost:8082
// @BasePath /
func main() {
	//
	// Worker
	//
	s := gocron.NewScheduler(time.UTC)
	_, err := s.Cron("*/1 * * * *").Do(task) // every minute
	if err != nil {
		log.Println("worker error", err.Error())
		return
	}
	s.StartAsync()

	//
	// Serve API
	//
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}
	http.HandleFunc("/docs/", httpSwagger.WrapHandler)
	http.HandleFunc("/memory", srvHandler.ServeHTTP)
	log.Fatalln(http.ListenAndServe(":"+port, nil))
}

func init() {
	srvHandler = handler.Service()
	if err := os.Chmod(handler.File, 0777); err != nil {
		log.Println("can not chmod file", err.Error())
		return
	}

	jsonFile, err := os.Open(handler.File)
	if err != nil {
		log.Println("can not open file", err.Error())
		return
	}

	defer jsonFile.Close()

	byt, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Println("can not read file", err.Error())
		return
	}
	var result model.Result
	if err = json.Unmarshal(byt, &result); err != nil {
		log.Println("can not decode data", err.Error())
		return
	}

	for k, v := range result {
		go srvHandler.Cache.Set(k, v, cache.NoExpiration)
	}
}
