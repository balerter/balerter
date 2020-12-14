package corestorage

import (
	alert2 "github.com/balerter/balerter/internal/alert"
)

type KV interface {
	Put(string, string) error
	Get(string) (string, error)
	Upsert(string, string) error
	Delete(string) error
	All() (map[string]string, error)
}

type Alert interface {
	GetOrNew(string) (*alert2.Alert, error)
	All() ([]*alert2.Alert, error)
	Release(a *alert2.Alert) error
	Get(string) (*alert2.Alert, error)
}

type CoreStorage interface {
	Name() string
	KV() KV
	Alert() Alert
	Stop() error
}
