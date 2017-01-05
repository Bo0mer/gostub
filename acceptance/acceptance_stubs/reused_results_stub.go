// Generated by 'github.com/mokiat/gostub'

package acceptance_stubs

import (
	sync "sync"

	alias1 "github.com/mokiat/gostub/acceptance"
)

type ReusedResultsStub struct {
	StubGUID            int
	FullNameStub        func() (result1 string, result2 string)
	fullNameMutex       sync.RWMutex
	fullNameArgsForCall []struct {
	}
	fullNameReturns struct {
		result1 string
		result2 string
	}
}

var _ alias1.ReusedResults = new(ReusedResultsStub)

func (stub *ReusedResultsStub) FullName() (string, string) {
	stub.fullNameMutex.Lock()
	defer stub.fullNameMutex.Unlock()
	stub.fullNameArgsForCall = append(stub.fullNameArgsForCall, struct {
	}{})
	if stub.FullNameStub != nil {
		return stub.FullNameStub()
	} else {
		return stub.fullNameReturns.result1, stub.fullNameReturns.result2
	}
}
func (stub *ReusedResultsStub) FullNameCallCount() int {
	stub.fullNameMutex.RLock()
	defer stub.fullNameMutex.RUnlock()
	return len(stub.fullNameArgsForCall)
}
func (stub *ReusedResultsStub) FullNameReturns(result1 string, result2 string) {
	stub.fullNameMutex.Lock()
	defer stub.fullNameMutex.Unlock()
	stub.fullNameReturns = struct {
		result1 string
		result2 string
	}{result1, result2}
}
