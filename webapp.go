package rocket

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime"

	"context"
)

type Handler func(context.Context, Context)

type ContextBuilder func(req *http.Request, args Args, view Renderer) (Context, *http.Request)

type WebApp struct {
	dispatcher
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

func (app *WebApp) SetContextBuilder(f ContextBuilder) {
	app.ctxBuilder = f
}

func (app *WebApp) RegisterController(c Controller) {
	app.mount(c.GetMount(), c.GetRoutes())
}

func (app *WebApp) Start(listener net.Listener) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.Handler)
	app.server = &http.Server{Handler: mux}

	log.Println("listen start:", listener.Addr().String())
	app.server.Serve(listener)
}

func (app *WebApp) Handler(w http.ResponseWriter, req *http.Request) {
	bind, args, found := app.Lookup(req.Method, req.URL.Path)

	if !found {
		http.NotFound(w, req)
		return
	}

	if app.ctxBuilder == nil {
		// set default context builder
		app.ctxBuilder = NewContext
	}
	c, req := app.ctxBuilder(req, args, bind.View)

	bind.HandleRequest(req.Context())

	// write response
	c.Res().Write(w)

}

func (app *WebApp) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	app.Handler(w, req)
}
