package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"os"
)

func TestGetLastModifiedTime(t *testing.T) {
	filePath := "test.txt"

	// Create a test file
	err := ioutil.WriteFile(filePath, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(filePath)

	lastModified := getLastModifiedTime(filePath)
	if lastModified.IsZero() {
		t.Error("getLastModifiedTime returned zero time")
	}
}

func TestRouteHandler(t *testing.T) {
	// Load routes for testing
	readRoutes()
	lastModifiedTime = getLastModifiedTime(JSONFile)

	testCases := []struct {
		name       string
		url        string
		statusCode int
	}{
		{
			name:       "Valid Route",
			url:        "/api/users",
			statusCode: http.StatusOK,
		},
		{
			name:       "Invalid Route",
			url:        "/api/unknown",
			statusCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tc.url, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(routeHandler)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.statusCode {
				t.Errorf("Handler returned wrong status code: got %v want %v",
					status, tc.statusCode)
			}
		})
	}
}

func TestReadRoutes(t *testing.T) {
	readRoutes()
	_, ok := routes.Load("/api/users")
	if !ok {
		t.Error("Failed to read routes from JSON file")
	}
}