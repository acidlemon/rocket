package rocket

import (
)

type RequestHandler interface {
	HandleRequest(c *Context)
}


