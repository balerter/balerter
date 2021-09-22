package runner

import (
	"context"
	"fmt"
	"github.com/balerter/balerter/internal/modules"
	"github.com/balerter/balerter/internal/script/script"
	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"testing"
	"time"
)

func TestNewRunner(t *testing.T) {
	r := New(0, nil, nil,
		nil, nil, "", nil, nil)
	assert.IsType(t, &Runner{}, r)
}

func TestRunner_filterScripts(t *testing.T) {
	type args struct {
		ss   []*script.Script
		name string
	}
	tests := []struct {
		name            string
		args            args
		wantScriptCount int
		wantScriptName  string
	}{
		{
			name: "no name",
			args: args{
				ss:   []*script.Script{{Name: "foo"}},
				name: "bar",
			},
			wantScriptCount: 0,
			wantScriptName:  "",
		},
		{
			name: "ok",
			args: args{
				ss:   []*script.Script{{Name: "foo"}, {Name: "bar"}, {Name: "baz"}},
				name: "bar",
			},
			wantScriptCount: 1,
			wantScriptName:  "bar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rnr := &Runner{}
			got := rnr.filterScripts(tt.args.ss, tt.args.name)
			if len(got) != tt.wantScriptCount {
				t.Errorf("unexpected script count = %d, want %d", len(got), tt.wantScriptCount)
			}
			if len(got) == 1 {
				if got[0].Name != tt.wantScriptName {
					t.Errorf("unexpected scipt name = %s, want %s", got[0].Name, tt.wantScriptName)
				}
			}
		})
	}
}

func TestRunner_watchJobs(t *testing.T) {
	r := &Runner{jobs: make(chan job, 1), logger: zap.NewNop()}

	go r.watchJobs()

	time.Sleep(time.Millisecond * 50)

	j := &jobMock{RunFunc: func() {}, NameFunc: func() string {
		return "job"
	}}
	r.jobs <- j

	time.Sleep(time.Millisecond * 50)

	assert.Equal(t, 1, len(j.RunCalls()))
	assert.Equal(t, 1, len(j.NameCalls()))
}

func TestRunner_Watch_once_error_get(t *testing.T) {
	sm := &scriptsManagerMock{
		GetFunc: func() ([]*script.Script, error) {
			return nil, fmt.Errorf("err1")
		},
	}

	core, logger := observer.New(zap.DebugLevel)

	rnr := &Runner{
		scriptsManager: sm,
		logger:         zap.New(core),
		cron:           cron.New(),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rnr.Watch(ctx, cancel, true)

	assert.Equal(t, 1, logger.FilterMessage("error get scripts").Len())
}

func TestRunner_Watch_once(t *testing.T) {
	sm := &scriptsManagerMock{
		GetFunc: func() ([]*script.Script, error) {
			return []*script.Script{{Name: "foo"}}, nil
		},
	}

	var called bool

	rnr := &Runner{
		scriptsManager: sm,
		logger:         zap.NewNop(),
		cron:           cron.New(),
		updateScriptsFunc: func(ctx context.Context, scripts []*script.Script, once bool) {
			called = true
			require.Equal(t, 1, len(scripts))
			assert.Equal(t, "foo", scripts[0].Name)
			assert.True(t, once)
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rnr.Watch(ctx, cancel, true)

	require.True(t, called)
}

func TestRunner_Watch_cliScript(t *testing.T) {
	sm := &scriptsManagerMock{
		GetFunc: func() ([]*script.Script, error) {
			return []*script.Script{{Name: "foo"}, {Name: "bar"}}, nil
		},
	}

	var called bool

	rnr := &Runner{
		cliScript:      "foo",
		scriptsManager: sm,
		logger:         zap.NewNop(),
		cron:           cron.New(),
		updateScriptsFunc: func(ctx context.Context, scripts []*script.Script, once bool) {
			called = true
			require.Equal(t, 1, len(scripts))
			require.Equal(t, "foo", scripts[0].Name)
			assert.True(t, once)
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rnr.Watch(ctx, cancel, true)

	require.True(t, called)
}

func TestRunner_Watch(t *testing.T) {
	sm := &scriptsManagerMock{
		GetFunc: func() ([]*script.Script, error) {
			return []*script.Script{{Name: "foo"}}, nil
		},
	}

	var called int

	rnr := &Runner{
		scriptsManager: sm,
		logger:         zap.NewNop(),
		cron:           cron.New(),
		updateScriptsFunc: func(ctx context.Context, scripts []*script.Script, once bool) {
			called++
			require.Equal(t, 1, len(scripts))
			assert.Equal(t, "foo", scripts[0].Name)
		},
		updateInterval: time.Millisecond * 300,
	}

	ctx, cancel := context.WithCancel(context.Background())

	go rnr.Watch(ctx, cancel, false)

	<-time.After(time.Millisecond * 500)
	cancel()

	require.Equal(t, 2, called)
}

func TestRunner_updateScripts_cancel_ctx(t *testing.T) {
	crn := cron.New()

	rnr := &Runner{
		cron:   crn,
		logger: zap.NewNop(),
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	rnr.updateScripts(ctx, []*script.Script{{}}, false)

	assert.Equal(t, 0, len(crn.Entries()))
}

func TestRunner_updateScripts_ignore(t *testing.T) {
	crn := cron.New()

	rnr := &Runner{
		cron:   crn,
		logger: zap.NewNop(),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rnr.updateScripts(ctx, []*script.Script{{Ignore: true}}, false)

	assert.Equal(t, 0, len(crn.Entries()))
}

func TestRunner_updateScripts_already_running(t *testing.T) {
	crn := cron.New()

	core, logger := observer.New(zap.DebugLevel)

	rnr := &Runner{
		cron:   crn,
		logger: zap.New(core),
		pool:   map[string]job{"9fc9328050ce91a3785b1042f9707c6ac5ae9fce": &Job{}},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rnr.updateScripts(ctx, []*script.Script{{Name: "foo"}}, false)

	assert.Equal(t, 0, len(crn.Entries()))
	assert.Equal(t, 1, logger.FilterMessage("script already scheduled").Len())
}

func TestRunner_updateScripts_once(t *testing.T) {
	crn := cron.New()

	coreModule := &modules.ModuleMock{
		NameFunc: func() string {
			return "m1"
		},
		GetLoaderFunc: func(_ modules.Job) lua.LGFunction {
			return func(ls *lua.LState) int {
				return 0
			}
		},
	}

	storagesModule := &modules.ModuleMock{
		NameFunc: func() string {
			return "m2"
		},
		GetLoaderFunc: func(_ modules.Job) lua.LGFunction {
			return func(ls *lua.LState) int {
				return 0
			}
		},
	}

	storagesManager := &storagesManagerMock{GetFunc: func() []modules.Module {
		return []modules.Module{storagesModule}
	}}

	dsModule := &modules.ModuleMock{
		NameFunc: func() string {
			return "m3"
		},
		GetLoaderFunc: func(_ modules.Job) lua.LGFunction {
			return func(ls *lua.LState) int {
				return 0
			}
		},
	}

	dsManager := &dsManagerMock{GetFunc: func() []modules.Module {
		return []modules.Module{dsModule}
	}}

	j := &jobMock{
		NameFunc: func() string {
			return "foo"
		},
		ScriptFunc: func() *script.Script {
			return &script.Script{}
		},
		SetLuaStateFunc: func(ls *lua.LState) {
			preload := ls.GetField(ls.GetField(ls.Get(lua.EnvironIndex), "package"), "preload")
			tbl, ok := preload.(*lua.LTable)
			require.True(t, ok)

			assert.Equal(t, lua.LTFunction, tbl.RawGetString("m1").Type())
			assert.Equal(t, lua.LTFunction, tbl.RawGetString("storage.m2").Type())
			assert.Equal(t, lua.LTFunction, tbl.RawGetString("datasource.m3").Type())
			assert.Equal(t, lua.LTFunction, tbl.RawGetString("api").Type())
		},
		RunFunc: func() {

		},
	}

	rnr := &Runner{
		coreModules:     []modules.Module{coreModule},
		storagesManager: storagesManager,
		dsManager:       dsManager,

		cron:   crn,
		logger: zap.NewNop(),
		pool:   map[string]job{},
		newJobFunc: func(s *script.Script, logger *zap.Logger) job {
			return j
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rnr.updateScripts(ctx, []*script.Script{{Name: "foo"}}, true)

	assert.Equal(t, 1, len(j.SetLuaStateCalls()))
	assert.Equal(t, 1, len(j.NameCalls()))
	assert.Equal(t, 1, len(j.RunCalls()))
	assert.Equal(t, 1, len(coreModule.NameCalls()))
	assert.Equal(t, 1, len(coreModule.GetLoaderCalls()))
	assert.Equal(t, 1, len(storagesModule.NameCalls()))
	assert.Equal(t, 1, len(storagesModule.GetLoaderCalls()))
	assert.Equal(t, 1, len(storagesManager.GetCalls()))
	assert.Equal(t, 1, len(dsModule.NameCalls()))
	assert.Equal(t, 1, len(dsModule.GetLoaderCalls()))
	assert.Equal(t, 1, len(dsManager.GetCalls()))
}

func TestRunner_updateScripts(t *testing.T) {
	coreModule := &modules.ModuleMock{
		NameFunc: func() string {
			return "m1"
		},
		GetLoaderFunc: func(_ modules.Job) lua.LGFunction {
			return func(ls *lua.LState) int {
				return 0
			}
		},
	}

	storagesModule := &modules.ModuleMock{
		NameFunc: func() string {
			return "m2"
		},
		GetLoaderFunc: func(_ modules.Job) lua.LGFunction {
			return func(ls *lua.LState) int {
				return 0
			}
		},
	}

	storagesManager := &storagesManagerMock{GetFunc: func() []modules.Module {
		return []modules.Module{storagesModule}
	}}

	dsModule := &modules.ModuleMock{
		NameFunc: func() string {
			return "m3"
		},
		GetLoaderFunc: func(_ modules.Job) lua.LGFunction {
			return func(ls *lua.LState) int {
				return 0
			}
		},
	}

	dsManager := &dsManagerMock{GetFunc: func() []modules.Module {
		return []modules.Module{dsModule}
	}}

	j := &jobMock{
		NameFunc: func() string {
			return "foo"
		},
		ScriptFunc: func() *script.Script {
			return &script.Script{}
		},
		SetLuaStateFunc: func(ls *lua.LState) {
			preload := ls.GetField(ls.GetField(ls.Get(lua.EnvironIndex), "package"), "preload")
			tbl, ok := preload.(*lua.LTable)
			require.True(t, ok)

			assert.Equal(t, lua.LTFunction, tbl.RawGetString("m1").Type())
			assert.Equal(t, lua.LTFunction, tbl.RawGetString("storage.m2").Type())
			assert.Equal(t, lua.LTFunction, tbl.RawGetString("datasource.m3").Type())
			assert.Equal(t, lua.LTFunction, tbl.RawGetString("api").Type())
		},
		RunFunc: func() {

		},
		SetEntryIDFunc: func(_ cron.EntryID) {

		},
	}

	crn := cron.New()
	s1EntryID, err := crn.AddFunc("1 1 1 1 1", func() {})
	require.NoError(t, err)

	// Hash 1dab309881b4ce59a311b1fe1f1bb33306b43c28
	s1Hash := "1dab309881b4ce59a311b1fe1f1bb33306b43c28"
	s1 := &script.Script{
		Name: "bar",
	}
	j1 := &jobMock{
		ScriptFunc: func() *script.Script {
			return s1
		},
		StopFunc: func() {

		},
		EntryIDFunc: func() cron.EntryID {
			return s1EntryID
		},
	}

	// Hash 1dab309881b4ce59a311b1fe1f1bb33306b43c28
	s2 := &script.Script{
		Name:      "foo",
		CronValue: "1 1 1 1 1",
	}

	rnr := &Runner{
		coreModules:     []modules.Module{coreModule},
		storagesManager: storagesManager,
		dsManager:       dsManager,

		cron:   crn,
		logger: zap.NewNop(),
		pool:   map[string]job{s1Hash: j1},
		newJobFunc: func(s *script.Script, logger *zap.Logger) job {
			return j
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rnr.updateScripts(ctx, []*script.Script{s2}, false)

	assert.Equal(t, 1, len(j.SetLuaStateCalls()))
	assert.Equal(t, 1, len(j.NameCalls()))
	assert.Equal(t, 1, len(j.SetEntryIDCalls()))
	assert.Equal(t, 1, len(coreModule.NameCalls()))
	assert.Equal(t, 1, len(coreModule.GetLoaderCalls()))
	assert.Equal(t, 1, len(storagesModule.NameCalls()))
	assert.Equal(t, 1, len(storagesModule.GetLoaderCalls()))
	assert.Equal(t, 1, len(storagesManager.GetCalls()))
	assert.Equal(t, 1, len(dsModule.NameCalls()))
	assert.Equal(t, 1, len(dsModule.GetLoaderCalls()))
	assert.Equal(t, 1, len(dsManager.GetCalls()))

	assert.Equal(t, 1, len(rnr.pool))
}

func TestRunner_Stop(t *testing.T) {
	crn := cron.New()
	entryID, err := crn.AddFunc("1 1 1 1 1", func() {})
	require.NoError(t, err)

	j := &jobMock{
		ScriptFunc: func() *script.Script {
			return &script.Script{
				Name: "s1",
			}
		},
		StopFunc: func() {

		},
		EntryIDFunc: func() cron.EntryID {
			return entryID
		},
	}

	rnr := &Runner{
		jobs:   make(chan job),
		cron:   crn,
		logger: zap.NewNop(),
		pool:   map[string]job{"1": j},
	}

	rnr.Stop()

	assert.Equal(t, 1, len(j.StopCalls()))
	assert.Equal(t, 1, len(j.EntryIDCalls()))
	assert.Equal(t, 1, len(j.ScriptCalls()))
	assert.Equal(t, 0, len(crn.Entries()))
}
