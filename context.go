package rocket

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

type contextKey string

const CONTEXT_KEY contextKey = "rocket.Context"

type Context interface {
	Res() *Response
	Req() *http.Request
	View() Renderer
	Args() Args
	Arg(string) string
	ArgInt(string) int64
	Params() Params
	Param(string) ([]string, bool)
	ParamInt(string) ([]int64, bool)
	ParamSingle(string) (string, bool)
	ParamSingleInt(string) (int64, bool)
	SetCookie(*http.Cookie)

	Redirect(string)

	Render(string, RenderVars)
	RenderText(string)
	RenderTexts([]string)
	RenderJSON(interface{})
	Halt(int, string)
}

type c struct {
	req    *http.Request
	res    *Response
	view   Renderer
	args   Args
	params Params
	Stash  map[string]interface{}
}

type Args map[string]string
type Params map[string][]string

func NewContext(request *http.Request, args Args, renderer Renderer) (Context, *http.Request) {
	request.ParseForm()
	params := map[string][]string(request.Form)

	myC := &c{
		req: request,
		res: &Response{
			StatusCode: 200,
			Header:     http.Header{},
		},
		args:   args,
		view:   renderer,
		params: params,
		Stash:  map[string]interface{}{},
	}

	ctx := context.WithValue(request.Context(), CONTEXT_KEY, myC)
	req := request.WithContext(ctx)

	return myC, req
}

func GetContext(ctx context.Context) Context {
	return ctx.Value(CONTEXT_KEY).(Context)
}

func (c *c) Res() *Response {
	return c.res
}

func (c *c) Req() *http.Request {
	return c.req
}

func (c *c) View() Renderer {
	return c.view
}

func (c *c) Args() Args {
	return c.args
}

// obsolete
//func (c *c) Arg(name string) (string, bool) {
//	value, ok := c.args[name]
//
//	return value, ok
//}

func (c *c) Arg(name string) string {
	if value, ok := c.args[name]; !ok {
		panic("Context.Arg could not found key: " + name)
	} else {
		return value
	}
}

func (c *c) ArgInt(name string) int64 {
	if value, ok := c.args[name]; !ok {
		panic("Context.ArgInt could not found key: " + name)
	} else {
		intval, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			panic(fmt.Sprint("Arg %s is expected as integer, actual value is %s", name, value))
		}
		return intval
	}
}

func (c *c) Params() Params {
	return c.params
}

func (c *c) Param(name string) ([]string, bool) {
	slice, ok := c.params[name]
	return slice, ok
}

func (c *c) ParamInt(name string) ([]int64, bool) {
	slice, ok := c.params[name]
	if !ok {
		return nil, false
	}

	result := make([]int64, 0, len(slice))
	for _, str := range slice {
		value, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, false
		}
		result = append(result, value)
	}

	return result, true
}

func (c *c) ParamSingle(name string) (string, bool) {
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

func (c *c) ParamSingleInt(name string) (int64, bool) {
	str, valid := c.ParamSingle(name)
	if !valid {
		return 0, false
	}

	value, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, false
	}

	return value, true
}

func (c *c) SetCookie(cookie *http.Cookie) {
	c.res.Header.Add("Set-Cookie", cookie.String())
}

func (c *c) Redirect(uri string) {
	c.res.StatusCode = http.StatusFound
	c.res.Header.Set("Location", uri)
	c.res.Body = []string{""}
}

func (c *c) RenderText(text string) {
	renderText := c.View().RenderText(text)
	c.res.Body = []string{renderText}
}

func (c *c) RenderTexts(texts []string) {
	renderTexts := c.View().RenderTexts(texts)
	c.res.Body = renderTexts
}

func (c *c) RenderJSON(data interface{}) {
	renderJson := c.View().RenderJSON(data)
	c.res.Body = []string{renderJson}
	c.res.Header.Add("Content-Type", "application/json")
}

func (c *c) Render(tmpl string, data RenderVars) {
	renderText := c.View().Render(tmpl, data)
	c.res.Body = []string{renderText}
}

// Halt does not render using view
func (c *c) Halt(code int, text string) {
	c.res.StatusCode = code
	c.res.Body = []string{text}
	panic(c)
}
