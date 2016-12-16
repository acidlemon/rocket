package rocket

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime"

	"context"

	"github.com/naoina/denco"
)

type Handler func(context.Context, Context)

type ContextBuilder func(ctx context.Context, req *http.Request, args Args, view Renderer) context.Context

type WebApp struct {
	routes     RouteMap
	routers    map[string]*denco.Router
	server     *http.Server
	ctxBuilder ContextBuilder
}

type bindObject struct {
	Method Handler
	View   Renderer
}

var (
	errorPage = `<!DOCTYPE html>
<html>
<head>
<title>Internal Server Error</title>
<link rel="stylesheet" href="//maxcdn.bootstrapcdn.com/bootstrap/3.2.0/css/bootstrap.min.css">
<style type="text/css">
body {
    margin:0 20px;
}
</style>
</head>
<body>
<div class="container">
<div class="page-header">
<h1>Internal Server Error</h1>
</div>
<div class="panel panel-danger">
<div class="panel-heading">reason: %v</div>
<div class="panel-body">
<pre>%v</pre>
</div>
</div>
</body>
</html>
`
)

func (b *bindObject) HandleRequest(ctx context.Context) {
	c := ctx.Value(CONTEXT_KEY).(Context)
	defer func() {
		if e := recover(); e != nil {
			buf := make([]byte, 4096)
			runtime.Stack(buf, false)
			stackMsg := string(buf)
			c.Res().StatusCode = http.StatusInternalServerError
			c.Res().Body = []string{fmt.Sprintf(errorPage, e, stackMsg)}
			fmt.Println("Error:", e)
			fmt.Println("Stack:\n", stackMsg)
		}
	}()
	b.Method(ctx, c)
}

func NewWebApp() *WebApp {
	app := new(WebApp)
	return app.Init()
}

func (app *WebApp) SetContextBuilder(f ContextBuilder) {
	app.ctxBuilder = f
}

func (app *WebApp) Init() *WebApp {
	app.routers = make(map[string]*denco.Router, 8)
	app.routes = map[string]map[string]*bindObject{
		MethodAny:          make(map[string]*bindObject),
		http.MethodGet:     make(map[string]*bindObject),
		http.MethodPost:    make(map[string]*bindObject),
		http.MethodHead:    make(map[string]*bindObject),
		http.MethodPut:     make(map[string]*bindObject),
		http.MethodPatch:   make(map[string]*bindObject),
		http.MethodDelete:  make(map[string]*bindObject),
		http.MethodOptions: make(map[string]*bindObject),
	}
	// CONNECT, TRACE is not supported

	app.ctxBuilder = NewContext

	return app
}

func (app *WebApp) RegisterController(c Controller) {
	rm := c.FetchRoutes()

	for method, r := range rm {
		for k, v := range r {
			app.routes[method][k] = v
		}
	}

	app.buildRouter()
}

func (app *WebApp) AddRoute(path string, bind Handler, view Renderer) {
	app.routes.AddRoute(path, bind, view)
	app.buildRouter()
}

func (app *WebApp) AddRouteMethod(method, path string, bind Handler, view Renderer) {
	app.routes.AddRouteMethod(method, path, bind, view)
	app.buildRouter()
}

func (app *WebApp) buildRouter() {
	for method, r := range app.routes {
		records := []denco.Record{}

		for k, v := range r {
			records = append(records, denco.NewRecord(k, v))
		}

		app.routers[method] = denco.New()
		err := app.routers[method].Build(records)
		if err != nil {
			panic(err)
		}
	}
}

func (app *WebApp) Start(listener net.Listener) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.Handler)
	app.server = &http.Server{Handler: mux}

	log.Println("listen start:", listener.Addr().String())
	app.server.Serve(listener)
}

func (app *WebApp) Handler(w http.ResponseWriter, req *http.Request) {
	bind, pathParams, found := app.routers[req.Method].Lookup(req.URL.Path)
	if !found {
		// fallback
		bind, pathParams, _ = app.routers[MethodAny].Lookup(req.URL.Path)
	}

	if bind == nil {
		http.NotFound(w, req)
		return
	}

	var args = Args{}
	for _, v := range pathParams {
		args[v.Name] = v.Value
	}

	ctx := req.Context()
	ctx = app.ctxBuilder(ctx, req, args, bind.(*bindObject).View)

	bind.(*bindObject).HandleRequest(ctx)

	// write response
	c := ctx.Value(CONTEXT_KEY).(Context)
	c.Res().Write(w)

}

func (app *WebApp) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	app.Handler(w, req)
}
