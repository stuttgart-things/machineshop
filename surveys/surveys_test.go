/*
Copyright © 2025 Patrick Hermann patrick.hermann@sva.de
*/

package surveys

import (
	"bytes"
	"strings"
	"testing"
	"text/template"
)

func TestRandomFromSlice(t *testing.T) {
	// Test case 1: Non-empty slice
	inputSlice := []string{"apple", "banana", "cherry"}
	randomValue := RandomFromSlice(inputSlice)

	// Check if the random value is in the slice
	found := false
	for _, v := range inputSlice {
		if v == randomValue {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected one of %v, but got %s", inputSlice, randomValue)
	}

	// Test case 2: Empty slice
	emptySlice := []string{}
	randomValue = RandomFromSlice(emptySlice)
	if randomValue != "" {
		t.Errorf("Expected empty string for empty slice, but got %s", randomValue)
	}
}

// TestRenderTemplateInlineWithFunction tests the rendering of a Go template
// with inline functions. It defines a simple template function map with
// functions to convert a string to uppercase and to add two integers. The
// template data includes calls to these functions, and the template variables
// provide the necessary data for rendering. The test captures the rendered
// output and compares it with the expected output to ensure correctness.
func TestRenderTemplateInlineWithFunction(t *testing.T) {
	// Define a simple template function map
	templateFunctions := template.FuncMap{
		"toUpperCase": func(input string) string {
			return strings.ToUpper(input)
		},
		"add": func(a, b int) int {
			return a + b
		},
	}

	// Define template data with function calls
	templateData := `
		Hello, {{ toUpperCase .Name }}!
		The sum of 2 and 3 is {{ add 2 3 }}.
	`

	// Define template variables
	templateVariables := map[string]interface{}{
		"Name": "world",
	}

	// Buffer to capture the rendered output
	var renderedOutput bytes.Buffer

	// Wrap the rendering in a function to capture the output
	render := func() {
		tmpl, err := template.New("template").Funcs(templateFunctions).Parse(templateData)
		if err != nil {
			t.Fatalf("Failed to parse template: %v", err)
		}

		err = tmpl.Execute(&renderedOutput, templateVariables)
		if err != nil {
			t.Fatalf("Failed to execute template: %v", err)
		}
	}

	// Call the render function
	render()

	// Expected output
	expectedOutput := `
		Hello, WORLD!
		The sum of 2 and 3 is 5.
	`

	// Compare rendered output with expected output
	if renderedOutput.String() != expectedOutput {
		t.Errorf("Rendered output does not match expected output. Got:\n%s\nExpected:\n%s\n",
			renderedOutput.String(), expectedOutput)
	}
}
