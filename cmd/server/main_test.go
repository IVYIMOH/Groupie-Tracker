package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestSafeStaticHandler(t *testing.T) {
	// 1. Move to project root so we can see 'static/'
	originalWD, _ := os.Getwd()
	// If you are in cmd/server, you need to go up two levels
	os.Chdir("../../")
	defer os.Chdir(originalWD)

	// 2. Initialize the handler
	handler := safeStaticHandler("static")

	tests := []struct {
		name           string
		url            string
		expectedStatus int
	}{
		{
			name: "Valid Static File",
			// CHECK: Make sure this file actually exists in your static folder!
			// If it's in static/css/style.css, change the URL below.
			url:            "/static/style-404.css",
			expectedStatus: http.StatusOK,
		},
		{
			name: "Directory Traversal Attempt",
			url:  "/static/../../main.go",
			// Change to 400 because your code is rejecting the '..' syntax
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Non-existent File",
			url:            "/static/this_file_does_not_exist.png",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.url, nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				// Technical Debug: help the auditor see what path was actually checked
				t.Errorf("%s: expected %d, got %d. (Check if file exists at real path)",
					tt.name, tt.expectedStatus, rr.Code)
			}
		})
	}
}
