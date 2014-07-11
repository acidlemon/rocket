package rocket

import (
	"testing"
	"github.com/acidlemon/rocket"
)

type MockView struct {
}

func (m *MockView) RenderText(text string) string {
	return text + text
}

func (m *MockView) RenderTexts(texts []string) []string {
	result := []string{}

	for _, v := range texts {
		result = append(result, v + v)
	}

	return result
}

func (m *MockView) RenderJSON(data rocket.RenderVars) string {
	return "MockView.RenderJSON"
}

func (m *MockView) Render(tmpl string, data rocket.RenderVars) string {
	return "MockView.Render(" + tmpl + ")"
}


func DummyContext() *rocket.Context {
	return &rocket.Context{
		View: &MockView{},
		Res: &rocket.Response{},
	}
}

func TestRenderer(t *testing.T) {
	c := DummyContext()

	c.RenderText("powawa")

	if c.Res.Body[0] != "powawapowawa" {
		t.Fatal("RenderText failed")
	}

	c.RenderTexts([]string{"hoge", "powawa"})

	if c.Res.Body[0] != "hogehoge" {
		t.Fatal("RenderTexts[0] is not hogehoge")
	}
	if c.Res.Body[1] != "powawapowawa" {
		t.Fatal("RenderTexts[1] is not powawapowawa")
	}

	c.RenderJSON(rocket.RenderVars{ "Cat": "nya" })

	if c.Res.Body[0] != "MockView.RenderJSON" {
		t.Fatal("RenderJSON does not work properly")
	}

	c.Render("powawa", rocket.RenderVars{"Cat":"mya-"})

	if c.Res.Body[0] != "MockView.Render(powawa)" {
		t.Fatal("Render does not work properly")
	}

}


