package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const ADDRESS = "127.0.0.1:8080"
const JSON_FILE = "structure.json"

type Route struct {
	Url      string      `json:"url"`
	Response interface{} `json:"response"`
}

type ResponseFunc func(w http.ResponseWriter, r *http.Request)

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
	paths := readPaths()

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
	// Read modified date of JSON_FILE
	// Then check if it changed in handler func
	// And re-read contents if necessary
	log.Println("Starting server on ", ADDRESS)
	log.Fatal(http.ListenAndServe(ADDRESS, http.HandlerFunc(myHandler)))
}
