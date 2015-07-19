package acceptance_stubs

import (
	sync "sync"

	alias1 "github.com/momchil-atanasov/gostub/acceptance/external/external_dup"
)

type InterfaceSupportStub struct {
	MethodStub func(arg1 interface {
		alias1.Runner
		ResolveAddress(alias1.Address) alias1.Address
	}) (result1 interface {
		alias1.Runner
		ProcessAddress(alias1.Address) alias1.Address
	})
	methodMutex       sync.RWMutex
	methodArgsForCall []struct {
		arg1 interface {
			alias1.Runner
			ResolveAddress(alias1.Address) alias1.Address
		}
	}
	methodReturns struct {
		result1 interface {
			alias1.Runner
			ProcessAddress(alias1.Address) alias1.Address
		}
	}
}

func (stub *InterfaceSupportStub) Method(arg1 interface {
	alias1.Runner
	ResolveAddress(alias1.Address) alias1.Address
}) interface {
	alias1.Runner
	ProcessAddress(alias1.Address) alias1.Address
} {
	stub.methodMutex.Lock()
	defer stub.methodMutex.Unlock()
	stub.methodArgsForCall = append(stub.methodArgsForCall, struct {
		arg1 interface {
			alias1.Runner
			ResolveAddress(alias1.Address) alias1.Address
		}
	}{arg1})
	if stub.MethodStub != nil {
		return stub.MethodStub(arg1)
	} else {
		return stub.methodReturns.result1
	}
}
func (stub *InterfaceSupportStub) MethodCallCount() int {
	stub.methodMutex.RLock()
	defer stub.methodMutex.RUnlock()
	return len(stub.methodArgsForCall)
}
func (stub *InterfaceSupportStub) MethodArgsForCall(index int) interface {
	alias1.Runner
	ResolveAddress(alias1.Address) alias1.Address
} {
	stub.methodMutex.RLock()
	defer stub.methodMutex.RUnlock()
	return stub.methodArgsForCall[index].arg1
}
func (stub *InterfaceSupportStub) MethodReturns(result1 interface {
	alias1.Runner
	ProcessAddress(alias1.Address) alias1.Address
}) {
	stub.methodMutex.Lock()
	defer stub.methodMutex.Unlock()
	stub.methodReturns = struct {
		result1 interface {
			alias1.Runner
			ProcessAddress(alias1.Address) alias1.Address
		}
	}{result1}
}
