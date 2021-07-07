package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const ADDRESS = "127.0.0.1:8080"
const JSON_FILE = "structure.json"

var paths map[string]interface{}

type Route struct {
	Url      string      `json:"url"`
	Response interface{} `json:"response"`
}

type ResponseFunc func(w http.ResponseWriter, r *http.Request)

func getLastModifiedTime(filePath string) time.Time {
	file, err := os.Stat(JSON_FILE)

	if err != nil {
		log.Fatal(err)
	}

	return file.ModTime()
}

func readPaths() map[string]interface{} {
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

func myHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	if val, ok := paths[r.URL.String()]; ok {
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
	paths = readPaths()

	go func() {
		lastModifiedTime := getLastModifiedTime(JSON_FILE)

		for {
			time.Sleep(1 * time.Second)
			newLastModifiedTime := getLastModifiedTime(JSON_FILE)
			if newLastModifiedTime.After(lastModifiedTime) {
				paths = readPaths()
				lastModifiedTime = newLastModifiedTime
			}
		}
	}()

	log.Println("Starting server on ", ADDRESS)
	log.Fatal(http.ListenAndServe(ADDRESS, http.HandlerFunc(myHandler)))
}
