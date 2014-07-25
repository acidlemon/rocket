package rocket

import (
	"net"
	"net/http"

	"github.com/naoina/denco"
	//	"github.com/acidlemon/go-dumper"
)

type Handler func(CtxData)

type CtxBuilder func(req *http.Request, params denco.Params, view Renderer) CtxData

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

func (b *bindObject) HandleRequest(c CtxData) {
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
	app.server.Serve(listener)
}

func (app *WebApp) Handler(w http.ResponseWriter, req *http.Request) {
	bind, params, _ := app.router.Lookup(req.URL.Path)

	var c CtxData
	c = app.ctxBuilder(req, params, bind.(*bindObject).View)

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
