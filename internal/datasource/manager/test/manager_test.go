package manager

import (
	"fmt"
	"github.com/balerter/balerter/internal/config/datasources"
	"github.com/balerter/balerter/internal/config/datasources/clickhouse"
	"github.com/balerter/balerter/internal/config/datasources/loki"
	"github.com/balerter/balerter/internal/config/datasources/mysql"
	"github.com/balerter/balerter/internal/config/datasources/postgres"
	"github.com/balerter/balerter/internal/config/datasources/prometheus"
	"github.com/balerter/balerter/internal/modules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestNew(t *testing.T) {
	m := New(zap.NewNop())

	assert.IsType(t, &Manager{}, m)
}

func TestClean(t *testing.T) {
	m := New(zap.NewNop())
	m1 := &modules.ModuleMock{}
	m2 := &modules.ModuleMock{}
	m.modules["m1"] = m1
	m.modules["m2"] = m2

	m1.On("Clean")
	m2.On("Clean")

	m.Clean()

	m1.AssertCalled(t, "Clean")
	m2.AssertCalled(t, "Clean")

	m1.AssertExpectations(t)
	m2.AssertExpectations(t)
}

func TestResult(t *testing.T) {
	m := New(zap.NewNop())
	m1 := &modules.ModuleMock{}
	m2 := &modules.ModuleMock{}
	m.modules["m1"] = m1
	m.modules["m2"] = m2

	r1 := modules.TestResult{ModuleName: "r1"}
	r2 := modules.TestResult{ModuleName: "r2"}

	m1.On("Result").Return([]modules.TestResult{r1}, nil)
	m2.On("Result").Return([]modules.TestResult{r2}, nil)

	res, err := m.Result()
	require.NoError(t, err)

	m1.AssertCalled(t, "Result")
	m2.AssertCalled(t, "Result")

	m1.AssertExpectations(t)
	m2.AssertExpectations(t)

	r1expect := modules.TestResult{ModuleName: "datasource.r1"}
	r2expect := modules.TestResult{ModuleName: "datasource.r2"}
	assert.Equal(t, 2, len(res))
	assert.Contains(t, res, r1expect)
	assert.Contains(t, res, r2expect)
}

func TestResult_error(t *testing.T) {
	m := New(zap.NewNop())
	m1 := &modules.ModuleMock{}
	m.modules["m1"] = m1

	m1.On("Result").Return(nil, fmt.Errorf("error1"))

	_, err := m.Result()
	require.Error(t, err)
	assert.Equal(t, "error1", err.Error())

	m1.AssertCalled(t, "Result")
	m1.AssertExpectations(t)
}

func TestInit(t *testing.T) {
	m := New(zap.NewNop())

	cfg := &datasources.DataSources{
		Clickhouse: []clickhouse.Clickhouse{{Name: "ch1"}},
		Prometheus: []prometheus.Prometheus{{Name: "prom1"}},
		Postgres:   []postgres.Postgres{{Name: "pg1"}},
		MySQL:      []mysql.Mysql{{Name: "mysql1"}},
		Loki:       []loki.Loki{{Name: "loki1"}},
	}

	err := m.Init(cfg)
	require.NoError(t, err)

	assert.Equal(t, 5, len(m.modules))

	mod, ok := m.modules["clickhouse.ch1"]
	assert.True(t, ok)
	require.NotNil(t, mod)
	assert.Equal(t, "clickhouse.ch1", mod.Name())

	mod, ok = m.modules["prometheus.prom1"]
	assert.True(t, ok)
	require.NotNil(t, mod)
	assert.Equal(t, "prometheus.prom1", mod.Name())

	mod, ok = m.modules["postgres.pg1"]
	assert.True(t, ok)
	require.NotNil(t, mod)
	assert.Equal(t, "postgres.pg1", mod.Name())

	mod, ok = m.modules["mysql.mysql1"]
	assert.True(t, ok)
	require.NotNil(t, mod)
	assert.Equal(t, "mysql.mysql1", mod.Name())

	mod, ok = m.modules["loki.loki1"]
	assert.True(t, ok)
	require.NotNil(t, mod)
	assert.Equal(t, "loki.loki1", mod.Name())
}

func TestGet(t *testing.T) {
	m := &Manager{
		modules: map[string]modules.ModuleTest{
			"m1": &modules.ModuleMock{},
			"m2": &modules.ModuleMock{},
		},
	}

	mods := m.Get()

	assert.Equal(t, 2, len(mods))
}
