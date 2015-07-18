package acceptance_stubs

import (
	sync "sync"

	alias1 "github.com/momchil-atanasov/gostub/acceptance/external/external_dup"
)

type ArraySupportStub struct {
	MethodStub        func(arg1 [3]alias1.Address) (result1 [3]alias1.Address)
	methodMutex       sync.RWMutex
	methodArgsForCall []struct {
		arg1 [3]alias1.Address
	}
	methodReturns struct {
		result1 [3]alias1.Address
	}
}

func (stub *ArraySupportStub) Method(arg1 [3]alias1.Address) [3]alias1.Address {
	stub.methodMutex.Lock()
	defer stub.methodMutex.Unlock()
	stub.methodArgsForCall = append(stub.methodArgsForCall, struct {
		arg1 [3]alias1.Address
	}{arg1})
	if stub.MethodStub != nil {
		return stub.MethodStub(arg1)
	} else {
		return stub.methodReturns.result1
	}
}
func (stub *ArraySupportStub) MethodCallCount() int {
	stub.methodMutex.RLock()
	defer stub.methodMutex.RUnlock()
	return len(stub.methodArgsForCall)
}
func (stub *ArraySupportStub) MethodArgsForCall(index int) [3]alias1.Address {
	stub.methodMutex.RLock()
	defer stub.methodMutex.RUnlock()
	return stub.methodArgsForCall[index].arg1
}
func (stub *ArraySupportStub) MethodReturns(result1 [3]alias1.Address) {
	stub.methodMutex.Lock()
	defer stub.methodMutex.Unlock()
	stub.methodReturns = struct {
		result1 [3]alias1.Address
	}{result1}
}
