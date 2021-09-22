// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package modules

import (
	"sync"
	"time"

	"github.com/balerter/balerter/internal/script/script"
)

// JobMock is a mock implementation of Job.
//
// 	func TestSomethingThatUsesJob(t *testing.T) {
//
// 		// make and configure a mocked Job
// 		mockedJob := &JobMock{
// 			GetPriorExecutionTimeFunc: func() time.Duration {
// 				panic("mock out the GetPriorExecutionTime method")
// 			},
// 			ScriptFunc: func() *script.Script {
// 				panic("mock out the Script method")
// 			},
// 		}
//
// 		// use mockedJob in code that requires Job
// 		// and then make assertions.
//
// 	}
type JobMock struct {
	// GetPriorExecutionTimeFunc mocks the GetPriorExecutionTime method.
	GetPriorExecutionTimeFunc func() time.Duration

	// ScriptFunc mocks the Script method.
	ScriptFunc func() *script.Script

	// calls tracks calls to the methods.
	calls struct {
		// GetPriorExecutionTime holds details about calls to the GetPriorExecutionTime method.
		GetPriorExecutionTime []struct {
		}
		// Script holds details about calls to the Script method.
		Script []struct {
		}
	}
	lockGetPriorExecutionTime sync.RWMutex
	lockScript                sync.RWMutex
}

// GetPriorExecutionTime calls GetPriorExecutionTimeFunc.
func (mock *JobMock) GetPriorExecutionTime() time.Duration {
	if mock.GetPriorExecutionTimeFunc == nil {
		panic("JobMock.GetPriorExecutionTimeFunc: method is nil but Job.GetPriorExecutionTime was just called")
	}
	callInfo := struct {
	}{}
	mock.lockGetPriorExecutionTime.Lock()
	mock.calls.GetPriorExecutionTime = append(mock.calls.GetPriorExecutionTime, callInfo)
	mock.lockGetPriorExecutionTime.Unlock()
	return mock.GetPriorExecutionTimeFunc()
}

// GetPriorExecutionTimeCalls gets all the calls that were made to GetPriorExecutionTime.
// Check the length with:
//     len(mockedJob.GetPriorExecutionTimeCalls())
func (mock *JobMock) GetPriorExecutionTimeCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockGetPriorExecutionTime.RLock()
	calls = mock.calls.GetPriorExecutionTime
	mock.lockGetPriorExecutionTime.RUnlock()
	return calls
}

// Script calls ScriptFunc.
func (mock *JobMock) Script() *script.Script {
	if mock.ScriptFunc == nil {
		panic("JobMock.ScriptFunc: method is nil but Job.Script was just called")
	}
	callInfo := struct {
	}{}
	mock.lockScript.Lock()
	mock.calls.Script = append(mock.calls.Script, callInfo)
	mock.lockScript.Unlock()
	return mock.ScriptFunc()
}

// ScriptCalls gets all the calls that were made to Script.
// Check the length with:
//     len(mockedJob.ScriptCalls())
func (mock *JobMock) ScriptCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockScript.RLock()
	calls = mock.calls.Script
	mock.lockScript.RUnlock()
	return calls
}
