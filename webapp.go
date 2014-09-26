package rocket

import (
	"fmt"
	"net"
	"net/http"
	"runtime"

	"github.com/naoina/denco"
	//	"github.com/acidlemon/go-dumper"
)

type Handler func(CtxData)

type CtxBuilder func(req *http.Request, args Args, view Renderer) CtxData

type WebApp struct {
	router      *denco.Router
	routes      map[string]*bindObject
	server      *http.Server
	ctxBuilder  CtxBuilder
	middlewares []MiddlewareHandler
	middleware  middleware
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

func (b *bindObject) HandleRequest(c CtxData) {
	defer func() {
		if e := recover(); e != nil {
			buf := make([]byte, 4096)
			runtime.Stack(buf, false)
			stackMsg := string(buf)
			c.Res().StatusCode = http.StatusInternalServerError
			c.Res().Body = []string{ fmt.Sprintf(errorPage, e, stackMsg)}
			fmt.Println("Error:", e)
			fmt.Println("Stack:\n", stackMsg)
		}
	}()
	b.Method(c)
}

func NewWebApp() *WebApp {
	app := new(WebApp)
	return app.Init()
}

func (app *WebApp) SetContextBuilder(f CtxBuilder) {
	app.ctxBuilder = f
}

func (app *WebApp) Init() *WebApp {
	app.router = denco.New()
	app.routes = make(map[string]*bindObject)
	app.ctxBuilder = NewContext
	app.middlewares = []MiddlewareHandler{}

	return app
}

func (app *WebApp) RegisterController(c Dispatcher){
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
	app.Use(WrapMiddlewareHandler(http.HandlerFunc(app.Handler)))

	mux := http.NewServeMux()
	mux.Handle("/", app)
	app.server = &http.Server{Handler: mux}

	fmt.Println("listen start:", listener.Addr().String())
	app.server.Serve(listener)
}

func (app *WebApp) Handler(w http.ResponseWriter, req *http.Request) {
	bind, pathParams, _ := app.router.Lookup(req.URL.Path)

	if bind == nil {
		http.NotFound(w, req);
		return
	}

	var args = Args{}
	for _, v := range pathParams {
		args[v.Name] = v.Value
	}

	var c CtxData
	c = app.ctxBuilder(req, args, bind.(*bindObject).View)

	bind.(*bindObject).HandleRequest(c)

	// write response
	c.Res().Write(w)
}

func (app *WebApp) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	app.middleware.ServeHTTP(w, req)
}

func (app *WebApp) Use(mh MiddlewareHandler) {
	app.middlewares = append(app.middlewares, mh)
	app.middleware = build(app.middlewares)
}
