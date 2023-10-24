package graph

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/dissipative/opabinia/internal/infra/fs"
	"github.com/dissipative/opabinia/internal/markdown/parser"

	"golang.org/x/sync/semaphore"
)

type DependencyType int

const (
	asset DependencyType = iota
	markdown
)

type DependencyTree struct {
	filename string
	depType  DependencyType
	children []*DependencyTree
}

var maxConcurrentRoutines = int64(runtime.NumCPU())

// NewDependencyTree creates a new DependencyTree for the given filename.
// If the file is a local markdown file, it extracts all the links within it and
// creates child dependencies recursively.
func NewDependencyTree(filename string, paths Paths) (*DependencyTree, error) {
	if !fs.IsLocalMarkdownFile(filename) {
		return &DependencyTree{
			filename: filename,
			depType:  asset,
			children: nil,
		}, nil
	}

	var children []*DependencyTree
	var ee []error

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", filename, err)
	}

	links := parser.ExtractLinks(data)
	data = nil // explicitly remove file contents from memory

	errChan := make(chan error, len(links))
	childChan := make(chan *DependencyTree, len(links))
	wg := new(sync.WaitGroup)
	sem := semaphore.NewWeighted(maxConcurrentRoutines) // limiting concurrency

	for _, link := range links {
		if !fs.IsLocalFile(link) {
			// web urls are not needed
			continue
		}

		path := fs.NormalizePath(filename, link)
		if paths.IsSet(path) {
			// circular dependency, it is ok for hypertext; just omit
			continue
		}
		paths.Set(path)

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer sem.Release(1)

			err := sem.Acquire(context.Background(), 1)
			if err != nil {
				errChan <- err
				return
			}

			child, err := NewDependencyTree(path, paths)
			if err != nil {
				errChan <- err
				return
			}
			childChan <- child
		}()
	}

	go func() {
		wg.Wait()
		close(errChan)
		close(childChan)
	}()

	for e := range errChan {
		ee = append(ee, e)
	}

	for ch := range childChan {
		children = append(children, ch)
	}

	if len(ee) > 0 {
		return nil, errors.Join(ee...)
	}

	return &DependencyTree{
		filename: filename,
		depType:  markdown,
		children: children,
	}, nil
}

func (d *DependencyTree) IsMarkdown() bool {
	return d.depType == markdown
}

// processDependency processes the DependencyTree, rendering it to html if it is a markdown file
// and copying it to a destination directory.
// If it's a markdown file, it also processes its child dependencies recursively.
// Errors encountered during processing are sent to the provided error channel.
func (d *DependencyTree) processDependency(r Renderer, dst string, errChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	newName := filepath.Join(dst, d.filename)

	// create dir for new file
	err := fs.CreateDir(filepath.Dir(newName))
	if err != nil && !errors.Is(err, fs.ErrDirExist) {
		errChan <- err
		return
	}

	if !d.IsMarkdown() {
		err = fs.CopyFile(d.filename, newName)
		if err != nil {
			errChan <- err
			return
		}
	} else {
		content, err := r.RenderPage(d.filename)
		if err != nil {
			errChan <- err
			return
		}

		newName = filepath.Join(dst, strings.TrimSuffix(d.filename, ".md")+".html")
		err = os.WriteFile(newName, content, fs.DefaultFilePerm)
		content = nil // explicitly remove content from memory as soon as possible
		if err != nil {
			errChan <- err
			return
		}

		for _, child := range d.children {
			wg.Add(1)
			go child.processDependency(r, dst, errChan, wg)
		}
	}
}
