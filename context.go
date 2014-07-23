package rocket

import (
	"net/http"

	"github.com/naoina/denco"
)

type CtxData interface {
	Res() *Response
	Req() *http.Request
	View() Renderer
	Params() denco.Params

	Render(string, RenderVars)
	RenderText(string)
	RenderTexts([]string)
	RenderJSON(RenderVars)
}

type Context struct {
	req    *http.Request
	res    *Response
	view   Renderer
	params denco.Params
	Stash  map[string]interface{}
}

func NewContext(request *http.Request, params denco.Params, renderer Renderer) CtxData {
	c := &Context{
		req: request,
		res: &Response{
			StatusCode: 404,
		},
		params: params,
		view:   renderer,
		Stash:  map[string]interface{}{},
	}

	return c
}

func (c *Context) Res() *Response {
	return c.res
}

func (c *Context) Req() *http.Request {
	return c.req
}

func (c *Context) View() Renderer {
	return c.view
}

func (c *Context) Params() denco.Params {
	return c.params
}

func (c *Context) RenderText(text string) {
	renderText := c.View().RenderText(text)
	c.Res().Body = []string{renderText}
}

func (c *Context) RenderTexts(texts []string) {
	renderTexts := c.View().RenderTexts(texts)
	c.Res().Body = renderTexts
}

func (c *Context) RenderJSON(data RenderVars) {
	renderJson := c.View().RenderJSON(data)
	c.Res().Body = []string{renderJson}
}

func (c *Context) Render(tmpl string, data RenderVars) {
	renderText := c.View().Render(tmpl, data)
	c.Res().Body = []string{renderText}
}
