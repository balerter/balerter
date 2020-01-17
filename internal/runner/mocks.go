package runner

import (
	"github.com/balerter/balerter/internal/script/script"
	"github.com/stretchr/testify/mock"
)

type scriptManagerMock struct {
	mock.Mock
}

func (m *scriptManagerMock) Get() ([]*script.Script, error) {
	args := m.Called()
	return args.Get(0).([]*script.Script), args.Error(1)
}
