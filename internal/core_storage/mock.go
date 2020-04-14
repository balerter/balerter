package core_storage

import (
	"github.com/balerter/balerter/internal/alert/alert"
	"github.com/stretchr/testify/mock"
)

type CoreStorageMock struct {
	name string
	mock.Mock
	alert *CoreStorageAlertMock
	kv    *CoreStorageKVMock
}

func NewMock(name string) *CoreStorageMock {
	m := &CoreStorageMock{
		name:  name,
		alert: &CoreStorageAlertMock{},
		kv:    &CoreStorageKVMock{},
	}

	return m
}

func (m *CoreStorageMock) KVMock() *CoreStorageKVMock {
	return m.kv
}

func (m *CoreStorageMock) AlertMock() *CoreStorageAlertMock {
	return m.alert
}

func (m *CoreStorageMock) KV() CoreStorageKV {
	return m.kv
}

func (m *CoreStorageMock) Stop() error {
	args := m.Called()
	return args.Error(0)
}

func (m *CoreStorageMock) Name() string {
	return m.name
}

func (m *CoreStorageMock) Alert() CoreStorageAlert {
	return m.alert
}

type CoreStorageAlertMock struct {
	mock.Mock
}

func (m *CoreStorageAlertMock) GetOrNew(string) (*alert.Alert, error) {
	args := m.Called()
	return args.Get(0).(*alert.Alert), args.Error(1)
}

func (m *CoreStorageAlertMock) All() ([]*alert.Alert, error) {
	args := m.Called()
	return args.Get(0).([]*alert.Alert), args.Error(1)
}

func (m *CoreStorageAlertMock) Release(a *alert.Alert) {
	args := m.Called(a)
	_ = args
}

func (m *CoreStorageAlertMock) Get(s string) (*alert.Alert, error) {
	args := m.Called(s)
	a := args.Get(0)
	e := args.Error(1)
	if a == nil {
		return nil, e
	}
	return a.(*alert.Alert), e
}

type CoreStorageKVMock struct {
	mock.Mock
}

func (m *CoreStorageKVMock) All() (map[string]string, error) {
	args := m.Called()
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *CoreStorageKVMock) Get(_ string) (string, error) {
	return "", nil
}

func (m *CoreStorageKVMock) Delete(_ string) error {
	return nil
}

func (m *CoreStorageKVMock) Put(_, _ string) error {
	return nil
}

func (m *CoreStorageKVMock) Upsert(_, _ string) error {
	return nil
}
