package corestorage

import (
	"github.com/balerter/balerter/internal/alert/alert"
	"github.com/stretchr/testify/mock"
)

type Mock struct {
	name string
	mock.Mock
	alert *AlertMock
	kv    *KVMock
}

func NewMock(name string) *Mock {
	m := &Mock{
		name:  name,
		alert: &AlertMock{},
		kv:    &KVMock{},
	}

	return m
}

func (m *Mock) KVMock() *KVMock {
	return m.kv
}

func (m *Mock) AlertMock() *AlertMock {
	return m.alert
}

func (m *Mock) KV() KV {
	return m.kv
}

func (m *Mock) Stop() error {
	args := m.Called()
	return args.Error(0)
}

func (m *Mock) Name() string {
	return m.name
}

func (m *Mock) Alert() Alert {
	return m.alert
}

type AlertMock struct {
	mock.Mock
}

func (m *AlertMock) GetOrNew(string) (*alert.Alert, error) {
	args := m.Called()
	a := args.Get(0)
	if a == nil {
		return nil, args.Error(1)
	}
	return a.(*alert.Alert), args.Error(1)
}

func (m *AlertMock) All() ([]*alert.Alert, error) {
	args := m.Called()
	return args.Get(0).([]*alert.Alert), args.Error(1)
}

func (m *AlertMock) Release(a *alert.Alert) error {
	args := m.Called(a)
	return args.Error(0)
}

func (m *AlertMock) Get(s string) (*alert.Alert, error) {
	args := m.Called(s)
	a := args.Get(0)
	e := args.Error(1)
	if a == nil {
		return nil, e
	}
	return a.(*alert.Alert), e
}

type KVMock struct {
	mock.Mock
}

func (m *KVMock) All() (map[string]string, error) {
	args := m.Called()
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *KVMock) Get(_ string) (string, error) {
	return "", nil
}

func (m *KVMock) Delete(_ string) error {
	return nil
}

func (m *KVMock) Put(_, _ string) error {
	return nil
}

func (m *KVMock) Upsert(_, _ string) error {
	return nil
}
