/*
Copyright © 2025 Patrick Hermann patrick.hermann@sva.de
*/

package surveys

import (
	"os"
	"testing"
)

// TestReadYAML tests the ReadYAMLToMap function
func TestReadYAML(t *testing.T) {
	// Example YAML content
	yamlContent := `
renderInline:
  gitlab:
    - "tests/git.tpl"
  github:
    - "tests/git.tpl"
`

	// Write YAML content to a temporary file for testing
	tempFile := "test.yaml"
	if err := os.WriteFile(tempFile, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Error writing test file: %v", err)
	}
	defer os.Remove(tempFile) // Clean up the temporary file

	// Use the function to read the YAML
	data, err := ReadYAMLToMap(tempFile)
	if err != nil {
		t.Fatalf("Error reading YAML: %v", err)
	}

	// Validate the "renderInline" section
	renderInline, ok := data["renderInline"].(map[string]interface{})
	if !ok {
		t.Fatalf("renderInline section not found or invalid")
	}

	// Validate each key and its associated values
	tests := map[string][]string{
		"gitlab": {"tests/git.tpl"},
		"github": {"tests/git.tpl"},
	}

	for key, expectedFiles := range tests {
		values, exists := renderInline[key]
		if !exists {
			t.Errorf("Key %s not found in renderInline", key)
			continue
		}

		files, ok := values.([]interface{})
		if !ok {
			t.Errorf("Values for key %s are not of the expected type", key)
			continue
		}

		if len(files) != len(expectedFiles) {
			t.Errorf("Key %s: expected %d files, got %d", key, len(expectedFiles), len(files))
			continue
		}

		for i, file := range files {
			if file != expectedFiles[i] {
				t.Errorf("Key %s: expected file %s, got %s", key, expectedFiles[i], file)
			}
		}
	}
}
