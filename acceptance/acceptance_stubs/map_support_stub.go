package acceptance_stubs

import (
	sync "sync"

	alias1 "github.com/momchil-atanasov/gostub/acceptance/external/external_dup"
)

type MapSupportStub struct {
	MethodStub        func(arg1 map[alias1.Address]alias1.Address) (result1 map[alias1.Address]alias1.Address)
	methodMutex       sync.RWMutex
	methodArgsForCall []struct {
		arg1 map[alias1.Address]alias1.Address
	}
	methodReturns struct {
		result1 map[alias1.Address]alias1.Address
	}
}

func (stub *MapSupportStub) Method(arg1 map[alias1.Address]alias1.Address) map[alias1.Address]alias1.Address {
	stub.methodMutex.Lock()
	defer stub.methodMutex.Unlock()
	stub.methodArgsForCall = append(stub.methodArgsForCall, struct {
		arg1 map[alias1.Address]alias1.Address
	}{arg1})
	if stub.MethodStub != nil {
		return stub.MethodStub(arg1)
	} else {
		return stub.methodReturns.result1
	}
}
func (stub *MapSupportStub) MethodCallCount() int {
	stub.methodMutex.RLock()
	defer stub.methodMutex.RUnlock()
	return len(stub.methodArgsForCall)
}
func (stub *MapSupportStub) MethodArgsForCall(index int) map[alias1.Address]alias1.Address {
	stub.methodMutex.RLock()
	defer stub.methodMutex.RUnlock()
	return stub.methodArgsForCall[index].arg1
}
func (stub *MapSupportStub) MethodReturns(result1 map[alias1.Address]alias1.Address) {
	stub.methodMutex.Lock()
	defer stub.methodMutex.Unlock()
	stub.methodReturns = struct {
		result1 map[alias1.Address]alias1.Address
	}{result1}
}