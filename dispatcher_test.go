package rocket

import (
	"net/http"
	"reflect"
	"testing"

	"context"
)

func TestDispatcher(t *testing.T) {
	d := dispatcher{}

	var f, f2 Handler
	f = func(ctx context.Context, c Context) {
		c.RenderText("home called")
	}
	f2 = func(ctx context.Context, c Context) {
		c.RenderText("patch home called")
	}

	d.AddRoute("/home", f, &View{})

	bind, args, found := d.Lookup("GET", "/home")
	if !found {
		t.Fatal("/home is not found")
	}

	if len(args) != 0 {
		t.Fatal("unexpected args")
	}

	if reflect.ValueOf(bind.Method).Pointer() != reflect.ValueOf(f).Pointer() {
		t.Fatal("invalid handler")
	}

	d.AddRouteMethod(http.MethodPatch, "/:type", f2, &View{})
	d.buildRouter() // build again

	bind, args, found = d.Lookup("PATCH", "/home")
	if !found {
		t.Fatal("/home is not found")
	}

	if len(args) != 1 {
		t.Fatal("unexpected args")
	}

	if reflect.ValueOf(bind.Method).Pointer() != reflect.ValueOf(f2).Pointer() {
		t.Fatal("invalid handler")
	}

}
