package memory

import (
	"fmt"
	"sync"
)

type Memory struct {
	mx      sync.RWMutex
	storage map[string]string
}

func New() *Memory {
	m := &Memory{
		storage: make(map[string]string),
	}

	return m
}

func (m *Memory) Name() string {
	return "memory"
}

func (m *Memory) Put(name string, val string) error {
	m.mx.Lock()
	defer m.mx.Unlock()

	if _, ok := m.storage[name]; ok {
		return fmt.Errorf("variable already exists")
	}

	m.storage[name] = val

	return nil
}

func (m *Memory) Upsert(name string, val string) error {
	m.mx.Lock()
	defer m.mx.Unlock()

	m.storage[name] = val

	return nil
}

func (m *Memory) Get(name string) (string, error) {
	m.mx.RLock()
	defer m.mx.RUnlock()

	v, ok := m.storage[name]
	if !ok {
		return "", fmt.Errorf("variable not exists")
	}

	return v, nil
}

func (m *Memory) Delete(name string) error {
	m.mx.Lock()
	defer m.mx.Unlock()

	if _, ok := m.storage[name]; !ok {
		return fmt.Errorf("variable not exists")
	}

	delete(m.storage, name)

	return nil
}
