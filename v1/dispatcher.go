package rocket

import (
	"path"
)

type Dispatcher interface {
	AddRoute(path string, bind Handler, view Renderer)
	FetchRoutes() map[string]*bindObject
}

type Controller struct {
	routes map[string]*bindObject
	Mount string
}

func NewController() *Controller {
	return &Controller{
		routes: make(map[string]*bindObject),
		Mount: "",
	}
}

func (c *Controller) AddRoute(path string, handler Handler, view Renderer) {
	c.routes[path] = &bindObject{handler, view}
}

func (c *Controller) FetchRoutes() map[string]*bindObject {
	if c.Mount != "" {
		routes := make(map[string]*bindObject)
		for k, v := range c.routes {
			routes[path.Join(c.Mount, k)] = v
		}
		return routes
	} else {
		return c.routes
	}
}


