package rocket

import (
	"fmt"
	"net"
	"net/http"
	"github.com/naoina/kocha-urlrouter"
	_ "github.com/naoina/kocha-urlrouter/tst"
//	"github.com/acidlemon/go-dumper"
)

type Handler func(*Context)

type WebApp struct {
	router urlrouter.URLRouter
	routes map[string]*bindObject
}

type bindObject struct {
	Method Handler
	View Renderer
}

func (b *bindObject) HandleRequest(c *Context) {
	fmt.Println("HandleRequest Called")
	b.Method(c)
}


func NewWebApp() *WebApp {
	app := new(WebApp)
	return app.Init()
}

func (app *WebApp) Init() *WebApp {
	router := urlrouter.NewURLRouter("tst")

	app.router = router
	app.routes = make(map[string]*bindObject)

	return app
}

func (app *WebApp) AddRoute(path string, bind func(*Context), view Renderer) {
	app.routes[path] = &bindObject{bind, view}
}

func (app *WebApp) BuildRouter() {
	records := []urlrouter.Record{}

	for k, v := range app.routes {
		fmt.Printf("add %v\n", k)
		records = append(records, urlrouter.NewRecord(k, v))
	}

	app.router.Build(records)
}


func (app* WebApp) Start(listener net.Listener) {
	http.HandleFunc("/", app.handler)
	http.Serve(listener, nil)
}


func (app *WebApp) handler(w http.ResponseWriter, req *http.Request) {
	bind, _ := app.router.Lookup(req.URL.Path)

	// TODO Context Generatorを外から渡せるようにする(デフォルト実装は提供)
	c := &Context{
		Req: req,
		Res: &Response{
			StatusCode: 404,
		},
		View: bind.(*bindObject).View,
		Stash: map[string]interface{}{},
	}

	bind.(*bindObject).HandleRequest(c)

	// write response
	c.Res.Write(&w)
}


