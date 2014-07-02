package sixfold

import (
	"net/http"
)

type RequestHandler interface {
	HandleRequest(w http.ResponseWriter, r *http.Request)
}

// Reference Implementation
//type Controller struct {
//	HandleMethod method(func(w http.ResponseWriter, r *http.Request)}
//	Object interface{}
//}

//func (c Controller) HandleRequest(w http.ResponseWriter, r *http.Request) {
//	c.HandleMethod(c.Object, w, r)
//}

