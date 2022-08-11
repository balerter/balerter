package memory

import (
	"github.com/balerter/balerter/internal/alert"
	coreStorage "github.com/balerter/balerter/internal/corestorage"
	"net/http"
	"sync"
)

type storageKV struct {
	mxKV sync.RWMutex
	kv   map[string]string
}

type storageAlert struct {
	mxAlerts sync.RWMutex
	alerts   map[string]*alert.Alert
}

// Memory represent inMemory storage engine
type Memory struct {
	kv    *storageKV
	alert *storageAlert
}

// New creates new Memory storage
func New() *Memory {
	m := &Memory{
		kv: &storageKV{
			kv: make(map[string]string),
		},
		alert: &storageAlert{
			alerts: make(map[string]*alert.Alert),
		},
	}

	return m
}

func (m *storageKV) RunApiHandler(rw http.ResponseWriter, req *http.Request) {
	http.Error(rw, "coreapi is not supported for this module", http.StatusNotImplemented)
}

func (m *storageAlert) RunApiHandler(rw http.ResponseWriter, req *http.Request) {
	http.Error(rw, "coreapi is not supported for this module", http.StatusNotImplemented)
}

// Name returns the name of the storage
func (m *Memory) Name() string {
	return "memory"
}

// KV returns KV storage
func (m *Memory) KV() coreStorage.KV {
	return m.kv
}

// Alert returns Alert storage
func (m *Memory) Alert() coreStorage.Alert {
	return m.alert
}

// Stop the engine
func (m *Memory) Stop() error {
	return nil
}
