package graph

import "sync"

// Paths is a thread-safe map of paths.
type Paths struct {
	store map[string]struct{}
	*sync.RWMutex
}

// NewPaths creates a new Paths.
func NewPaths(store map[string]struct{}) *Paths {
	return &Paths{
		store:   store,
		RWMutex: new(sync.RWMutex),
	}
}

// Set sets the given path.
func (p *Paths) Set(path string) {
	p.Lock()
	defer p.Unlock()

	p.store[path] = struct{}{}
}

// IsSet returns true if the given path is set.
func (p *Paths) IsSet(path string) bool {
	p.RLock()
	defer p.RUnlock()

	_, exists := p.store[path]
	return exists
}
