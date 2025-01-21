/*
Copyright © 2025 Patrick Hermann patrick.hermann@sva.de
*/

package internal

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendToHomerun(t *testing.T) {
	// Mock server to simulate the homerun service
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate the Content-Type header
		if r.Header.Get("Content-Type") != contentType {
			t.Errorf("Expected Content-Type %s, got %s", contentType, r.Header.Get("Content-Type"))
		}

		// Validate the X-Auth-Token header
		if r.Header.Get("X-Auth-Token") != "test-token" {
			t.Errorf("Expected X-Auth-Token test-token, got %s", r.Header.Get("X-Auth-Token"))
		}

		// Write a response
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"success"}`))
	}))
	defer mockServer.Close()

	// Test data
	destination := mockServer.URL
	token := "test-token"
	renderedBody := []byte(`{"message":"hello"}`)

	// Call the function
	response, resp := SendToHomerun(destination, token, renderedBody)

	fmt.Println(resp)
	// Verify the response
	expectedResponse := `{"status":"success"}`
	if string(response) != expectedResponse {
		t.Errorf("Expected response %s, got %s", expectedResponse, string(response))
	}
}
