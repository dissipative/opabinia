package app

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dissipative/opabinia/internal/infra/fs"

	"github.com/dissipative/opabinia/layout"
)

type embeddedFile struct {
	name    string
	dir     string
	content []byte
}

var projectFiles = []embeddedFile{
	{
		name:    "default.tmpl",
		dir:     "templates",
		content: layout.Default,
	},
	{
		name:    "index.md",
		dir:     ".",
		content: layout.Index,
	},
	{
		name:    "other_page.md",
		dir:     "pages",
		content: layout.OtherPage,
	},
	{
		name:    "styles.css",
		dir:     "assets",
		content: nil,
	},
	{
		name:    "apple-touch-icon.png",
		dir:     "assets/favicon",
		content: layout.AppleTouchIcon,
	},
	{
		name:    "favicon-16x16.png",
		dir:     "assets/favicon",
		content: layout.Favicon16x16,
	},
	{
		name:    "favicon-32x32.png",
		dir:     "assets/favicon",
		content: layout.Favicon32x32,
	},
	{
		name:    "opabinia.yml",
		content: layout.Config,
	},
}

// DoInit initializes a new project with a specified structure and default content.
// It performs the following steps:
// 1. Reads the provided project Name from command-line arguments.
// 2. Creates the root project directory.
// 3. Creates subdirectories and populates them with templates, sample files, and assets.
func (a *App) DoInit() error {
	name, err := a.readProjectName()
	if err != nil {
		return err
	}

	err = fs.CreateDir(name)
	if err != nil {
		return err
	}

	for _, file := range projectFiles {
		dirName := filepath.Join(name, file.dir)
		err = fs.CreateDir(dirName)
		if err != nil && !errors.Is(err, fs.ErrDirExist) {
			return err
		}

		err = os.WriteFile(filepath.Join(dirName, file.name), file.content, fs.DefaultFilePerm)
		if err != nil {
			return fmt.Errorf("failed to write file %s: %w", filepath.Join(dirName, file.name), err)
		}
	}

	fmt.Printf("New project initialized to <%s> directory.\n", name)

	return nil
}

// readProjectName retrieves the project Name from command-line arguments.
// It searches for a valid Name provided after the "-i" or "--init" flags.
func (a *App) readProjectName() (string, error) {
	var name string
	for i, arg := range os.Args {
		if isInit(arg) && i+1 < len(os.Args) {
			name = os.Args[i+1]
			break
		}
	}

	if name == "" {
		return "", errors.New("please provide project Name")
	}
	if strings.HasPrefix(name, "-") {
		return "", errors.New("please provide project Name, option given instead")
	}

	return name, nil
}

// isInit checks if the given argument corresponds to the initialization command ("-i" or "--init").
func isInit(arg string) bool {
	return arg == "-i" || arg == "--init"
}
