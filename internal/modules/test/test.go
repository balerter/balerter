package test

import (
	"github.com/balerter/balerter/internal/mock"
	"github.com/balerter/balerter/internal/modules"
	"github.com/balerter/balerter/internal/script/script"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"strings"
)

type modulesManager interface {
	Get() []modules.ModuleTest
}

type Test struct {
	dsManager      modulesManager
	storageManager modulesManager
	alertMgr       *mock.ModuleMock
	logModule      *mock.ModuleMock
	logger         *zap.Logger

	datasource map[string]modules.ModuleTest
	storage    map[string]modules.ModuleTest
}

func New(dsManager modulesManager, storageManager modulesManager, alertMgr, logModule *mock.ModuleMock, logger *zap.Logger) *Test {
	t := &Test{
		dsManager:      dsManager,
		storageManager: storageManager,
		alertMgr:       alertMgr,
		logModule:      logModule,
		logger:         logger,

		datasource: make(map[string]modules.ModuleTest),
		storage:    make(map[string]modules.ModuleTest),
	}

	for _, module := range t.dsManager.Get() {
		t.datasource[module.Name()] = module
	}

	for _, module := range t.storageManager.Get() {
		t.storage[module.Name()] = module
	}

	return t
}

func (t *Test) Name() string {
	return "test"
}

func (t *Test) Stop() error {
	return nil
}

func (t *Test) getStorage(s *script.Script) lua.LGFunction {
	return func(L *lua.LState) int {
		nameL := L.Get(1)
		if nameL.Type() == lua.LTNil {
			t.logger.Error("module test.storage should have 1 argument")
			return 0
		}

		name := strings.TrimSpace(nameL.String())
		if name == "" {
			t.logger.Error("module test.storage should have 1 not empty argument")
			return 0
		}

		m, ok := t.storage[name]
		if !ok {
			t.logger.Error("storage not found", zap.String("name", name))
			return 0
		}

		m.GetLoader(s)(L)

		return 1
	}
}

func (t *Test) getAlertMgr(s *script.Script) lua.LGFunction {
	return func(L *lua.LState) int {
		return t.alertMgr.GetLoader(s)(L)
	}
}

func (t *Test) getLogModule(s *script.Script) lua.LGFunction {
	return func(L *lua.LState) int {
		return t.logModule.GetLoader(s)(L)
	}
}

func (t *Test) getDatasource(s *script.Script) lua.LGFunction {
	return func(L *lua.LState) int {
		nameL := L.Get(1)
		if nameL.Type() == lua.LTNil {
			t.logger.Error("module test.datasource should have 1 argument")
			return 0
		}

		name := strings.TrimSpace(nameL.String())
		if name == "" {
			t.logger.Error("module test.datasource should have 1 not empty argument")
			return 0
		}

		m, ok := t.datasource[name]
		if !ok {
			t.logger.Error("datasource not found", zap.String("name", name))
			return 0
		}

		m.GetLoader(s)(L)

		return 1
	}
}

func (t *Test) GetLoader(script *script.Script) lua.LGFunction {
	return func(L *lua.LState) int {
		var exports = map[string]lua.LGFunction{
			"datasource": t.getDatasource(script),
			"storage":    t.getStorage(script),
			"alert":      t.getAlertMgr(script),
			"log":        t.getLogModule(script),
			//	"run": t.run(script.Name),
		}

		mod := L.SetFuncs(L.NewTable(), exports)

		mod.RawSetString("AnyValue", lua.LString(mock.AnyValue))

		L.Push(mod)
		return 1

	}
}

//func (t *Test) run(_ string) lua.LGFunction {
//	return func(L *lua.LState) int {
//
//		name := L.Get(1)
//		if name.Type() == lua.LTNil {
//			t.logger.Error("error get test name")
//			L.Push(lua.LString("error get test name"))
//			return 1
//		}
//
//		lFunc := L.Get(2)
//		if lFunc.Type() != lua.LTFunction {
//			t.logger.Error("error get test function")
//			L.Push(lua.LString("error get test function"))
//			return 1
//		}
//
//		//t.logger.Debug()
//
//		f := lFunc.(*lua.LFunction)
//
//		fState := lua.NewState()
//		fState.Push(f)
//
//		inTst := &lua.LTable{}
//		inTst.RawSet(lua.LString("assert"), fState.NewFunction(t.aa))
//
//		fState.Push(inTst)
//
//		err := fState.PCall(1, lua.MultRet, nil)
//		if err != nil {
//			t.logger.Error("error run test function", zap.Error(err))
//			L.Push(lua.LString("error run test function: " + err.Error()))
//			return 1
//		}
//
//		return 0
//	}
//}
