package rocket

import (
	"net"
	"net/http"
	"github.com/naoina/kocha-urlrouter"
	_ "github.com/naoina/kocha-urlrouter/tst"
//	"github.com/acidlemon/go-dumper"
)

type Handler func(CtxData)

type CtxBuilder func(req *http.Request, view Renderer) CtxData

type WebApp struct {
	router urlrouter.URLRouter
	routes map[string]*bindObject
	server *http.Server
	ctxBuilder CtxBuilder
}

type bindObject struct {
	Method Handler
	View Renderer
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
	router := urlrouter.NewURLRouter("tst")

	app.router = router
	app.routes = make(map[string]*bindObject)
	app.ctxBuilder = NewContext

	return app
}

func (app *WebApp) AddRoute(path string, bind Handler, view Renderer) {
	app.routes[path] = &bindObject{bind, view}
}

func (app *WebApp) BuildRouter() {
	records := []urlrouter.Record{}

	for k, v := range app.routes {
		records = append(records, urlrouter.NewRecord(k, v))
	}

	app.router.Build(records)
}

func (app *WebApp) Start(listener net.Listener) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.Handler)
	app.server = &http.Server{Handler: mux}
	app.server.Serve(listener)
}

func (app *WebApp) Handler(w http.ResponseWriter, req *http.Request) {
	bind, _ := app.router.Lookup(req.URL.Path)

	var c CtxData
	c = app.ctxBuilder(req, bind.(*bindObject).View)

	bind.(*bindObject).HandleRequest(c)

	// write response
	c.Res().Write(w)
}

func (app *WebApp) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	app.Handler(w, req)
}


