package rocket // import "gopkg.in/acidlemon/rocket.v1"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"path/filepath"
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
	TemplateDelims []string // Delims[0] = left, Delims[1] = right
}

func (v *View) RenderText(text string) string {
	return text
}

func (v *View) RenderTexts(texts []string) []string {
	return texts
}

func (v *View) delims() (string, string) {
	left := ""
	right := ""

	if v.TemplateDelims != nil {
		if len(v.TemplateDelims) > 0 {
			left = v.TemplateDelims[0]
		}
		if len(v.TemplateDelims) > 1 {
			right = v.TemplateDelims[1]
		}
	}

	return left, right
}

func (v *View) Render(tmplFile string, bind RenderVars) string {
	buf := new(bytes.Buffer)
	var err error

	tmpl := template.Must(
		template.New(filepath.Base(tmplFile)).Delims(v.delims()).ParseFiles(tmplFile))

	if v.BasicTemplates != nil {
		for _, v := range v.BasicTemplates {
			tmpl, err = tmpl.ParseFiles(v)
			if err != nil {
				panic(err)
			}
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
