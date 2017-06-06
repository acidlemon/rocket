package rocket

import (
	"fmt"
	"net/http"

	"github.com/naoina/denco"
)

const MethodAny string = "any"

type Dispatcher interface {
	AddRoute(path string, bind Handler, m ...Middleware)
	AddMethodRoute(method, path string, bind Handler, m ...Middleware)
	Lookup(method, path string) (*bindObject, Args, bool)
}

type dispatcher struct {
	routes  map[string]map[string]interface{} // map[httpMethod]map[route]
	routers map[string]*denco.Router
	view    Renderer
}

func (d *dispatcher) init() {
	d.routes = map[string]map[string]interface{}{
		MethodAny: make(map[string]interface{}),
		// CONNECT, TRACE is not supported
		http.MethodGet:     make(map[string]interface{}),
		http.MethodPost:    make(map[string]interface{}),
		http.MethodHead:    make(map[string]interface{}),
		http.MethodPut:     make(map[string]interface{}),
		http.MethodPatch:   make(map[string]interface{}),
		http.MethodDelete:  make(map[string]interface{}),
		http.MethodOptions: make(map[string]interface{}),
	}
}

func (d *dispatcher) AddRoute(path string, bind Handler, m ...Middleware) {
	if d.routes == nil {
		d.init()
	}

	d.routes[MethodAny][path] = &bindObject{bind, d.view}
}

func (d *dispatcher) AddMethodRoute(method, path string, bind Handler, m ...Middleware) {
	if d.routes == nil {
		d.init()
	}

	switch method {
	case http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions:
		d.routes[method][path] = &bindObject{bind, d.view}
	default:
		// not supported method
		panic(fmt.Sprintf(`HTTP method %s is not supported`, method))
	}
}

func (d *dispatcher) Lookup(method, path string) (*bindObject, Args, bool) {
	// build it first
	if d.routers == nil {
		d.buildRouter()
	}

	bind, pathParams, found := d.routers[method].Lookup(path)
	if !found {
		// fallback
		bind, pathParams, found = d.routers[MethodAny].Lookup(path)
	}

	if !found {
		return nil, Args{}, false
	}

	var args = Args{}
	for _, v := range pathParams {
		args[v.Name] = v.Value
	}

	return bind.(*bindObject), args, true
}

func (d *dispatcher) buildRouter() {
	if d.routes == nil {
		d.init()
	}
	d.routers = make(map[string]*denco.Router, 8)

	for method, r := range d.routes {
		records := []denco.Record{}

		for k, v := range r {
			records = append(records, denco.NewRecord(k, v))
		}

		d.routers[method] = denco.New()
		err := d.routers[method].Build(records)
		if err != nil {
			panic(err)
		}
	}
}

func (d *dispatcher) mount(mountOn string, target map[string]map[string]interface{}) {
	if d.routes == nil {
		d.init()
	}

	for method, route := range target {
		for path, value := range route {
			d.routes[method][mountOn+path] = value
		}
	}
}

type controller struct {
	dispatcher
	mount string
}

func NewController(view Renderer) Controller {
	return &controller{
		dispatcher: dispatcher{view: view},
	}
}

type Controller interface {
	Dispatcher
	SetMount(string)
	GetMount() string
	GetRoutes() map[string]map[string]interface{}
}

func (c *controller) SetMount(m string) {
	c.mount = m
}

func (c *controller) GetMount() string {
	return c.mount
}

func (c *controller) GetRoutes() map[string]map[string]interface{} {
	return c.dispatcher.routes
}
