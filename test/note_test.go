package test

import (
	"fmt"
	"testing"
    "github.com/nekowawolf/notes/module"
)

func TestInsertNote(t *testing.T) {
	title := "Unit Test Note"
	content := "This is a test note for unit testing."

	result, err := module.InsertNote(title, content)
	if err != nil {
		t.Errorf("Failed to insert note: %v", err)
		return
	}

	fmt.Printf("Inserted Note ID: %v\n", result)
}