package app

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestCommander_DoInit(t *testing.T) {
	// Setup: Mock App and create temporary directory for testing
	c := &App{}
	dir := "InitTestProject"
	defer os.RemoveAll(dir)

	// Set fake arguments for the initialization command
	os.Args = []string{Name, "-i", dir}

	// Execute the DoInit function
	err := c.DoInit()
	if err != nil {
		t.Fatalf("DoInit() failed with error: %v", err)
	}

	// Check if the directories and files have been created
	for _, file := range projectFiles {
		_, err := os.Stat(filepath.Join(dir, file.dir, file.name))
		if os.IsNotExist(err) {
			t.Errorf("Expected file %s to exist, but it doesn't", filepath.Join(file.dir, file.name))
		}
	}
}

func TestCommander_readProjectName(t *testing.T) {
	c := &App{}

	tests := []struct {
		args     []string
		expected string
		err      error
	}{
		{[]string{Name, "-i", "projectName"}, "projectName", nil},
		{[]string{Name, "--init", "projectName"}, "projectName", nil},
		{[]string{Name, "--init"}, "", errors.New("please provide project Name")},
		{[]string{Name, "--init", "-anotherFlag"}, "", errors.New("please provide project Name, option given instead")},
	}

	for _, tt := range tests {
		os.Args = tt.args
		name, err := c.readProjectName()

		if name != tt.expected {
			t.Errorf("Expected project Name %s but got %s", tt.expected, name)
		}

		if err != nil && err.Error() != tt.err.Error() {
			t.Errorf("Expected error %v but got %v", tt.err, err)
		}
	}
}
