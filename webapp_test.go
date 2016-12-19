package rocket

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"context"
	"testing"
)

var view *View = &View{}

func prepareWebApp() *WebApp {
	app := &WebApp{}
	app.SetContextBuilder(NewContext)
	return app
}

func TestBasic(t *testing.T) {
	app := prepareWebApp()
	app.AddRoute("/", func(ctx context.Context, c Context) {
		c.Res().StatusCode = http.StatusOK
		c.RenderText("Hello World!!")
	}, view)
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
	app := prepareWebApp()

	// default handler
	app.AddRoute("/:name", func(ctx context.Context, c Context) {
		c.Res().StatusCode = http.StatusOK
		c.RenderText(fmt.Sprintf("Hello %s!!", c.Args()["name"]))
	}, view)
	// POST handler
	app.AddRouteMethod(http.MethodPost, "/:name", func(ctx context.Context, c Context) {
		c.Res().StatusCode = http.StatusOK
		c.RenderText(fmt.Sprintf("Hello POST %s!!", c.Args()["name"]))
	}, view)

	{
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

	{
		req, err := http.NewRequest("POST", "/foobar", nil)
		if err != nil {
			t.Fatal(err)
		}
		rec := httptest.NewRecorder()
		app.Handler(rec, req)
		if rec.Code != 200 {
			t.Errorf("expected %v, but got %v", 200, rec.Code)
		}
		if rec.Body.String() != "Hello POST foobar!!" {
			t.Errorf("expected %v, but got %v", "Hello POST foobar!!", rec.Body.String())
		}
	}

}
