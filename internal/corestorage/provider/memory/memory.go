package memory

import (
	"github.com/balerter/balerter/internal/alert"
	coreStorage "github.com/balerter/balerter/internal/corestorage"
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

type Memory struct {
	kv    *storageKV
	alert *storageAlert
}

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

func (m *Memory) Name() string {
	return "memory"
}

func (m *Memory) KV() coreStorage.KV {
	return m.kv
}

func (m *Memory) Alert() coreStorage.Alert {
	return m.alert
}

func (m *Memory) Stop() error {
	return nil
}
