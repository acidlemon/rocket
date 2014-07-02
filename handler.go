package sixfold

import (
	"fmt"
	"net"
	"net/http"
	"github.com/naoina/kocha-urlrouter"
	_ "github.com/naoina/kocha-urlrouter/doublearray"
)

var theApp *WebApp

type Handler func(w http.ResponseWriter, r *http.Request)

type WebApp struct {
	router urlrouter.URLRouter
	routes map[string]route
}

type route struct {
	c RequestHandler
}

type BindObject struct {
	Method Handler
//	w http.ResponseWriter
//	r *http.Request
}

func (b *BindObject) HandleRequest(w http.ResponseWriter, r *http.Request) {
	b.Method(w, r)
}


func NewWebApp() *WebApp {
	app := new(WebApp)
	return app.Init()
}

func (app *WebApp) Init() *WebApp{
	router := urlrouter.NewURLRouter("doublearray")

	app.router = router
	app.routes = make(map[string]route)

	theApp = app

	return app
}

//func (app *WebApp) AddRoute(path string, c RequestHandler) {
func (app *WebApp) AddRoute(path string, bind func(w http.ResponseWriter, r *http.Request)) {
	app.routes[path] = route{ &BindObject{bind} }
}

func (app *WebApp) BuildRouter() {
	records := []urlrouter.Record{}

	for k, v := range app.routes {
		records = append(records, urlrouter.NewRecord(k, &v))
	}

	app.router.Build(records)
}


func (app* WebApp) Start(listener net.Listener) {
	http.HandleFunc("/", handler)
	http.Serve(listener, nil)
}

func handler(w http.ResponseWriter, req *http.Request) {
	theApp.handler(w, req)
}

func (app *WebApp) handler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("req.URL.Path: ", req.URL.Path)
	r, _ := app.router.Lookup(req.URL.Path)
	fmt.Printf("%v\n", r.(*route).c.(*BindObject).Method)

	r.(*route).c.HandleRequest(w, req)
}


