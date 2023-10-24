package fs

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateDir(t *testing.T) {
	dirName := "testDir"
	parentDir := "parentDir"

	// Cleanup after tests
	defer func() {
		_ = os.RemoveAll(dirName)
		_ = os.RemoveAll(parentDir)
	}()

	// Test creating a new directory
	err := CreateDir(dirName)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Test attempting to create a directory that already exists
	err = CreateDir(dirName)
	if !errors.Is(err, ErrDirExist) {
		t.Fatalf("Expected ErrDirExist, got %v", err)
	}

	// Test attempting to create a directory where parent doesn't exist
	err = CreateDir(filepath.Join(parentDir, dirName))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestNormalizePath(t *testing.T) {
	tests := []struct {
		parent   string
		filename string
		expected string
	}{
		{
			parent:   "/home/user/docs/page.md",
			filename: "index.md",
			expected: "/home/user/docs/index.md",
		},
		{
			parent:   "/home/user/docs/page.md",
			filename: "../index.md",
			expected: "/home/user/index.md",
		},
		{
			parent:   "/home/user/docs/page.md",
			filename: "../../index.md",
			expected: "/home/index.md",
		},
		{
			parent:   "/home/user/docs/page.md",
			filename: "../../../index.md",
			expected: "/index.md",
		},
		// for the following case: the function does not preserve "../" for the level exceeding root in the parent path;
		// that's why the result is not "/../index.md" but "/index.md"
		{
			parent:   "/home/user/docs.page.md",
			filename: "../../../../index.md",
			expected: "/index.md",
		},
	}
	for _, test := range tests {
		t.Run(test.filename, func(t *testing.T) {
			result := NormalizePath(test.parent, test.filename)
			if result != test.expected {
				t.Errorf("For parent %s and filename %s, expected %s but got %s", test.parent, test.filename, test.expected, result)
			}
		})
	}
}
