package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

const JSONFile = "structure.json"

var routes sync.Map
var lastModifiedTime time.Time

type Route struct {
	URL      string      `json:"url"`
	Response interface{} `json:"response"`
}

func getLastModifiedTime(filePath string) time.Time {
	file, err := os.Stat(filePath)
	if err != nil {
		log.Fatal(err)
	}

	return file.ModTime()
}

func readRoutes() {
	jsonFile, err := os.Open(JSONFile)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result []Route
	json.Unmarshal(byteValue, &result)

	for _, route := range result {
		routes.Store(route.URL, route.Response)
	}
}

func routeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == "OPTIONS" {
		return
	}

	val, ok := routes.Load(r.URL.String())
	if ok {
		json.NewEncoder(w).Encode(val)
	} else {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusNotFound)
	}
}

func main() {
	host := flag.String("host", "127.0.0.1", "Specify host to serve backend")
	port := flag.String("port", "8080", "Specify port to use")
	flag.Parse()

	address := *host + ":" + *port
	readRoutes()
	lastModifiedTime = getLastModifiedTime(JSONFile)

	go func() {
		for {
			time.Sleep(1 * time.Second)
			newLastModifiedTime := getLastModifiedTime(JSONFile)
			if newLastModifiedTime.After(lastModifiedTime) {
				routes.Range(func(key, _ interface{}) bool {
					routes.Delete(key)
					return true
				})
				readRoutes()
				lastModifiedTime = newLastModifiedTime
			}
		}
	}()

	log.Println("Starting server on ", address)
	log.Fatal(http.ListenAndServe(address, http.HandlerFunc(routeHandler)))
}