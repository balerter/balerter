package corestorage

import (
	"github.com/balerter/balerter/internal/alert"
	//alertModule "github.com/balerter/balerter/internal/modules/alert"
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

func (m *AlertMock) Index(levels []alert.Level) (alert.Alerts, error) {
	args := m.Called(levels)
	a := args.Get(0)
	if a == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(alert.Alerts), args.Error(1)
}

func (m *AlertMock) Update(name string, level alert.Level) (*alert.Alert, bool, error) {
	args := m.Called(name, level)
	a := args.Get(0)
	if a == nil {
		return nil, args.Bool(1), args.Error(2)
	}
	return args.Get(0).(*alert.Alert), args.Bool(1), args.Error(2)
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
