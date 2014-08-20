package rocket

type Dispatcher interface {
	AddRoute(path string, bind Handler, view Renderer)
	FetchRoutes() map[string]*bindObject
}

type Controller struct {
	routes map[string]*bindObject
}

func NewController() *Controller {
	return &Controller{routes: make(map[string]*bindObject)}
}

func (c *Controller) AddRoute(path string, handler Handler, view Renderer) {
	c.routes[path] = &bindObject{handler, view}
}

func (c *Controller) FetchRoutes() map[string]*bindObject {
	return c.routes
}


