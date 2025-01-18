/*
Copyright © 2025 Patrick Hermann patrick.hermann@sva.de
*/

package surveys

import (
	"testing"
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
