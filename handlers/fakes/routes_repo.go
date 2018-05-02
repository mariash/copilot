// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"

	"code.cloudfoundry.org/copilot/models"
)

type RoutesRepo struct {
	UpsertStub        func(route *models.Route)
	upsertMutex       sync.RWMutex
	upsertArgsForCall []struct {
		route *models.Route
	}
	DeleteStub        func(guid models.RouteGUID)
	deleteMutex       sync.RWMutex
	deleteArgsForCall []struct {
		guid models.RouteGUID
	}
	SyncStub        func(routes []*models.Route)
	syncMutex       sync.RWMutex
	syncArgsForCall []struct {
		routes []*models.Route
	}
	GetStub        func(guid models.RouteGUID) (*models.Route, bool)
	getMutex       sync.RWMutex
	getArgsForCall []struct {
		guid models.RouteGUID
	}
	getReturns struct {
		result1 *models.Route
		result2 bool
	}
	getReturnsOnCall map[int]struct {
		result1 *models.Route
		result2 bool
	}
	ListStub        func() map[string]string
	listMutex       sync.RWMutex
	listArgsForCall []struct{}
	listReturns     struct {
		result1 map[string]string
	}
	listReturnsOnCall map[int]struct {
		result1 map[string]string
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *RoutesRepo) Upsert(route *models.Route) {
	fake.upsertMutex.Lock()
	fake.upsertArgsForCall = append(fake.upsertArgsForCall, struct {
		route *models.Route
	}{route})
	fake.recordInvocation("Upsert", []interface{}{route})
	fake.upsertMutex.Unlock()
	if fake.UpsertStub != nil {
		fake.UpsertStub(route)
	}
}

func (fake *RoutesRepo) UpsertCallCount() int {
	fake.upsertMutex.RLock()
	defer fake.upsertMutex.RUnlock()
	return len(fake.upsertArgsForCall)
}

func (fake *RoutesRepo) UpsertArgsForCall(i int) *models.Route {
	fake.upsertMutex.RLock()
	defer fake.upsertMutex.RUnlock()
	return fake.upsertArgsForCall[i].route
}

func (fake *RoutesRepo) Delete(guid models.RouteGUID) {
	fake.deleteMutex.Lock()
	fake.deleteArgsForCall = append(fake.deleteArgsForCall, struct {
		guid models.RouteGUID
	}{guid})
	fake.recordInvocation("Delete", []interface{}{guid})
	fake.deleteMutex.Unlock()
	if fake.DeleteStub != nil {
		fake.DeleteStub(guid)
	}
}

func (fake *RoutesRepo) DeleteCallCount() int {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	return len(fake.deleteArgsForCall)
}

func (fake *RoutesRepo) DeleteArgsForCall(i int) models.RouteGUID {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	return fake.deleteArgsForCall[i].guid
}

func (fake *RoutesRepo) Sync(routes []*models.Route) {
	var routesCopy []*models.Route
	if routes != nil {
		routesCopy = make([]*models.Route, len(routes))
		copy(routesCopy, routes)
	}
	fake.syncMutex.Lock()
	fake.syncArgsForCall = append(fake.syncArgsForCall, struct {
		routes []*models.Route
	}{routesCopy})
	fake.recordInvocation("Sync", []interface{}{routesCopy})
	fake.syncMutex.Unlock()
	if fake.SyncStub != nil {
		fake.SyncStub(routes)
	}
}

func (fake *RoutesRepo) SyncCallCount() int {
	fake.syncMutex.RLock()
	defer fake.syncMutex.RUnlock()
	return len(fake.syncArgsForCall)
}

func (fake *RoutesRepo) SyncArgsForCall(i int) []*models.Route {
	fake.syncMutex.RLock()
	defer fake.syncMutex.RUnlock()
	return fake.syncArgsForCall[i].routes
}

func (fake *RoutesRepo) Get(guid models.RouteGUID) (*models.Route, bool) {
	fake.getMutex.Lock()
	ret, specificReturn := fake.getReturnsOnCall[len(fake.getArgsForCall)]
	fake.getArgsForCall = append(fake.getArgsForCall, struct {
		guid models.RouteGUID
	}{guid})
	fake.recordInvocation("Get", []interface{}{guid})
	fake.getMutex.Unlock()
	if fake.GetStub != nil {
		return fake.GetStub(guid)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getReturns.result1, fake.getReturns.result2
}

func (fake *RoutesRepo) GetCallCount() int {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	return len(fake.getArgsForCall)
}

func (fake *RoutesRepo) GetArgsForCall(i int) models.RouteGUID {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	return fake.getArgsForCall[i].guid
}

func (fake *RoutesRepo) GetReturns(result1 *models.Route, result2 bool) {
	fake.GetStub = nil
	fake.getReturns = struct {
		result1 *models.Route
		result2 bool
	}{result1, result2}
}

func (fake *RoutesRepo) GetReturnsOnCall(i int, result1 *models.Route, result2 bool) {
	fake.GetStub = nil
	if fake.getReturnsOnCall == nil {
		fake.getReturnsOnCall = make(map[int]struct {
			result1 *models.Route
			result2 bool
		})
	}
	fake.getReturnsOnCall[i] = struct {
		result1 *models.Route
		result2 bool
	}{result1, result2}
}

func (fake *RoutesRepo) List() map[string]string {
	fake.listMutex.Lock()
	ret, specificReturn := fake.listReturnsOnCall[len(fake.listArgsForCall)]
	fake.listArgsForCall = append(fake.listArgsForCall, struct{}{})
	fake.recordInvocation("List", []interface{}{})
	fake.listMutex.Unlock()
	if fake.ListStub != nil {
		return fake.ListStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.listReturns.result1
}

func (fake *RoutesRepo) ListCallCount() int {
	fake.listMutex.RLock()
	defer fake.listMutex.RUnlock()
	return len(fake.listArgsForCall)
}

func (fake *RoutesRepo) ListReturns(result1 map[string]string) {
	fake.ListStub = nil
	fake.listReturns = struct {
		result1 map[string]string
	}{result1}
}

func (fake *RoutesRepo) ListReturnsOnCall(i int, result1 map[string]string) {
	fake.ListStub = nil
	if fake.listReturnsOnCall == nil {
		fake.listReturnsOnCall = make(map[int]struct {
			result1 map[string]string
		})
	}
	fake.listReturnsOnCall[i] = struct {
		result1 map[string]string
	}{result1}
}

func (fake *RoutesRepo) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.upsertMutex.RLock()
	defer fake.upsertMutex.RUnlock()
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	fake.syncMutex.RLock()
	defer fake.syncMutex.RUnlock()
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	fake.listMutex.RLock()
	defer fake.listMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *RoutesRepo) recordInvocation(key string, args []interface{}) {
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
