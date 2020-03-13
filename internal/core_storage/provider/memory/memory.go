package memory

import (
	"github.com/balerter/balerter/internal/alert/alert"
	"sync"
)

type Memory struct {
	mxKV sync.RWMutex
	kv   map[string]string

	mxAlerts sync.RWMutex
	alerts   map[string]*alert.Alert
}

func New() *Memory {
	m := &Memory{
		kv:     make(map[string]string),
		alerts: make(map[string]*alert.Alert),
	}

	return m
}

func (m *Memory) Name() string {
	return "memory"
}
