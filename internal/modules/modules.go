package modules

import (
	"github.com/stretchr/testify/mock"
	lua "github.com/yuin/gopher-lua"
)

type Module interface {
	Name() string
	GetLoader() lua.LGFunction
	Stop() error
}

type ModuleMock struct {
	mock.Mock
}

func (m *ModuleMock) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *ModuleMock) GetLoader() lua.LGFunction {
	args := m.Called()
	return args.Get(0).(lua.LGFunction)
}

func (m *ModuleMock) Stop() error {
	args := m.Called()
	return args.Error(0)
}
