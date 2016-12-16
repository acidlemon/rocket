package rocket

import (
	"fmt"
	"net/http"
)

const MethodAny string = "any"

type Dispatcher interface {
	AddRoute(path string, bind Handler, view Renderer)
	AddRouteMethod(method, path string, bind Handler, view Renderer)
}

type RouteMap map[string]map[string]*bindObject // map[httpMethod]map[route]

func newDispatcher() RouteMap {
	return RouteMap{
		MethodAny:          make(map[string]*bindObject),
		http.MethodGet:     make(map[string]*bindObject),
		http.MethodPost:    make(map[string]*bindObject),
		http.MethodHead:    make(map[string]*bindObject),
		http.MethodPut:     make(map[string]*bindObject),
		http.MethodPatch:   make(map[string]*bindObject),
		http.MethodDelete:  make(map[string]*bindObject),
		http.MethodOptions: make(map[string]*bindObject),
	}
}

func (d RouteMap) AddRoute(path string, bind Handler, view Renderer) {
	d[MethodAny][path] = &bindObject{bind, view}
}

func (d RouteMap) AddRouteMethod(method, path string, bind Handler, view Renderer) {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions:
		d[method][path] = &bindObject{bind, view}
	default:
		// not supported method
		panic(fmt.Sprintf(`HTTP method %s is not supported`, method))
	}
}

type controller struct {
	RouteMap
	mount string
}

type Controller interface {
	Dispatcher
	SetMount(string)
	GetMount() string
	FetchRoutes() RouteMap
}

func NewController() Controller {
	return &controller{
		RouteMap: newDispatcher(),
	}
}

func (c *controller) FetchRoutes() RouteMap {
	/*
		if c.Mount != "" {
			routes := make(map[string]*bindObject)
			for k, v := range c.routes {
				routes[path.Join(c.Mount, k)] = v
			}
			return routes
		} else {
			return c.routes
		}
	*/
	return c.RouteMap
}

func (c *controller) SetMount(m string) {
	c.mount = m
}

func (c *controller) GetMount() string {
	return c.mount
}
