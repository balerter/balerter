// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package modules

import (
	"sync"

	"github.com/balerter/balerter/internal/script/script"
	lua "github.com/yuin/gopher-lua"
)

// ModuleMock is a mock implementation of Module.
//
// 	func TestSomethingThatUsesModule(t *testing.T) {
//
// 		// make and configure a mocked Module
// 		mockedModule := &ModuleMock{
// 			GetLoaderFunc: func(scriptMoqParam *script.Script) lua.LGFunction {
// 				panic("mock out the GetLoader method")
// 			},
// 			NameFunc: func() string {
// 				panic("mock out the Name method")
// 			},
// 			StopFunc: func() error {
// 				panic("mock out the Stop method")
// 			},
// 		}
//
// 		// use mockedModule in code that requires Module
// 		// and then make assertions.
//
// 	}
type ModuleMock struct {
	// GetLoaderFunc mocks the GetLoader method.
	GetLoaderFunc func(scriptMoqParam *script.Script) lua.LGFunction

	// NameFunc mocks the Name method.
	NameFunc func() string

	// StopFunc mocks the Stop method.
	StopFunc func() error

	// calls tracks calls to the methods.
	calls struct {
		// GetLoader holds details about calls to the GetLoader method.
		GetLoader []struct {
			// ScriptMoqParam is the scriptMoqParam argument value.
			ScriptMoqParam *script.Script
		}
		// Name holds details about calls to the Name method.
		Name []struct {
		}
		// Stop holds details about calls to the Stop method.
		Stop []struct {
		}
	}
	lockGetLoader sync.RWMutex
	lockName      sync.RWMutex
	lockStop      sync.RWMutex
}

// GetLoader calls GetLoaderFunc.
func (mock *ModuleMock) GetLoader(scriptMoqParam *script.Script) lua.LGFunction {
	if mock.GetLoaderFunc == nil {
		panic("ModuleMock.GetLoaderFunc: method is nil but Module.GetLoader was just called")
	}
	callInfo := struct {
		ScriptMoqParam *script.Script
	}{
		ScriptMoqParam: scriptMoqParam,
	}
	mock.lockGetLoader.Lock()
	mock.calls.GetLoader = append(mock.calls.GetLoader, callInfo)
	mock.lockGetLoader.Unlock()
	return mock.GetLoaderFunc(scriptMoqParam)
}

// GetLoaderCalls gets all the calls that were made to GetLoader.
// Check the length with:
//     len(mockedModule.GetLoaderCalls())
func (mock *ModuleMock) GetLoaderCalls() []struct {
	ScriptMoqParam *script.Script
} {
	var calls []struct {
		ScriptMoqParam *script.Script
	}
	mock.lockGetLoader.RLock()
	calls = mock.calls.GetLoader
	mock.lockGetLoader.RUnlock()
	return calls
}

// Name calls NameFunc.
func (mock *ModuleMock) Name() string {
	if mock.NameFunc == nil {
		panic("ModuleMock.NameFunc: method is nil but Module.Name was just called")
	}
	callInfo := struct {
	}{}
	mock.lockName.Lock()
	mock.calls.Name = append(mock.calls.Name, callInfo)
	mock.lockName.Unlock()
	return mock.NameFunc()
}

// NameCalls gets all the calls that were made to Name.
// Check the length with:
//     len(mockedModule.NameCalls())
func (mock *ModuleMock) NameCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockName.RLock()
	calls = mock.calls.Name
	mock.lockName.RUnlock()
	return calls
}

// Stop calls StopFunc.
func (mock *ModuleMock) Stop() error {
	if mock.StopFunc == nil {
		panic("ModuleMock.StopFunc: method is nil but Module.Stop was just called")
	}
	callInfo := struct {
	}{}
	mock.lockStop.Lock()
	mock.calls.Stop = append(mock.calls.Stop, callInfo)
	mock.lockStop.Unlock()
	return mock.StopFunc()
}

// StopCalls gets all the calls that were made to Stop.
// Check the length with:
//     len(mockedModule.StopCalls())
func (mock *ModuleMock) StopCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockStop.RLock()
	calls = mock.calls.Stop
	mock.lockStop.RUnlock()
	return calls
}
