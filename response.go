package rocket // import "gopkg.in/acidlemon/rocket.v2"

import (
	"net/http"
)

type Response struct {
	StatusCode int
	Body       []string
	Header     http.Header
}

func (res *Response) Write(w http.ResponseWriter) {
	for k, v := range res.Header {
		for _, value := range v {
			w.Header().Add(k, value)
		}
	}

	w.WriteHeader(res.StatusCode)

	for _, str := range res.Body {
		w.Write([]byte(str))
	}
}
