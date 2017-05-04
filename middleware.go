package rocket

import "net/http"

type Middleware interface {
	Handler(w http.ResponseWriter, r *http.Request)
}
