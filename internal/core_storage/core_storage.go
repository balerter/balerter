package core_storage

import (
	"github.com/balerter/balerter/internal/alert/alert"
)

type CoreStorageKV interface {
	Name() string
	Put(string, string) error
	Get(string) (string, error)
	Upsert(string, string) error
	Delete(string) error
}

type CoreStorageAlert interface {
	GetOrNew(string) (*alert.Alert, error)
	All() ([]*alert.Alert, error)
	Release(a *alert.Alert)
}

type CoreStorage interface {
	CoreStorageKV
	CoreStorageAlert
}
