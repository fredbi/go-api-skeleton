// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"github.com/fredbi/go-api-skeleton/api/pkg/repos"
	"github.com/jmoiron/sqlx"
	"sync"
)

// Ensure, that RunnableRepoMock does implement repos.RunnableRepo.
// If this is not the case, regenerate this file with moq.
var _ repos.RunnableRepo = &RunnableRepoMock{}

// RunnableRepoMock is a mock implementation of repos.RunnableRepo.
//
//	func TestSomethingThatUsesRunnableRepo(t *testing.T) {
//
//		// make and configure a mocked repos.RunnableRepo
//		mockedRunnableRepo := &RunnableRepoMock{
//			DBFunc: func() *sqlx.DB {
//				panic("mock out the DB method")
//			},
//			SampleFunc: func() repos.SampleRepo {
//				panic("mock out the Sample method")
//			},
//			StartFunc: func() error {
//				panic("mock out the Start method")
//			},
//			StopFunc: func() error {
//				panic("mock out the Stop method")
//			},
//		}
//
//		// use mockedRunnableRepo in code that requires repos.RunnableRepo
//		// and then make assertions.
//
//	}
type RunnableRepoMock struct {
	// DBFunc mocks the DB method.
	DBFunc func() *sqlx.DB

	// SampleFunc mocks the Sample method.
	SampleFunc func() repos.SampleRepo

	// StartFunc mocks the Start method.
	StartFunc func() error

	// StopFunc mocks the Stop method.
	StopFunc func() error

	// calls tracks calls to the methods.
	calls struct {
		// DB holds details about calls to the DB method.
		DB []struct {
		}
		// Sample holds details about calls to the Sample method.
		Sample []struct {
		}
		// Start holds details about calls to the Start method.
		Start []struct {
		}
		// Stop holds details about calls to the Stop method.
		Stop []struct {
		}
	}
	lockDB     sync.RWMutex
	lockSample sync.RWMutex
	lockStart  sync.RWMutex
	lockStop   sync.RWMutex
}

// DB calls DBFunc.
func (mock *RunnableRepoMock) DB() *sqlx.DB {
	if mock.DBFunc == nil {
		panic("RunnableRepoMock.DBFunc: method is nil but RunnableRepo.DB was just called")
	}
	callInfo := struct {
	}{}
	mock.lockDB.Lock()
	mock.calls.DB = append(mock.calls.DB, callInfo)
	mock.lockDB.Unlock()
	return mock.DBFunc()
}

// DBCalls gets all the calls that were made to DB.
// Check the length with:
//
//	len(mockedRunnableRepo.DBCalls())
func (mock *RunnableRepoMock) DBCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockDB.RLock()
	calls = mock.calls.DB
	mock.lockDB.RUnlock()
	return calls
}

// Sample calls SampleFunc.
func (mock *RunnableRepoMock) Sample() repos.SampleRepo {
	if mock.SampleFunc == nil {
		panic("RunnableRepoMock.SampleFunc: method is nil but RunnableRepo.Sample was just called")
	}
	callInfo := struct {
	}{}
	mock.lockSample.Lock()
	mock.calls.Sample = append(mock.calls.Sample, callInfo)
	mock.lockSample.Unlock()
	return mock.SampleFunc()
}

// SampleCalls gets all the calls that were made to Sample.
// Check the length with:
//
//	len(mockedRunnableRepo.SampleCalls())
func (mock *RunnableRepoMock) SampleCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockSample.RLock()
	calls = mock.calls.Sample
	mock.lockSample.RUnlock()
	return calls
}

// Start calls StartFunc.
func (mock *RunnableRepoMock) Start() error {
	if mock.StartFunc == nil {
		panic("RunnableRepoMock.StartFunc: method is nil but RunnableRepo.Start was just called")
	}
	callInfo := struct {
	}{}
	mock.lockStart.Lock()
	mock.calls.Start = append(mock.calls.Start, callInfo)
	mock.lockStart.Unlock()
	return mock.StartFunc()
}

// StartCalls gets all the calls that were made to Start.
// Check the length with:
//
//	len(mockedRunnableRepo.StartCalls())
func (mock *RunnableRepoMock) StartCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockStart.RLock()
	calls = mock.calls.Start
	mock.lockStart.RUnlock()
	return calls
}

// Stop calls StopFunc.
func (mock *RunnableRepoMock) Stop() error {
	if mock.StopFunc == nil {
		panic("RunnableRepoMock.StopFunc: method is nil but RunnableRepo.Stop was just called")
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
//
//	len(mockedRunnableRepo.StopCalls())
func (mock *RunnableRepoMock) StopCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockStop.RLock()
	calls = mock.calls.Stop
	mock.lockStop.RUnlock()
	return calls
}