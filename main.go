package main

import (
	// Go imports
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	// External imports
	"github.com/go-co-op/gocron"
	"github.com/patrickmn/go-cache"
	httpSwagger "github.com/swaggo/http-swagger"

	// Internal imports
	"workout/memory-store-service/constant"
	_ "workout/memory-store-service/docs"
	"workout/memory-store-service/internal/handler"
	"workout/memory-store-service/model"
)

var srvHandler *handler.Handler

var task = func() {
	result := make(model.Result)
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
		fmt.Println("can not decode data", err.Error())
		return
	}

	if err := os.WriteFile(constant.TmpDataFile, b, os.ModePerm); err != nil {
		fmt.Println("can not write files", err.Error())
		return
	}
}

// @title Key Value Store Restful API
// @version 1.0
// @description Key value store restful api
// @BasePath /
func main() {
	//
	// Worker
	//
	s := gocron.NewScheduler(time.UTC)
	_, err := s.Cron("*/1 * * * *").Do(task) // every minute
	if err != nil {
		fmt.Println("worker error", err.Error())
		return
	}
	s.StartAsync()

	//
	// Logs
	//
	handler.OpenLogFile(constant.ServerLogFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	//
	// Serve API
	//
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		fmt.Printf("defaulting to port %s\n", port)
	}

	http.HandleFunc("/docs/", httpSwagger.WrapHandler)
	http.HandleFunc("/memory", srvHandler.ServeHTTP)

	log.Fatalln(http.ListenAndServe(":"+port, handler.LogRequestMiddleware(http.DefaultServeMux)))
}

func init() {
	srvHandler = handler.Service()

	jsonFile, err := os.Open(constant.TmpDataFile)
	if err != nil {
		fmt.Println("can not open file", err.Error())
		return
	}

	if err := os.Chmod(constant.TmpDataFile, 0777); err != nil {
		fmt.Println("can not chmod file", err.Error())
		return
	}

	defer jsonFile.Close()

	byt, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("can not read file", err.Error())
		return
	}

	var result model.Result
	if err = json.Unmarshal(byt, &result); err != nil {
		fmt.Println("can not decode data", err.Error())
		return
	}

	for k, v := range result {
		go srvHandler.Cache.Set(k, v, cache.NoExpiration)
	}
}
