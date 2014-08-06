package rocket

import (
	"net/http"
)

type CtxData interface {
	Res() *Response
	Req() *http.Request
	View() Renderer
	Args() Args
	Arg(string) (string, bool)
	Params() Params
	Param(string) ([]string, bool)
	ParamSingle(string) (string, bool)

	Render(string, RenderVars)
	RenderText(string)
	RenderTexts([]string)
	RenderJSON(RenderVars)
}

type Context struct {
	req    *http.Request
	res    *Response
	view   Renderer
	args   Args
	params Params
	Stash  map[string]interface{}
}

type Args map[string]string
type Params map[string][]string

func NewContext(request *http.Request, args Args, renderer Renderer) CtxData {
	request.ParseForm()
	params := map[string][]string(request.Form)

	c := &Context{
		req: request,
		res: &Response{
			StatusCode: 200,
		},
		args: args,
		view: renderer,
		params: params,
		Stash: map[string]interface{}{},
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

func (c *Context) Args() Args {
	return c.args
}

func (c *Context) Arg(name string) (string, bool) {
	value, ok := c.args[name]

	return value, ok
}

func (c *Context) Params() Params {
	return c.params
}

func (c *Context) Param(name string) ([]string, bool) {
	slice, ok := c.params[name]
	return slice, ok
}

func (c *Context) ParamSingle(name string) (string, bool) {
	var value string
	valid := false
	if slice, ok := c.params[name]; ok {
		if len(slice) > 0 {
			value = slice[0]
			valid = true
		}
	}

	return value, valid
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
