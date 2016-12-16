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
	router     *denco.Router
	routes     map[string]*bindObject
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
	app.router = denco.New()
	app.routes = make(map[string]*bindObject)
	app.ctxBuilder = NewContext

	return app
}

func (app *WebApp) RegisterController(c Dispatcher) {
	r := c.FetchRoutes()

	for k, v := range r {
		app.routes[k] = v
	}
}

func (app *WebApp) AddRoute(path string, bind Handler, view Renderer) {
	app.routes[path] = &bindObject{bind, view}
}

func (app *WebApp) BuildRouter() {
	records := []denco.Record{}

	for k, v := range app.routes {
		records = append(records, denco.NewRecord(k, v))
	}

	app.router.Build(records)
}

func (app *WebApp) Start(listener net.Listener) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.Handler)
	app.server = &http.Server{Handler: mux}

	log.Println("listen start:", listener.Addr().String())
	app.server.Serve(listener)
}

func (app *WebApp) Handler(w http.ResponseWriter, req *http.Request) {
	bind, pathParams, _ := app.router.Lookup(req.URL.Path)

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
