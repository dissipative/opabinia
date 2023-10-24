package graph

import (
	"sync"

	"github.com/dissipative/opabinia/internal/infra/fs"
)

// Graph is a dependency graph.
type Graph struct {
	root  *DependencyTree
	paths map[string]struct{}
}

// BuildGraph builds a graph from the given entry point.
func BuildGraph(entryPoint string) (*Graph, error) {
	paths := NewPaths(map[string]struct{}{
		fs.NormalizePath(".", entryPoint): {},
	})

	root, err := NewDependencyTree(entryPoint, *paths)
	if err != nil {
		return nil, err
	}

	return &Graph{root: root, paths: paths.store}, nil
}

// Renderer is the interface that wraps the RenderPage method.
type Renderer interface {
	RenderPage(mdFilename string) ([]byte, error)
}

// Process processes the graph by rendering each markdown file in parallel.
func (g *Graph) Process(r Renderer, dst string) error {
	wg := new(sync.WaitGroup)
	errChan := make(chan error, len(g.paths))

	wg.Add(1)
	go g.root.processDependency(r, dst, errChan, wg)
	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		return err
	}

	return nil
}
