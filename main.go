package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const ADDRESS = "127.0.0.1:8080"
const JSON_FILE = "structure.json"

type Route struct {
	Url      string      `json: url`
	Response interface{} `json: response`
}

type ResponseFunc func(w http.ResponseWriter, r *http.Request)

func responseWrapper(route Route) ResponseFunc {
	var response = func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		log.Println(r.Method, r.URL, string(body))

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(route.Response)
	}

	return response
}

func loadPaths(router *mux.Router) {
	jsonFile, err := os.Open(JSON_FILE)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result []Route
	json.Unmarshal([]byte(byteValue), &result)

	for _, route := range result {
		router.HandleFunc(route.Url, responseWrapper(route))
	}
}

func main() {
	log.Println("Starting HTTP server on port 8080")
	router := mux.NewRouter()
	loadPaths(router)
	log.Fatal(http.ListenAndServe(":8080", router))
}
