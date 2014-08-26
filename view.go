package rocket

import (
	"fmt"
	"encoding/json"
	"html/template"
	"bytes"
//	"github.com/acidlemon/go-dumper"
)

type RenderVars map[string]interface{}

type Renderer interface {
	Render(string, RenderVars) string
	RenderText(string) string
	RenderTexts([]string) []string
	RenderJSON(RenderVars) string
}

type View struct {
	BasicTemplates []string
}

func (v *View) RenderText(text string) string {
	return text
}

func (v *View) RenderTexts(texts []string) []string {
	return texts
}

func (v *View) Render(tmplFile string, bind RenderVars) string {
	buf := new(bytes.Buffer)
	var err error
	tmpl := template.Must(template.ParseFiles(tmplFile)) // TODO Cache?
	for _, v := range v.BasicTemplates {
		tmpl, err = tmpl.ParseFiles(v)
		if err != nil {
			panic(err)
		}
	}
	err = tmpl.Execute(buf, bind)
	if err != nil {
		panic(fmt.Sprintf("render error: err=%v", err))
	}

	return buf.String()
}

func (v *View) RenderJSON(data RenderVars) string {
	text, err := json.Marshal(data)

	if err != nil {
		panic(fmt.Sprintf("cannot render json: err=%v", err))
	}

	return string(text)
}


