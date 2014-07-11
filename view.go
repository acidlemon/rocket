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
}

func (v *View) RenderText(text string) string {
	return text
}

func (v *View) RenderTexts(texts []string) []string {
	return texts
}

func (v *View) Render(tmplFile string, bind RenderVars) string {
	buf := new(bytes.Buffer)
	tmpl := template.Must(template.ParseFiles(tmplFile)) // TODO Cache?
	err := tmpl.Execute(buf, bind)
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


