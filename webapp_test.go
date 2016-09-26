package rocket

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"context"
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
	app.AddRoute("/", func(ctx context.Context, c Context) {
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

func TestQueryArgs(t *testing.T) {
	app := newTestApp()
	app.AddRoute("/:name", func(ctx context.Context, c Context) {
		c.Res().StatusCode = http.StatusOK
		c.RenderText(fmt.Sprintf("Hello %s!!", c.Args()["name"]))
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
