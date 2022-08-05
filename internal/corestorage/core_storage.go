package corestorage

import (
	"github.com/balerter/balerter/internal/alert"
	"net/http"
)

//go:generate moq -out module_alert.go -skip-ensure -fmt goimports . Alert
//go:generate moq -out module_kv.go -skip-ensure -fmt goimports . KV
//go:generate moq -out module_core_storage.go -skip-ensure -fmt goimports . CoreStorage

// KV is an interface for KV storage
type KV interface {
	Put(string, string) error
	Get(string) (string, error)
	Upsert(string, string) error
	Delete(string) error
	All() (map[string]string, error)
	RunApiHandler(rw http.ResponseWriter, req *http.Request)
}

// Alert is an interface for Alert storage
type Alert interface {
	// Update exists alert or create new
	Update(name string, level alert.Level) (*alert.Alert, bool, error)
	Index(levels []alert.Level) (alert.Alerts, error)
	Get(name string) (*alert.Alert, error)
	RunApiHandler(rw http.ResponseWriter, req *http.Request)
}

// CoreStorage is an interface for the CoreStorage
type CoreStorage interface {
	Name() string
	KV() KV
	Alert() Alert
	Stop() error
}
