// Code generated by counterfeiter. DO NOT EDIT.
package filefakes

import (
	"context"
	"sync"

	v1 "github.com/nginx/agent/v3/api/grpc/mpi/v1"
)

type FakeFileManagerServiceInterface struct {
	ClearCacheStub        func()
	clearCacheMutex       sync.RWMutex
	clearCacheArgsForCall []struct {
	}
	ConfigApplyStub        func(context.Context, *v1.ConfigApplyRequest) (bool, error)
	configApplyMutex       sync.RWMutex
	configApplyArgsForCall []struct {
		arg1 context.Context
		arg2 *v1.ConfigApplyRequest
	}
	configApplyReturns struct {
		result1 bool
		result2 error
	}
	configApplyReturnsOnCall map[int]struct {
		result1 bool
		result2 error
	}
	RollbackStub        func(context.Context, string) error
	rollbackMutex       sync.RWMutex
	rollbackArgsForCall []struct {
		arg1 context.Context
		arg2 string
	}
	rollbackReturns struct {
		result1 error
	}
	rollbackReturnsOnCall map[int]struct {
		result1 error
	}
	SetIsConnectedStub        func(bool)
	setIsConnectedMutex       sync.RWMutex
	setIsConnectedArgsForCall []struct {
		arg1 bool
	}
	UpdateFileStub        func(context.Context, string, *v1.File) error
	updateFileMutex       sync.RWMutex
	updateFileArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 *v1.File
	}
	updateFileReturns struct {
		result1 error
	}
	updateFileReturnsOnCall map[int]struct {
		result1 error
	}
	UpdateOverviewStub        func(context.Context, string, []*v1.File) error
	updateOverviewMutex       sync.RWMutex
	updateOverviewArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 []*v1.File
	}
	updateOverviewReturns struct {
		result1 error
	}
	updateOverviewReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeFileManagerServiceInterface) ClearCache() {
	fake.clearCacheMutex.Lock()
	fake.clearCacheArgsForCall = append(fake.clearCacheArgsForCall, struct {
	}{})
	stub := fake.ClearCacheStub
	fake.recordInvocation("ClearCache", []interface{}{})
	fake.clearCacheMutex.Unlock()
	if stub != nil {
		fake.ClearCacheStub()
	}
}

func (fake *FakeFileManagerServiceInterface) ClearCacheCallCount() int {
	fake.clearCacheMutex.RLock()
	defer fake.clearCacheMutex.RUnlock()
	return len(fake.clearCacheArgsForCall)
}

func (fake *FakeFileManagerServiceInterface) ClearCacheCalls(stub func()) {
	fake.clearCacheMutex.Lock()
	defer fake.clearCacheMutex.Unlock()
	fake.ClearCacheStub = stub
}

func (fake *FakeFileManagerServiceInterface) ConfigApply(arg1 context.Context, arg2 *v1.ConfigApplyRequest) (bool, error) {
	fake.configApplyMutex.Lock()
	ret, specificReturn := fake.configApplyReturnsOnCall[len(fake.configApplyArgsForCall)]
	fake.configApplyArgsForCall = append(fake.configApplyArgsForCall, struct {
		arg1 context.Context
		arg2 *v1.ConfigApplyRequest
	}{arg1, arg2})
	stub := fake.ConfigApplyStub
	fakeReturns := fake.configApplyReturns
	fake.recordInvocation("ConfigApply", []interface{}{arg1, arg2})
	fake.configApplyMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeFileManagerServiceInterface) ConfigApplyCallCount() int {
	fake.configApplyMutex.RLock()
	defer fake.configApplyMutex.RUnlock()
	return len(fake.configApplyArgsForCall)
}

func (fake *FakeFileManagerServiceInterface) ConfigApplyCalls(stub func(context.Context, *v1.ConfigApplyRequest) (bool, error)) {
	fake.configApplyMutex.Lock()
	defer fake.configApplyMutex.Unlock()
	fake.ConfigApplyStub = stub
}

func (fake *FakeFileManagerServiceInterface) ConfigApplyArgsForCall(i int) (context.Context, *v1.ConfigApplyRequest) {
	fake.configApplyMutex.RLock()
	defer fake.configApplyMutex.RUnlock()
	argsForCall := fake.configApplyArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeFileManagerServiceInterface) ConfigApplyReturns(result1 bool, result2 error) {
	fake.configApplyMutex.Lock()
	defer fake.configApplyMutex.Unlock()
	fake.ConfigApplyStub = nil
	fake.configApplyReturns = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *FakeFileManagerServiceInterface) ConfigApplyReturnsOnCall(i int, result1 bool, result2 error) {
	fake.configApplyMutex.Lock()
	defer fake.configApplyMutex.Unlock()
	fake.ConfigApplyStub = nil
	if fake.configApplyReturnsOnCall == nil {
		fake.configApplyReturnsOnCall = make(map[int]struct {
			result1 bool
			result2 error
		})
	}
	fake.configApplyReturnsOnCall[i] = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *FakeFileManagerServiceInterface) Rollback(arg1 context.Context, arg2 string) error {
	fake.rollbackMutex.Lock()
	ret, specificReturn := fake.rollbackReturnsOnCall[len(fake.rollbackArgsForCall)]
	fake.rollbackArgsForCall = append(fake.rollbackArgsForCall, struct {
		arg1 context.Context
		arg2 string
	}{arg1, arg2})
	stub := fake.RollbackStub
	fakeReturns := fake.rollbackReturns
	fake.recordInvocation("Rollback", []interface{}{arg1, arg2})
	fake.rollbackMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeFileManagerServiceInterface) RollbackCallCount() int {
	fake.rollbackMutex.RLock()
	defer fake.rollbackMutex.RUnlock()
	return len(fake.rollbackArgsForCall)
}

func (fake *FakeFileManagerServiceInterface) RollbackCalls(stub func(context.Context, string) error) {
	fake.rollbackMutex.Lock()
	defer fake.rollbackMutex.Unlock()
	fake.RollbackStub = stub
}

func (fake *FakeFileManagerServiceInterface) RollbackArgsForCall(i int) (context.Context, string) {
	fake.rollbackMutex.RLock()
	defer fake.rollbackMutex.RUnlock()
	argsForCall := fake.rollbackArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeFileManagerServiceInterface) RollbackReturns(result1 error) {
	fake.rollbackMutex.Lock()
	defer fake.rollbackMutex.Unlock()
	fake.RollbackStub = nil
	fake.rollbackReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeFileManagerServiceInterface) RollbackReturnsOnCall(i int, result1 error) {
	fake.rollbackMutex.Lock()
	defer fake.rollbackMutex.Unlock()
	fake.RollbackStub = nil
	if fake.rollbackReturnsOnCall == nil {
		fake.rollbackReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.rollbackReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeFileManagerServiceInterface) SetIsConnected(arg1 bool) {
	fake.setIsConnectedMutex.Lock()
	fake.setIsConnectedArgsForCall = append(fake.setIsConnectedArgsForCall, struct {
		arg1 bool
	}{arg1})
	stub := fake.SetIsConnectedStub
	fake.recordInvocation("SetIsConnected", []interface{}{arg1})
	fake.setIsConnectedMutex.Unlock()
	if stub != nil {
		fake.SetIsConnectedStub(arg1)
	}
}

func (fake *FakeFileManagerServiceInterface) SetIsConnectedCallCount() int {
	fake.setIsConnectedMutex.RLock()
	defer fake.setIsConnectedMutex.RUnlock()
	return len(fake.setIsConnectedArgsForCall)
}

func (fake *FakeFileManagerServiceInterface) SetIsConnectedCalls(stub func(bool)) {
	fake.setIsConnectedMutex.Lock()
	defer fake.setIsConnectedMutex.Unlock()
	fake.SetIsConnectedStub = stub
}

func (fake *FakeFileManagerServiceInterface) SetIsConnectedArgsForCall(i int) bool {
	fake.setIsConnectedMutex.RLock()
	defer fake.setIsConnectedMutex.RUnlock()
	argsForCall := fake.setIsConnectedArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeFileManagerServiceInterface) UpdateFile(arg1 context.Context, arg2 string, arg3 *v1.File) error {
	fake.updateFileMutex.Lock()
	ret, specificReturn := fake.updateFileReturnsOnCall[len(fake.updateFileArgsForCall)]
	fake.updateFileArgsForCall = append(fake.updateFileArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 *v1.File
	}{arg1, arg2, arg3})
	stub := fake.UpdateFileStub
	fakeReturns := fake.updateFileReturns
	fake.recordInvocation("UpdateFile", []interface{}{arg1, arg2, arg3})
	fake.updateFileMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeFileManagerServiceInterface) UpdateFileCallCount() int {
	fake.updateFileMutex.RLock()
	defer fake.updateFileMutex.RUnlock()
	return len(fake.updateFileArgsForCall)
}

func (fake *FakeFileManagerServiceInterface) UpdateFileCalls(stub func(context.Context, string, *v1.File) error) {
	fake.updateFileMutex.Lock()
	defer fake.updateFileMutex.Unlock()
	fake.UpdateFileStub = stub
}

func (fake *FakeFileManagerServiceInterface) UpdateFileArgsForCall(i int) (context.Context, string, *v1.File) {
	fake.updateFileMutex.RLock()
	defer fake.updateFileMutex.RUnlock()
	argsForCall := fake.updateFileArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeFileManagerServiceInterface) UpdateFileReturns(result1 error) {
	fake.updateFileMutex.Lock()
	defer fake.updateFileMutex.Unlock()
	fake.UpdateFileStub = nil
	fake.updateFileReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeFileManagerServiceInterface) UpdateFileReturnsOnCall(i int, result1 error) {
	fake.updateFileMutex.Lock()
	defer fake.updateFileMutex.Unlock()
	fake.UpdateFileStub = nil
	if fake.updateFileReturnsOnCall == nil {
		fake.updateFileReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.updateFileReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeFileManagerServiceInterface) UpdateOverview(arg1 context.Context, arg2 string, arg3 []*v1.File) error {
	var arg3Copy []*v1.File
	if arg3 != nil {
		arg3Copy = make([]*v1.File, len(arg3))
		copy(arg3Copy, arg3)
	}
	fake.updateOverviewMutex.Lock()
	ret, specificReturn := fake.updateOverviewReturnsOnCall[len(fake.updateOverviewArgsForCall)]
	fake.updateOverviewArgsForCall = append(fake.updateOverviewArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 []*v1.File
	}{arg1, arg2, arg3Copy})
	stub := fake.UpdateOverviewStub
	fakeReturns := fake.updateOverviewReturns
	fake.recordInvocation("UpdateOverview", []interface{}{arg1, arg2, arg3Copy})
	fake.updateOverviewMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeFileManagerServiceInterface) UpdateOverviewCallCount() int {
	fake.updateOverviewMutex.RLock()
	defer fake.updateOverviewMutex.RUnlock()
	return len(fake.updateOverviewArgsForCall)
}

func (fake *FakeFileManagerServiceInterface) UpdateOverviewCalls(stub func(context.Context, string, []*v1.File) error) {
	fake.updateOverviewMutex.Lock()
	defer fake.updateOverviewMutex.Unlock()
	fake.UpdateOverviewStub = stub
}

func (fake *FakeFileManagerServiceInterface) UpdateOverviewArgsForCall(i int) (context.Context, string, []*v1.File) {
	fake.updateOverviewMutex.RLock()
	defer fake.updateOverviewMutex.RUnlock()
	argsForCall := fake.updateOverviewArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeFileManagerServiceInterface) UpdateOverviewReturns(result1 error) {
	fake.updateOverviewMutex.Lock()
	defer fake.updateOverviewMutex.Unlock()
	fake.UpdateOverviewStub = nil
	fake.updateOverviewReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeFileManagerServiceInterface) UpdateOverviewReturnsOnCall(i int, result1 error) {
	fake.updateOverviewMutex.Lock()
	defer fake.updateOverviewMutex.Unlock()
	fake.UpdateOverviewStub = nil
	if fake.updateOverviewReturnsOnCall == nil {
		fake.updateOverviewReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.updateOverviewReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeFileManagerServiceInterface) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.clearCacheMutex.RLock()
	defer fake.clearCacheMutex.RUnlock()
	fake.configApplyMutex.RLock()
	defer fake.configApplyMutex.RUnlock()
	fake.rollbackMutex.RLock()
	defer fake.rollbackMutex.RUnlock()
	fake.setIsConnectedMutex.RLock()
	defer fake.setIsConnectedMutex.RUnlock()
	fake.updateFileMutex.RLock()
	defer fake.updateFileMutex.RUnlock()
	fake.updateOverviewMutex.RLock()
	defer fake.updateOverviewMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeFileManagerServiceInterface) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}