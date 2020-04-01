package registry

import (
	"errors"
	lua "github.com/yuin/gopher-lua"
)

var (
	ErrEntryIsNotRegistered = errors.New("response entry is not registered")
)

// Registry for registration responses by method name and arguments
// and asserts with calls
type Registry struct {
	responseEntries map[string]*responseEntry
}

type responseEntry struct {
	entries   map[string]*responseEntry
	responses [][]lua.LValue
}

func newResponseEntry() *responseEntry {
	e := &responseEntry{
		entries: make(map[string]*responseEntry),
	}
	return e
}

func New() *Registry {
	r := &Registry{
		responseEntries: map[string]*responseEntry{},
	}

	return r
}

func (r *Registry) Clean() {
	for key := range r.responseEntries {
		delete(r.responseEntries, key)
	}
}
