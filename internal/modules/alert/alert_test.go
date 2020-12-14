package alert

import (
	alert2 "github.com/balerter/balerter/internal/alert"
	"github.com/stretchr/testify/mock"
)

type managerMock struct {
	mock.Mock
}

func (m *managerMock) Call(alertName string, alertLevel alert2.Level, text string, options *alert2.Options) error {
	args := m.Called(alertName, alertLevel, text, options)
	return args.Error(0)
}

func (m *managerMock) Get(name string) (*alert2.Alert, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*alert2.Alert), args.Error(1)
}
