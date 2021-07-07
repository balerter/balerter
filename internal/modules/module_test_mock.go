// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package modules

import (
	"sync"

	"github.com/balerter/balerter/internal/script/script"
	lua "github.com/yuin/gopher-lua"
)

// ModuleTestMock is a mock implementation of ModuleTest.
//
// 	func TestSomethingThatUsesModuleTest(t *testing.T) {
//
// 		// make and configure a mocked ModuleTest
// 		mockedModuleTest := &ModuleTestMock{
// 			CleanFunc: func()  {
// 				panic("mock out the Clean method")
// 			},
// 			GetLoaderFunc: func(scriptMoqParam *script.Script) lua.LGFunction {
// 				panic("mock out the GetLoader method")
// 			},
// 			NameFunc: func() string {
// 				panic("mock out the Name method")
// 			},
// 			ResultFunc: func() ([]TestResult, error) {
// 				panic("mock out the Result method")
// 			},
// 		}
//
// 		// use mockedModuleTest in code that requires ModuleTest
// 		// and then make assertions.
//
// 	}
type ModuleTestMock struct {
	// CleanFunc mocks the Clean method.
	CleanFunc func()

	// GetLoaderFunc mocks the GetLoader method.
	GetLoaderFunc func(scriptMoqParam *script.Script) lua.LGFunction

	// NameFunc mocks the Name method.
	NameFunc func() string

	// ResultFunc mocks the Result method.
	ResultFunc func() ([]TestResult, error)

	// calls tracks calls to the methods.
	calls struct {
		// Clean holds details about calls to the Clean method.
		Clean []struct {
		}
		// GetLoader holds details about calls to the GetLoader method.
		GetLoader []struct {
			// ScriptMoqParam is the scriptMoqParam argument value.
			ScriptMoqParam *script.Script
		}
		// Name holds details about calls to the Name method.
		Name []struct {
		}
		// Result holds details about calls to the Result method.
		Result []struct {
		}
	}
	lockClean     sync.RWMutex
	lockGetLoader sync.RWMutex
	lockName      sync.RWMutex
	lockResult    sync.RWMutex
}

// Clean calls CleanFunc.
func (mock *ModuleTestMock) Clean() {
	if mock.CleanFunc == nil {
		panic("ModuleTestMock.CleanFunc: method is nil but ModuleTest.Clean was just called")
	}
	callInfo := struct {
	}{}
	mock.lockClean.Lock()
	mock.calls.Clean = append(mock.calls.Clean, callInfo)
	mock.lockClean.Unlock()
	mock.CleanFunc()
}

// CleanCalls gets all the calls that were made to Clean.
// Check the length with:
//     len(mockedModuleTest.CleanCalls())
func (mock *ModuleTestMock) CleanCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockClean.RLock()
	calls = mock.calls.Clean
	mock.lockClean.RUnlock()
	return calls
}

// GetLoader calls GetLoaderFunc.
func (mock *ModuleTestMock) GetLoader(scriptMoqParam *script.Script) lua.LGFunction {
	if mock.GetLoaderFunc == nil {
		panic("ModuleTestMock.GetLoaderFunc: method is nil but ModuleTest.GetLoader was just called")
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
//     len(mockedModuleTest.GetLoaderCalls())
func (mock *ModuleTestMock) GetLoaderCalls() []struct {
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
func (mock *ModuleTestMock) Name() string {
	if mock.NameFunc == nil {
		panic("ModuleTestMock.NameFunc: method is nil but ModuleTest.Name was just called")
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
//     len(mockedModuleTest.NameCalls())
func (mock *ModuleTestMock) NameCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockName.RLock()
	calls = mock.calls.Name
	mock.lockName.RUnlock()
	return calls
}

// Result calls ResultFunc.
func (mock *ModuleTestMock) Result() ([]TestResult, error) {
	if mock.ResultFunc == nil {
		panic("ModuleTestMock.ResultFunc: method is nil but ModuleTest.Result was just called")
	}
	callInfo := struct {
	}{}
	mock.lockResult.Lock()
	mock.calls.Result = append(mock.calls.Result, callInfo)
	mock.lockResult.Unlock()
	return mock.ResultFunc()
}

// ResultCalls gets all the calls that were made to Result.
// Check the length with:
//     len(mockedModuleTest.ResultCalls())
func (mock *ModuleTestMock) ResultCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockResult.RLock()
	calls = mock.calls.Result
	mock.lockResult.RUnlock()
	return calls
}
