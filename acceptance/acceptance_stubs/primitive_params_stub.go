package acceptance_stubs

import (
	sync "sync"
)

type PrimitiveParamsStub struct {
	SaveStub        func(arg1 int, arg2 string, arg3 float32)
	saveMutex       sync.RWMutex
	saveArgsForCall []struct {
		arg1 int
		arg2 string
		arg3 float32
	}
}

func (stub *PrimitiveParamsStub) Save(arg1 int, arg2 string, arg3 float32) {
	stub.saveMutex.Lock()
	defer stub.saveMutex.Unlock()
	stub.saveArgsForCall = append(stub.saveArgsForCall, struct {
		arg1 int
		arg2 string
		arg3 float32
	}{arg1, arg2, arg3})
	if stub.SaveStub != nil {
		stub.SaveStub(arg1, arg2, arg3)
	}
}
func (stub *PrimitiveParamsStub) SaveCallCount() int {
	stub.saveMutex.RLock()
	defer stub.saveMutex.RUnlock()
	return len(stub.saveArgsForCall)
}
func (stub *PrimitiveParamsStub) SaveArgsForCall(index int) (int, string, float32) {
	stub.saveMutex.RLock()
	defer stub.saveMutex.RUnlock()
	return stub.saveArgsForCall[index].arg1, stub.saveArgsForCall[index].arg2, stub.saveArgsForCall[index].arg3
}
