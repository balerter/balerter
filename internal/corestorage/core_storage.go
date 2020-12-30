package corestorage

import (
	"github.com/balerter/balerter/internal/alert"
)

type KV interface {
	Put(string, string) error
	Get(string) (string, error)
	Upsert(string, string) error
	Delete(string) error
	All() (map[string]string, error)
}

type Alert interface {
	GetOrNew(string) (*alert.Alert, error)
	All() ([]*alert.Alert, error)
	Store(a *alert.Alert) error
	Get(string) (*alert.Alert, error)
}

type CoreStorage interface {
	Name() string
	KV() KV
	Alert() Alert
	Stop() error
}
