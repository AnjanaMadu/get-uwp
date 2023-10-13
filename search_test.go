package main

import (
	"encoding/json"
	"testing"
)

func TestSearchStore(t *testing.T) {
	// Test case 1: Test searchStore with a valid query
	results, err := searchStore("calculator")
	if err != nil {
		t.Errorf("searchStore returned an error: %v", err)
	}
	if len(results) == 0 {
		t.Errorf("searchStore returned no results")
	}

	json, _ := json.MarshalIndent(results, "", "  ")
	t.Logf("%s", json)
}

func TestGetFiles(t *testing.T) {
	// Test case 1: Test getFiles with a valid product ID
	files, err := getFiles("9WZDNCRFHVN5")
	if err != nil {
		t.Errorf("getFiles returned an error: %v", err)
	}
	if len(files) == 0 {
		t.Errorf("getFiles returned no files")
	}

	json, _ := json.MarshalIndent(files, "", "  ")
	t.Logf("%s", json)
}
