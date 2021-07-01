package manager

import (
	"fmt"
	"github.com/balerter/balerter/internal/config/scripts"
	"github.com/balerter/balerter/internal/config/scripts/file"
	"github.com/balerter/balerter/internal/config/scripts/folder"
	"github.com/balerter/balerter/internal/script/script"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"sort"
	"testing"
)

type scriptsProviderMock struct {
	mock.Mock
}

func (m *scriptsProviderMock) Get() ([]*script.Script, error) {
	args := m.Called()
	return args.Get(0).([]*script.Script), args.Error(1)
}

func TestNew(t *testing.T) {
	m := New()
	assert.IsType(t, &Manager{}, m)
}

func TestManager_Init(t *testing.T) {
	m := New()
	err := m.Init(&scripts.Scripts{
		UpdateInterval: 0,
		Folder:         []folder.Folder{{Name: "a", Path: "", Mask: ""}},
		File:           []file.File{{Name: "b", Filename: ""}},
		Postgres:       nil,
	})
	require.NoError(t, err)
	_, ok := m.providers["a"]
	assert.True(t, ok)
	_, ok = m.providers["b"]
	assert.True(t, ok)
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

	names := make([]string, 0)
	for _, s := range ss {
		names = append(names, s.Name)
	}
	sort.Strings(names)

	assert.Equal(t, "s11", names[0])
	assert.Equal(t, "s12", names[1])
	assert.Equal(t, "s21", names[2])
	assert.Equal(t, "s22", names[3])

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

func Test_removeTests(t *testing.T) {
	ss := []*script.Script{{Name: "1", IsTest: true}, {Name: "2", IsTest: false}, {Name: "3", IsTest: true}, {Name: "4", IsTest: false}}
	res := removeTests(ss)
	require.Equal(t, 2, len(res))
	assert.Equal(t, "2", res[0].Name)
	assert.Equal(t, "4", res[1].Name)

	// only tests
	ss = []*script.Script{{Name: "1", IsTest: true}, {Name: "2", IsTest: true}, {Name: "3", IsTest: true}, {Name: "4", IsTest: true}}
	res = removeTests(ss)
	require.Equal(t, 0, len(res))

	// empty
	ss = []*script.Script{}
	res = removeTests(ss)
	require.Equal(t, 0, len(res))

	// first test
	ss = []*script.Script{{Name: "1", IsTest: true}, {Name: "2", IsTest: false}, {Name: "3", IsTest: false}, {Name: "4", IsTest: false}}
	res = removeTests(ss)
	require.Equal(t, 3, len(res))
	assert.Equal(t, "2", res[0].Name)
	assert.Equal(t, "3", res[1].Name)
	assert.Equal(t, "4", res[2].Name)

	// last test
	ss = []*script.Script{{Name: "1", IsTest: false}, {Name: "2", IsTest: false}, {Name: "3", IsTest: false}, {Name: "4", IsTest: true}}
	res = removeTests(ss)
	require.Equal(t, 3, len(res))
	assert.Equal(t, "1", res[0].Name)
	assert.Equal(t, "2", res[1].Name)
	assert.Equal(t, "3", res[2].Name)

	// one test
	ss = []*script.Script{{Name: "1", IsTest: true}}
	res = removeTests(ss)
	require.Equal(t, 0, len(res))

	// one not test
	ss = []*script.Script{{Name: "1", IsTest: false}}
	res = removeTests(ss)
	require.Equal(t, 1, len(res))
	assert.Equal(t, "1", res[0].Name)
}

func TestManager_GetWithTests_error(t *testing.T) {
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

	ss, err := mgr.GetWithTests()

	p1.MethodCalled("Get")
	p2.MethodCalled("Get")

	require.Error(t, err)
	assert.Equal(t, "errorGet", err.Error())
	assert.Nil(t, ss)
}

func TestManager_GetWithTests(t *testing.T) {
	s11 := &script.Script{Name: "s11"}
	s21 := &script.Script{Name: "s21"}

	p1 := &scriptsProviderMock{}
	p2 := &scriptsProviderMock{}

	p1.On("Get").Return([]*script.Script{s11}, nil)
	p2.On("Get").Return([]*script.Script{s21}, nil)

	mgr := &Manager{
		providers: map[string]Provider{
			"p1": p1,
			"p2": p2,
		},
	}

	ss, err := mgr.GetWithTests()

	p1.MethodCalled("Get")
	p2.MethodCalled("Get")

	require.NoError(t, err)
	assert.Equal(t, 2, len(ss))

	for _, s := range ss {
		if s.Name != "s11" && s.Name != "s21" {
			t.Fatalf("unexpected s name %s", s.Name)
		}
	}
}
