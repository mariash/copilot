// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"

	"code.cloudfoundry.org/copilot/api"
	"code.cloudfoundry.org/copilot/models"
)

type BackendSet struct {
	GetStub        func(guid models.DiegoProcessGUID) *api.BackendSet
	getMutex       sync.RWMutex
	getArgsForCall []struct {
		guid models.DiegoProcessGUID
	}
	getReturns struct {
		result1 *api.BackendSet
	}
	getReturnsOnCall map[int]struct {
		result1 *api.BackendSet
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *BackendSet) Get(guid models.DiegoProcessGUID) *api.BackendSet {
	fake.getMutex.Lock()
	ret, specificReturn := fake.getReturnsOnCall[len(fake.getArgsForCall)]
	fake.getArgsForCall = append(fake.getArgsForCall, struct {
		guid models.DiegoProcessGUID
	}{guid})
	fake.recordInvocation("Get", []interface{}{guid})
	fake.getMutex.Unlock()
	if fake.GetStub != nil {
		return fake.GetStub(guid)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.getReturns.result1
}

func (fake *BackendSet) GetCallCount() int {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	return len(fake.getArgsForCall)
}

func (fake *BackendSet) GetArgsForCall(i int) models.DiegoProcessGUID {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	return fake.getArgsForCall[i].guid
}

func (fake *BackendSet) GetReturns(result1 *api.BackendSet) {
	fake.GetStub = nil
	fake.getReturns = struct {
		result1 *api.BackendSet
	}{result1}
}

func (fake *BackendSet) GetReturnsOnCall(i int, result1 *api.BackendSet) {
	fake.GetStub = nil
	if fake.getReturnsOnCall == nil {
		fake.getReturnsOnCall = make(map[int]struct {
			result1 *api.BackendSet
		})
	}
	fake.getReturnsOnCall[i] = struct {
		result1 *api.BackendSet
	}{result1}
}

func (fake *BackendSet) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *BackendSet) recordInvocation(key string, args []interface{}) {
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
