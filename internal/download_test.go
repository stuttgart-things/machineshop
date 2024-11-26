package internal

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDownloadFile(t *testing.T) {
	// Create a test HTTP server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a successful response
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("Test content of the file"))
		if err != nil {
			fmt.Println("ERROR BY WRITING RESPONSE: ", err)
		}
	}))
	defer ts.Close()

	// Call the DownloadFile function with the test URL
	content, err := DownloadFile(ts.URL)
	if err != nil {
		t.Fatalf("DownloadFile returned an error: %v", err)
	}

	// Verify the content of the downloaded file
	expected := "Test content of the file"
	if content != expected {
		t.Errorf("Expected '%s', but got '%s'", expected, content)
	}
}
