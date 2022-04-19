// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package corestorage

import (
	"sync"
)

// KVMock is a mock implementation of KV.
//
// 	func TestSomethingThatUsesKV(t *testing.T) {
//
// 		// make and configure a mocked KV
// 		mockedKV := &KVMock{
// 			AllFunc: func() (map[string]string, error) {
// 				panic("mock out the All method")
// 			},
// 			DeleteFunc: func(s string) error {
// 				panic("mock out the Delete method")
// 			},
// 			GetFunc: func(s string) (string, error) {
// 				panic("mock out the Get method")
// 			},
// 			PutFunc: func(s1 string, s2 string) error {
// 				panic("mock out the Put method")
// 			},
// 			UpsertFunc: func(s1 string, s2 string) error {
// 				panic("mock out the Upsert method")
// 			},
// 		}
//
// 		// use mockedKV in code that requires KV
// 		// and then make assertions.
//
// 	}
type KVMock struct {
	// AllFunc mocks the All method.
	AllFunc func() (map[string]string, error)

	// DeleteFunc mocks the Delete method.
	DeleteFunc func(s string) error

	// GetFunc mocks the Get method.
	GetFunc func(s string) (string, error)

	// PutFunc mocks the Put method.
	PutFunc func(s1 string, s2 string) error

	// UpsertFunc mocks the Upsert method.
	UpsertFunc func(s1 string, s2 string) error

	// calls tracks calls to the methods.
	calls struct {
		// All holds details about calls to the All method.
		All []struct {
		}
		// Delete holds details about calls to the Delete method.
		Delete []struct {
			// S is the s argument value.
			S string
		}
		// Get holds details about calls to the Get method.
		Get []struct {
			// S is the s argument value.
			S string
		}
		// Put holds details about calls to the Put method.
		Put []struct {
			// S1 is the s1 argument value.
			S1 string
			// S2 is the s2 argument value.
			S2 string
		}
		// Upsert holds details about calls to the Upsert method.
		Upsert []struct {
			// S1 is the s1 argument value.
			S1 string
			// S2 is the s2 argument value.
			S2 string
		}
	}
	lockAll    sync.RWMutex
	lockDelete sync.RWMutex
	lockGet    sync.RWMutex
	lockPut    sync.RWMutex
	lockUpsert sync.RWMutex
}

// All calls AllFunc.
func (mock *KVMock) All() (map[string]string, error) {
	if mock.AllFunc == nil {
		panic("KVMock.AllFunc: method is nil but KV.All was just called")
	}
	callInfo := struct {
	}{}
	mock.lockAll.Lock()
	mock.calls.All = append(mock.calls.All, callInfo)
	mock.lockAll.Unlock()
	return mock.AllFunc()
}

// AllCalls gets all the calls that were made to All.
// Check the length with:
//     len(mockedKV.AllCalls())
func (mock *KVMock) AllCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockAll.RLock()
	calls = mock.calls.All
	mock.lockAll.RUnlock()
	return calls
}

// Delete calls DeleteFunc.
func (mock *KVMock) Delete(s string) error {
	if mock.DeleteFunc == nil {
		panic("KVMock.DeleteFunc: method is nil but KV.Delete was just called")
	}
	callInfo := struct {
		S string
	}{
		S: s,
	}
	mock.lockDelete.Lock()
	mock.calls.Delete = append(mock.calls.Delete, callInfo)
	mock.lockDelete.Unlock()
	return mock.DeleteFunc(s)
}

// DeleteCalls gets all the calls that were made to Delete.
// Check the length with:
//     len(mockedKV.DeleteCalls())
func (mock *KVMock) DeleteCalls() []struct {
	S string
} {
	var calls []struct {
		S string
	}
	mock.lockDelete.RLock()
	calls = mock.calls.Delete
	mock.lockDelete.RUnlock()
	return calls
}

// Get calls GetFunc.
func (mock *KVMock) Get(s string) (string, error) {
	if mock.GetFunc == nil {
		panic("KVMock.GetFunc: method is nil but KV.Get was just called")
	}
	callInfo := struct {
		S string
	}{
		S: s,
	}
	mock.lockGet.Lock()
	mock.calls.Get = append(mock.calls.Get, callInfo)
	mock.lockGet.Unlock()
	return mock.GetFunc(s)
}

// GetCalls gets all the calls that were made to Get.
// Check the length with:
//     len(mockedKV.GetCalls())
func (mock *KVMock) GetCalls() []struct {
	S string
} {
	var calls []struct {
		S string
	}
	mock.lockGet.RLock()
	calls = mock.calls.Get
	mock.lockGet.RUnlock()
	return calls
}

// Put calls PutFunc.
func (mock *KVMock) Put(s1 string, s2 string) error {
	if mock.PutFunc == nil {
		panic("KVMock.PutFunc: method is nil but KV.Put was just called")
	}
	callInfo := struct {
		S1 string
		S2 string
	}{
		S1: s1,
		S2: s2,
	}
	mock.lockPut.Lock()
	mock.calls.Put = append(mock.calls.Put, callInfo)
	mock.lockPut.Unlock()
	return mock.PutFunc(s1, s2)
}

// PutCalls gets all the calls that were made to Put.
// Check the length with:
//     len(mockedKV.PutCalls())
func (mock *KVMock) PutCalls() []struct {
	S1 string
	S2 string
} {
	var calls []struct {
		S1 string
		S2 string
	}
	mock.lockPut.RLock()
	calls = mock.calls.Put
	mock.lockPut.RUnlock()
	return calls
}

// Upsert calls UpsertFunc.
func (mock *KVMock) Upsert(s1 string, s2 string) error {
	if mock.UpsertFunc == nil {
		panic("KVMock.UpsertFunc: method is nil but KV.Upsert was just called")
	}
	callInfo := struct {
		S1 string
		S2 string
	}{
		S1: s1,
		S2: s2,
	}
	mock.lockUpsert.Lock()
	mock.calls.Upsert = append(mock.calls.Upsert, callInfo)
	mock.lockUpsert.Unlock()
	return mock.UpsertFunc(s1, s2)
}

// UpsertCalls gets all the calls that were made to Upsert.
// Check the length with:
//     len(mockedKV.UpsertCalls())
func (mock *KVMock) UpsertCalls() []struct {
	S1 string
	S2 string
} {
	var calls []struct {
		S1 string
		S2 string
	}
	mock.lockUpsert.RLock()
	calls = mock.calls.Upsert
	mock.lockUpsert.RUnlock()
	return calls
}
