package test

import (
	"github.com/balerter/balerter/internal/modules"
	"github.com/balerter/balerter/internal/script/script"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

type dsManager interface {
	GetMocks() []modules.ModuleTest
}

type Test struct {
	dsManager dsManager
	logger    *zap.Logger

	datasource map[string]modules.Module
}

func New(dsManager dsManager, logger *zap.Logger) *Test {
	t := &Test{
		dsManager: dsManager,
		logger:    logger,

		datasource: make(map[string]modules.Module),
	}

	return t
}

func (t *Test) Name() string {
	return "test"
}

func (t *Test) Stop() error {
	return nil
}

func (t *Test) getDatasource(s *script.Script) lua.LGFunction {
	return func(L *lua.LState) int {
		name := L.Get(1).String()

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
			//	"run": t.run(script.Name),
			"datasource": t.getDatasource(script),
		}

		for _, module := range t.dsManager.GetMocks() {
			t.datasource[module.Name()] = module
		}

		mod := L.SetFuncs(L.NewTable(), exports)

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
