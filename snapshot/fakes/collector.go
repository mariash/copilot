// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	sync "sync"

	models "code.cloudfoundry.org/copilot/models"
)

type Collector struct {
	CollectStub        func() []*models.RouteWithBackends
	collectMutex       sync.RWMutex
	collectArgsForCall []struct {
	}
	collectReturns struct {
		result1 []*models.RouteWithBackends
	}
	collectReturnsOnCall map[int]struct {
		result1 []*models.RouteWithBackends
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Collector) Collect() []*models.RouteWithBackends {
	fake.collectMutex.Lock()
	ret, specificReturn := fake.collectReturnsOnCall[len(fake.collectArgsForCall)]
	fake.collectArgsForCall = append(fake.collectArgsForCall, struct {
	}{})
	fake.recordInvocation("Collect", []interface{}{})
	fake.collectMutex.Unlock()
	if fake.CollectStub != nil {
		return fake.CollectStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.collectReturns
	return fakeReturns.result1
}

func (fake *Collector) CollectCallCount() int {
	fake.collectMutex.RLock()
	defer fake.collectMutex.RUnlock()
	return len(fake.collectArgsForCall)
}

func (fake *Collector) CollectCalls(stub func() []*models.RouteWithBackends) {
	fake.collectMutex.Lock()
	defer fake.collectMutex.Unlock()
	fake.CollectStub = stub
}

func (fake *Collector) CollectReturns(result1 []*models.RouteWithBackends) {
	fake.collectMutex.Lock()
	defer fake.collectMutex.Unlock()
	fake.CollectStub = nil
	fake.collectReturns = struct {
		result1 []*models.RouteWithBackends
	}{result1}
}

func (fake *Collector) CollectReturnsOnCall(i int, result1 []*models.RouteWithBackends) {
	fake.collectMutex.Lock()
	defer fake.collectMutex.Unlock()
	fake.CollectStub = nil
	if fake.collectReturnsOnCall == nil {
		fake.collectReturnsOnCall = make(map[int]struct {
			result1 []*models.RouteWithBackends
		})
	}
	fake.collectReturnsOnCall[i] = struct {
		result1 []*models.RouteWithBackends
	}{result1}
}

func (fake *Collector) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.collectMutex.RLock()
	defer fake.collectMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Collector) recordInvocation(key string, args []interface{}) {
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
