package rocket

import (
	"context"
	"net/http"
	"testing"
)

type MockView struct {
}

func (m *MockView) RenderText(text string) string {
	return text + text
}

func (m *MockView) RenderTexts(texts []string) []string {
	result := []string{}

	for _, v := range texts {
		result = append(result, v+v)
	}

	return result
}

func (m *MockView) RenderJSON(data interface{}) string {
	return "MockView.RenderJSON"
}

func (m *MockView) Render(tmpl string, data RenderVars) string {
	return "MockView.Render(" + tmpl + ")"
}

func DummyContext() *c {
	req := &http.Request{}
	view := &MockView{}
	args := Args{}

	ctx := context.Background()
	ctx = NewContext(ctx, req, args, view)
	c := ctx.Value(CONTEXT_KEY).(*c)

	return c
}

func TestRenderer(t *testing.T) {
	c := DummyContext()

	c.RenderText("powawa")

	if c.Res().Body[0] != "powawapowawa" {
		t.Fatal("RenderText failed")
	}

	c.RenderTexts([]string{"hoge", "powawa"})

	if c.Res().Body[0] != "hogehoge" {
		t.Fatal("RenderTexts[0] is not hogehoge")
	}
	if c.Res().Body[1] != "powawapowawa" {
		t.Fatal("RenderTexts[1] is not powawapowawa")
	}

	c.RenderJSON(RenderVars{"Cat": "nya"})

	if c.Res().Body[0] != "MockView.RenderJSON" {
		t.Fatal("RenderJSON does not work properly")
	}

	c.Render("powawa", RenderVars{"Cat": "mya-"})

	if c.Res().Body[0] != "MockView.Render(powawa)" {
		t.Fatal("Render does not work properly")
	}
}

func TestContextAccessors(t *testing.T) {
	c := DummyContext()

	if c.Req() == nil {
		t.Fatal("something wrong on c.Req()")
	}

	if c.Res() == nil {
		t.Fatal("c.Res() returns nil")
	}

	if c.View() == nil {
		t.Fatal("something wrong on c.View()")
	}

	if c.Params() == nil {
		t.Fatal("something wrong on c.Params()")
	}
}
