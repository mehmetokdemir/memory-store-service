package main

import (
	// Go imports
	"encoding/json"
	"fmt"
	httpSwagger "github.com/swaggo/http-swagger"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
	"workout/memory-store-service/constant"

	"github.com/go-co-op/gocron"
	"github.com/patrickmn/go-cache"
	// Internal imports
	_ "workout/memory-store-service/docs"
	"workout/memory-store-service/internal/handler"
	"workout/memory-store-service/model"
)

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
		fmt.Println("can not decode data", err.Error())
		return
	}
	if err := os.WriteFile(constant.TmpDataFile, b, os.ModePerm); err != nil {
		fmt.Println("can not write files", err.Error())
		return
	}

	fmt.Println("worker result", result)
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

	log.Fatalln(http.ListenAndServe(":"+port, handler.LogRequest(http.DefaultServeMux)))
}

func init() {
	srvHandler = handler.Service()
	if err := os.Chmod(constant.TmpDataFile, 0777); err != nil {
		fmt.Println("can not chmod file", err.Error())
		return
	}

	jsonFile, err := os.Open(constant.TmpDataFile)
	if err != nil {
		fmt.Println("can not open file", err.Error())
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
		fmt.Println("k", k)
		go srvHandler.Cache.Set(k, v, cache.NoExpiration)
	}
}
