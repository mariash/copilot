package handlers

import (
	"fmt"
	"sync"

	bbsmodels "code.cloudfoundry.org/bbs/models"

	"code.cloudfoundry.org/lager"
)

const CF_APP_PORT = 8080

type CAPIProcess struct {
	GUID             CAPIProcessGUID
	DiegoProcessGUID DiegoProcessGUID
}

type CAPIProcessGUID string
type DiegoProcessGUID string
type Hostname string
type RouteGUID string

type RoutesRepo struct {
	Repo map[RouteGUID]*Route
	sync.Mutex
}

func (r *RoutesRepo) Upsert(route *Route) {
	r.Lock()
	r.Repo[route.GUID] = route
	r.Unlock()
}

func (r *RoutesRepo) Delete(guid RouteGUID) {
	r.Lock()
	delete(r.Repo, guid)
	r.Unlock()
}

func (r *RoutesRepo) Get(guid RouteGUID) (*Route, bool) {
	r.Lock()
	route, ok := r.Repo[guid]
	r.Unlock()
	return route, ok
}

// TODO: probably remove or clean this up, currently using for debugging
func (r *RoutesRepo) List() map[string]string {
	list := make(map[string]string)

	r.Lock()
	for k, v := range r.Repo {
		list[string(k)] = v.Host
	}
	r.Unlock()

	return list
}

//go:generate counterfeiter -o fakes/routes_repo.go --fake-name RoutesRepo . routesRepoInterface
type routesRepoInterface interface {
	List() map[string]string
	Upsert(route *Route)
	Delete(guid RouteGUID)
	Get(guid RouteGUID) (*Route, bool)
}

type RouteMappingsRepo struct {
	Repo map[string]RouteMapping
	sync.Mutex
}

func (m *RouteMappingsRepo) Map(routeMapping RouteMapping) {
	m.Lock()
	m.Repo[routeMapping.Key()] = routeMapping
	m.Unlock()
}

func (m *RouteMappingsRepo) Unmap(routeMapping RouteMapping) {
	m.Lock()
	delete(m.Repo, routeMapping.Key())
	m.Unlock()
}

func (m *RouteMappingsRepo) List() map[string]RouteMapping {
	list := make(map[string]RouteMapping)

	m.Lock()
	for k, v := range m.Repo {
		list[k] = v
	}
	m.Unlock()

	return list
}

//go:generate counterfeiter -o fakes/route_mappings_repo.go --fake-name RouteMappingsRepo . routeMappingsRepoInterface
type routeMappingsRepoInterface interface {
	Map(routeMapping RouteMapping)
	Unmap(routeMapping RouteMapping)
	List() map[string]RouteMapping
}

func (p DiegoProcessGUID) Hostname() string {
	label := string(p)

	if len(label) >= 63 {
		label = label[0:62]
	}

	return fmt.Sprintf("%s.cfapps.internal", label)
}

type BBSClient interface {
	ActualLRPGroups(lager.Logger, bbsmodels.ActualLRPFilter) ([]*bbsmodels.ActualLRPGroup, error)
}

type Route struct {
	GUID RouteGUID
	Host string
}

func (r *Route) Hostname() string {
	return r.Host
}

type RouteMapping struct {
	RouteGUID   RouteGUID
	CAPIProcess *CAPIProcess
}

func (r *RouteMapping) Key() string {
	return string(r.RouteGUID) + "-" + string(r.CAPIProcess.GUID)
}
