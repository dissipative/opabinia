package app

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dissipative/opabinia/pkg/graph"

	"github.com/dissipative/opabinia/internal/infra/fs"

	mdRenderer "github.com/dissipative/opabinia/internal/markdown/renderer"
)

func (a *App) Compile() error {
	err := fs.CreateDir(fs.Dist)
	if err != nil {
		if !errors.Is(err, fs.ErrDirExist) {
			return err
		}

		proceed, e := a.tryRewriteDist()
		if e != nil {
			return err
		}
		if !proceed {
			return nil
		}
	}

	g, err := graph.BuildGraph("index.md") // todo: entrypoint as argument?

	renderer, err := mdRenderer.NewDefaultRenderer(a.config, a.logger)
	if err != nil {
		return err
	}

	err = g.Process(renderer, fs.Dist)
	if err != nil {
		return err
	}

	err = CopyFavicons(a.config.FaviconDir)
	if err != nil {
		return err
	}

	for _, asset := range []string{a.config.CustomCSS, a.config.CustomCSS, a.config.Webmanifest} {
		err = CreateAsset(asset)
		if err != nil {
			return err
		}
	}

	fmt.Printf("Compilation done to <%s> directory.\n", fs.Dist)

	return err
}

func CreateAsset(filename string) error {
	if filename == "" {
		return nil
	}
	filename = strings.TrimLeft(filename, "/")
	dst := filepath.Join(fs.Dist, filename)

	err := fs.CreateDir(filepath.Dir(dst))
	if err != nil && !errors.Is(err, fs.ErrDirExist) {
		return err
	}

	err = fs.CopyFile(filename, dst)
	if err != nil {
		return err
	}

	return nil
}

func CopyFavicons(path string) error {
	path = strings.TrimLeft(path, "/")
	dst := filepath.Join(fs.Dist, path)

	err := fs.CreateDir(dst)
	if err != nil && !errors.Is(err, fs.ErrDirExist) {
		return err
	}

	var files []string
	err = filepath.Walk(path, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// nested dirs are ignored
			return nil
		}

		files = append(files, file)
		return nil
	})

	for _, f := range files {
		err = fs.CopyFile(f, filepath.Join(fs.Dist, f))
		if err != nil {
			return err
		}
	}

	return nil
}

// tryRewriteDist checks if the "dist" directory exists and prompts the user
// if they would like to rewrite it. If yes, the directory is removed and recreated.
func (a *App) tryRewriteDist() (proceed bool, err error) {
	fmt.Printf("Directory \"%s\" exists. Rewrite? Y/n\n", fs.Dist)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	rewrite := scanner.Text()
	if strings.ToLower(rewrite) != "y" {
		return false, nil
	}

	err = os.RemoveAll(fs.Dist)
	if err != nil {
		return false, err
	}
	err = fs.CreateDir(fs.Dist) // retry creating dist
	if err != nil {
		return false, err
	}

	return true, nil
}
