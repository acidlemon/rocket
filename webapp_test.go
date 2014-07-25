package rocket

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"testing"
)

type TestApp struct {
	WebApp
}

var view *View = &View{}

func newTestApp() TestApp {
	app := TestApp{}
	app.Init()
	return app
}

func TestBasic(t *testing.T) {
	app := newTestApp()
	app.AddRoute("/", func(c CtxData) {
		c.Res().StatusCode = http.StatusOK
		c.RenderText("Hello World!!")
	}, view)
	app.BuildRouter()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()
	app.Handler(rec, req)
	if rec.Code != 200 {
		t.Errorf("expected %v, but got %v", 200, rec.Code)
	}
	if rec.Body.String() != "Hello World!!" {
		t.Errorf("expected %v, but got %v", "Hello World!!", rec.Body.String())
	}
}

func TestQueryParams(t *testing.T) {
	app := newTestApp()
	app.AddRoute("/:name", func(c CtxData) {
		c.Res().StatusCode = http.StatusOK
		c.RenderText(fmt.Sprintf("Hello %s!!", c.Params().Get("name")))
	}, view)
	app.BuildRouter()
	req, err := http.NewRequest("GET", "/powawa", nil)
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()
	app.Handler(rec, req)
	if rec.Code != 200 {
		t.Errorf("expected %v, but got %v", 200, rec.Code)
	}
	if rec.Body.String() != "Hello powawa!!" {
		t.Errorf("expected %v, but got %v", "Hello powawa!!", rec.Body.String())
	}
}

func TestMiddleware(t *testing.T) {
	result := ""
	rec := httptest.NewRecorder()
	app := newTestApp()
	app.Use(MiddlewareHandlerFunc(func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		result += "foo"
		next(w, req)
		result += "ban"
	}))
	app.Use(MiddlewareHandlerFunc(func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		result += "bar"
		next(w, req)
		result += "baz"
	}))
	app.Use(MiddlewareHandlerFunc(func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		result += "bat"
		next(w, req)
	}))
	app.ServeHTTP(rec, (*http.Request)(nil))
	if result != "foobarbatbazban" {
		t.Errorf("expected %v, but got %v", "foobarbatbazban", result)
	}
}
