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
	Url      string      `json:"url"`
	Response interface{} `json:"response"`
}

type ResponseFunc func(w http.ResponseWriter, r *http.Request)

func responseWrapper(route Route) ResponseFunc {
	var response = func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		log.Println(r.Method, r.URL, string(body))

		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(route.Response)
	}

	return response
}

func loadPaths() *mux.Router {
	router := mux.NewRouter()

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

	return router
}

func main() {
	router := loadPaths()
	server := http.Server{Addr: ADDRESS, Handler: router}

	// watcher, err := fsnotify.NewWatcher()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer watcher.Close()

	// done := make(chan bool)
	// go func() {
	// 	for {
	// 		select {
	// 		case event, ok := <-watcher.Events:
	// 			if !ok {
	// 				return
	// 			}
	// 			log.Println("event:", event)
	// 			if event.Op&fsnotify.Write == fsnotify.Write {
	// 				log.Println("modified file:", event.Name)
	// 				log.Println("Restaring server...")

	// 				server.Shutdown(context.Background())
	// 				router := loadPaths()
	// 				server.Handler = router
	// 				server.ListenAndServe()
	// 			}
	// 		case err, ok := <-watcher.Errors:
	// 			if !ok {
	// 				return
	// 			}
	// 			log.Println("error:", err)
	// 		}
	// 	}
	// }()
	// err = watcher.Add(JSON_FILE)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	log.Println("Starting HTTP server on port 8080")
	// log.Fatal(http.ListenAndServe(":8080", router))
	log.Fatal(server.ListenAndServe())
	// <-done
}
