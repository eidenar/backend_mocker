package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const JSON_FILE = "structure.json"

var routes map[string]interface{}

type Route struct {
	Url      string      `json:"url"`
	Response interface{} `json:"response"`
}

func getLastModifiedTime(filePath string) time.Time {
	file, err := os.Stat(filePath)
	if err != nil {
		log.Fatal(err)
	}

	return file.ModTime()
}

func readRoutes() map[string]interface{} {
	jsonFile, err := os.Open(JSON_FILE)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result []Route
	json.Unmarshal([]byte(byteValue), &result)

	mapping := make(map[string]interface{})
	for _, route := range result {
		mapping[route.Url] = route.Response
	}

	return mapping
}

func routeHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	if val, ok := routes[r.URL.String()]; ok {
		log.Println(r.Method, r.URL, string(body))

		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(val)
	} else {
		log.Println(r.Method, r.URL, "404 NOT FOUND")
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusNotFound)
	}
}

func main() {
	host := flag.String("host", "127.0.0.1", "Specify host to serve backend")
	port := flag.String("port", "8080", "Specify port to use")
	flag.Parse()

	address := *host + ":" + *port
	routes = readRoutes() // Read routes for 1st time

	go func() {
		lastModifiedTime := getLastModifiedTime(JSON_FILE)

		for {
			time.Sleep(1 * time.Second)
			newLastModifiedTime := getLastModifiedTime(JSON_FILE)
			if newLastModifiedTime.After(lastModifiedTime) {
				routes = readRoutes()
				lastModifiedTime = newLastModifiedTime
			}
		}
	}()

	log.Println("Starting server on ", address)
	log.Fatal(http.ListenAndServe(address, http.HandlerFunc(routeHandler)))
}
