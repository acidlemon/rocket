package rocket

import (
	"net/http"
)

type Response struct {
	StatusCode int
	Body []string
	Header http.Header
}

// もしかして参照渡しじゃなくても大丈夫なのかな引数
func (res *Response) Write(w *http.ResponseWriter) {
	writer := *w
	for k, v := range res.Header {
		for _, value := range v {
			writer.Header().Add(k, value)
		}
	}

	writer.WriteHeader(res.StatusCode)

	for _, str := range res.Body {
		writer.Write([]byte(str))
	}
}




