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
	// Update exists alert or create new
	Update(name string, level alert.Level) (*alert.Alert, bool, error)
	Index(levels []alert.Level) (alert.Alerts, error)
}

type CoreStorage interface {
	Name() string
	KV() KV
	Alert() Alert
	Stop() error
}
