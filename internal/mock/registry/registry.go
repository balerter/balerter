package registry

import (
	"errors"
	lua "github.com/yuin/gopher-lua"
)

var (
	// ErrEntryIsNotRegistered represents EntryIsNotRegistered error
	ErrEntryIsNotRegistered = errors.New("response entry is not registered")
)

// Registry for registration responses by method name and arguments
// and asserts with calls
type Registry struct {
	responseEntries map[string]*responseEntry
	assertEntries   map[string]*assertEntry
	calls           []call
}

type call struct {
	method string
	args   []lua.LValue
}

type assertEntry struct {
	entries map[string]*assertEntry
	asserts []bool
}

type responseEntry struct {
	entries   map[string]*responseEntry
	responses [][]lua.LValue
}

func newAssertEntry() *assertEntry {
	e := &assertEntry{
		entries: make(map[string]*assertEntry),
	}
	return e
}

func newResponseEntry() *responseEntry {
	e := &responseEntry{
		entries: make(map[string]*responseEntry),
	}
	return e
}

// New creates new registry
func New() *Registry {
	r := &Registry{
		responseEntries: map[string]*responseEntry{},
		assertEntries:   map[string]*assertEntry{},
	}

	return r
}

// Clean the registry
func (r *Registry) Clean() {
	for key := range r.responseEntries {
		delete(r.responseEntries, key)
	}
	for key := range r.assertEntries {
		delete(r.assertEntries, key)
	}
	r.calls = r.calls[:0]
}
