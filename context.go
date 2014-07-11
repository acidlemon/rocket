package rocket

import (
//	"fmt"
	"net/http"
//	"github.com/acidlemon/go-dumper"

)

type Context struct {
	Req *http.Request
	Res *Response
	Writer *http.ResponseWriter
	View Renderer
	Stash map[string]interface{}
}

func (c *Context) RenderText(text string) {
	renderText := c.View.RenderText(text)
	c.Res.Body = []string{renderText}
}

func (c *Context) RenderTexts(texts []string) {
	renderTexts := c.View.RenderTexts(texts)
	c.Res.Body = renderTexts
}

func (c *Context) RenderJSON(data RenderVars) {
	renderJson := c.View.RenderJSON(data)
	c.Res.Body = []string{renderJson}
}

func (c *Context) Render(tmpl string, data RenderVars) {
	renderText := c.View.Render(tmpl, data)
	c.Res.Body = []string{renderText}
}


