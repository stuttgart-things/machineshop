/*
Copyright © 2025 Patrick Hermann patrick.hermann@sva.de
*/

package internal

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestExecuteShellScript tests the ExecuteShellScript function
func TestExecuteShellScript(t *testing.T) {
	// Create a temporary directory for log files
	tempDir, err := ioutil.TempDir("", "test-logs")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up after the test

	// Define test cases
	tests := []struct {
		name         string
		scriptName   string
		script       string
		expectError  bool
		expectedLogs []string
	}{
		{
			name:        "Valid Script",
			scriptName:  "test-valid",
			script:      `echo "Hello, World!"`,
			expectError: false,
			expectedLogs: []string{
				"Hello, World!",
			},
		},
	}

	// Run each test case
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logFile := filepath.Join(tempDir, tc.scriptName+".log")

			// Capture the panic to test error handling
			defer func() {
				if r := recover(); r != nil {
					if !tc.expectError {
						t.Fatalf("Unexpected panic: %v", r)
					}
				}
			}()

			// Call the function
			ExecuteShellScript(tc.scriptName, tc.script, tempDir)

			// Read the log file
			logContents, err := ioutil.ReadFile(logFile)
			if err != nil {
				t.Fatalf("Failed to read log file: %v", err)
			}

			// Check expected log contents
			for _, expectedLog := range tc.expectedLogs {
				if !strings.Contains(string(logContents), expectedLog) {
					t.Errorf("Expected log to contain: %q, but got: %s", expectedLog, string(logContents))
				}
			}
		})
	}
}
