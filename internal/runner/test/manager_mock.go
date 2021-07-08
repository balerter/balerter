// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package test

import (
	"sync"

	"github.com/balerter/balerter/internal/modules"
)

// managerMock is a mock implementation of manager.
//
// 	func TestSomethingThatUsesmanager(t *testing.T) {
//
// 		// make and configure a mocked manager
// 		mockedmanager := &managerMock{
// 			CleanFunc: func()  {
// 				panic("mock out the Clean method")
// 			},
// 			GetFunc: func() []modules.ModuleTest {
// 				panic("mock out the Get method")
// 			},
// 			ResultFunc: func() ([]modules.TestResult, error) {
// 				panic("mock out the Result method")
// 			},
// 		}
//
// 		// use mockedmanager in code that requires manager
// 		// and then make assertions.
//
// 	}
type managerMock struct {
	// CleanFunc mocks the Clean method.
	CleanFunc func()

	// GetFunc mocks the Get method.
	GetFunc func() []modules.ModuleTest

	// ResultFunc mocks the Result method.
	ResultFunc func() ([]modules.TestResult, error)

	// calls tracks calls to the methods.
	calls struct {
		// Clean holds details about calls to the Clean method.
		Clean []struct {
		}
		// Get holds details about calls to the Get method.
		Get []struct {
		}
		// Result holds details about calls to the Result method.
		Result []struct {
		}
	}
	lockClean  sync.RWMutex
	lockGet    sync.RWMutex
	lockResult sync.RWMutex
}

// Clean calls CleanFunc.
func (mock *managerMock) Clean() {
	if mock.CleanFunc == nil {
		panic("managerMock.CleanFunc: method is nil but manager.Clean was just called")
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
//     len(mockedmanager.CleanCalls())
func (mock *managerMock) CleanCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockClean.RLock()
	calls = mock.calls.Clean
	mock.lockClean.RUnlock()
	return calls
}

// Get calls GetFunc.
func (mock *managerMock) Get() []modules.ModuleTest {
	if mock.GetFunc == nil {
		panic("managerMock.GetFunc: method is nil but manager.Get was just called")
	}
	callInfo := struct {
	}{}
	mock.lockGet.Lock()
	mock.calls.Get = append(mock.calls.Get, callInfo)
	mock.lockGet.Unlock()
	return mock.GetFunc()
}

// GetCalls gets all the calls that were made to Get.
// Check the length with:
//     len(mockedmanager.GetCalls())
func (mock *managerMock) GetCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockGet.RLock()
	calls = mock.calls.Get
	mock.lockGet.RUnlock()
	return calls
}

// Result calls ResultFunc.
func (mock *managerMock) Result() ([]modules.TestResult, error) {
	if mock.ResultFunc == nil {
		panic("managerMock.ResultFunc: method is nil but manager.Result was just called")
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
//     len(mockedmanager.ResultCalls())
func (mock *managerMock) ResultCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockResult.RLock()
	calls = mock.calls.Result
	mock.lockResult.RUnlock()
	return calls
}