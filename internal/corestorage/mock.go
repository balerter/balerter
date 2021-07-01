package corestorage

import (
	"github.com/balerter/balerter/internal/alert"
	"github.com/stretchr/testify/mock"
)

// Mock is the mock for the CoreStorage
type Mock struct {
	name string
	mock.Mock
	alert *AlertMock
	kv    *KVMock
}

// NewMock creates new Mock
func NewMock(name string) *Mock {
	m := &Mock{
		name:  name,
		alert: &AlertMock{},
		kv:    &KVMock{},
	}

	return m
}

// KVMock returns KV storage from the Core
func (m *Mock) KVMock() *KVMock {
	return m.kv
}

// AlertMock returns Alert storage from the Core
func (m *Mock) AlertMock() *AlertMock {
	return m.alert
}

// KV returns KV storage
func (m *Mock) KV() KV {
	return m.kv
}

// Stop the mock
func (m *Mock) Stop() error {
	args := m.Called()
	return args.Error(0)
}

// Name returns storage name
func (m *Mock) Name() string {
	return m.name
}

// Alert return Alert storage
func (m *Mock) Alert() Alert {
	return m.alert
}

// AlertMock represents mock for Alert storage
type AlertMock struct {
	mock.Mock
}

// Get returns alert from the storage
func (m *AlertMock) Get(name string) (*alert.Alert, error) {
	args := m.Called(name)
	a := args.Get(0)
	if a == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*alert.Alert), args.Error(1)
}

// Index returns alerts list
func (m *AlertMock) Index(levels []alert.Level) (alert.Alerts, error) {
	args := m.Called(levels)
	a := args.Get(0)
	if a == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(alert.Alerts), args.Error(1)
}

// Update the alert
func (m *AlertMock) Update(name string, level alert.Level) (*alert.Alert, bool, error) {
	args := m.Called(name, level)
	a := args.Get(0)
	if a == nil {
		return nil, args.Bool(1), args.Error(2) //nolint:gomnd // arg position
	}
	return args.Get(0).(*alert.Alert), args.Bool(1), args.Error(2) //nolint:gomnd // arg position
}

// KVMock represents KV mock storage
type KVMock struct {
	mock.Mock
}

// All returns all KV pairs
func (m *KVMock) All() (map[string]string, error) {
	args := m.Called()
	return args.Get(0).(map[string]string), args.Error(1)
}

// Get returns KV pair by name
func (m *KVMock) Get(n string) (string, error) {
	a := m.Called(n)
	return a.String(0), a.Error(1)
}

// Delete KV pair from the storage
func (m *KVMock) Delete(v string) error {
	a := m.Called(v)
	return a.Error(0)
}

// Put KV pair to the storage
func (m *KVMock) Put(v1, v2 string) error {
	a := m.Called(v1, v2)
	return a.Error(0)
}

// Upsert KV pair in the storage
func (m *KVMock) Upsert(v1, v2 string) error {
	a := m.Called(v1, v2)
	return a.Error(0)
}
