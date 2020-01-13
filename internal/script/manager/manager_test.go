package manager

import (
	"fmt"
	"github.com/balerter/balerter/internal/script/script"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

type scriptsProviderMock struct {
	mock.Mock
}

func (m *scriptsProviderMock) Get() ([]*script.Script, error) {
	args := m.Called()
	return args.Get(0).([]*script.Script), args.Error(1)
}

func TestManager_Get_Success(t *testing.T) {
	s11 := &script.Script{Name: "s11"}
	s12 := &script.Script{Name: "s12"}
	s21 := &script.Script{Name: "s21"}
	s22 := &script.Script{Name: "s22"}

	p1 := &scriptsProviderMock{}
	p2 := &scriptsProviderMock{}

	p1.On("Get").Return([]*script.Script{s11, s12}, nil)
	p2.On("Get").Return([]*script.Script{s21, s22}, nil)

	mgr := &Manager{
		providers: map[string]Provider{
			"p1": p1,
			"p2": p2,
		},
	}

	ss, err := mgr.Get()

	p1.MethodCalled("Get")
	p2.MethodCalled("Get")

	require.NoError(t, err)
	require.Equal(t, 4, len(ss))
	assert.Equal(t, "s11", ss[0].Name)
	assert.Equal(t, "s12", ss[1].Name)
	assert.Equal(t, "s21", ss[2].Name)
	assert.Equal(t, "s22", ss[3].Name)

	p1.AssertExpectations(t)
	p2.AssertExpectations(t)
}

func TestManager_Get_Fail(t *testing.T) {
	s11 := &script.Script{Name: "s11"}
	s21 := &script.Script{Name: "s21"}

	p1 := &scriptsProviderMock{}
	p2 := &scriptsProviderMock{}

	e := fmt.Errorf("errorGet")

	p1.On("Get").Return([]*script.Script{s11}, nil)
	p2.On("Get").Return([]*script.Script{s21}, e)

	mgr := &Manager{
		providers: map[string]Provider{
			"p1": p1,
			"p2": p2,
		},
	}

	ss, err := mgr.Get()

	p1.MethodCalled("Get")
	p2.MethodCalled("Get")

	require.Error(t, err)
	assert.Equal(t, "errorGet", err.Error())
	assert.Nil(t, ss)
}
