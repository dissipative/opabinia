package fs

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const Dist = "dist"

const (
	DefaultFilePerm      = 0o644
	DefaultDirectoryPerm = 0o755
)

func CopyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("check path for <%s> [%w]", src, err)
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, DefaultFilePerm)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)

	return err
}

var ErrDirExist = errors.New("directory already exists")

func CreateDir(name string) error {
	info, err := os.Stat(name)
	if err == nil && info.IsDir() {
		return fmt.Errorf("<%s>: %w", name, ErrDirExist)
	}
	if os.IsNotExist(err) {
		err = os.MkdirAll(name, DefaultDirectoryPerm)
		if err != nil {
			return err
		}
	}

	return nil
}

// IsLocalMarkdownFile checks if the given filename is a local markdown file.
// It checks for the ".md" file extension and ensures it's not a URL.
func IsLocalMarkdownFile(filename string) bool {
	return strings.HasSuffix(filename, ".md") && IsLocalFile(filename)
}

// IsLocalFile checks if the given filename is a local file and not a URL.
func IsLocalFile(filename string) bool {
	return !strings.Contains(filename, ":")
}

// NormalizePath adjusts the given filename to a standardized format.
// If the filename contains relative references like "../", it resolves them
// based on the provided parent directory. The function ensures the returned
// path is in the format "dir/dir/file.ext".
func NormalizePath(parent string, filename string) string {
	baseDir := filepath.Dir(parent)

	// Iterate over the "../" prefixes and move up the directory tree
	for strings.HasPrefix(filename, "../") {
		filename = strings.TrimPrefix(filename, "../")
		baseDir = filepath.Dir(baseDir)

	}

	return filepath.Join(baseDir, filepath.Clean(filename))
}
