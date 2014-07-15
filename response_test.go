package rocket

import (
	"fmt"
	"sort"
	"net/http"
	"testing"
)

type MockWriter struct {
	header http.Header
	Headers []string
	Body string
}

func (w *MockWriter) Header() http.Header{
	return w.header
}

func (w *MockWriter) WriteHeader(statusCode int) {
	w.Headers = []string{}

	w.Headers = append(w.Headers, fmt.Sprintf("HTTP/1.1 %d", statusCode))

	for k, v := range w.header {
		var line string
		for _, v2 := range v {
			line += v2
		}
		w.Headers = append(w.Headers, fmt.Sprintf("%s: %s", k, line))
	}
}

func (w *MockWriter) Write(b []byte) (int, error){
	w.Body += string(b)
	return len(b), nil
}


func TestResponse(t *testing.T) {
	res := Response{
		StatusCode: 200,
		Body: []string{"test\n", "today is rainy\n"},
		Header: http.Header{
			"Content-Type": {"text/plain", ";charset=utf8"},
			"Server" : {"Powawa/2.2.4"},
			"Connection": {"close"},
			"Date" : {"Mon, 14 Jul 2014 02:27:00 GMT"},
		},
	}

	writer := &MockWriter{header: http.Header{}}

	res.Write(writer)

	if writer.Body != "test\ntoday is rainy\n" {
		fmt.Println("body: ", writer.Body)
		t.Fatalf("Body does not write correctly")
	}

	// 乱暴だけどソートして順番を保証して1個ずつチェックする
	sorted_headers := sort.StringSlice(writer.Headers)
	sort.Sort(sorted_headers)

	if len(sorted_headers) != 5 {
		t.Fatal("Header line does not have 5 lines")
	}

	if sorted_headers[0] != "Connection: close" {
		t.Fatal("invalid header[0]")
	}
	if sorted_headers[1] != "Content-Type: text/plain;charset=utf8" {
		t.Fatal("invalid header[1]")
	}
	if sorted_headers[2] != "Date: Mon, 14 Jul 2014 02:27:00 GMT" {
		t.Fatal("invalid header[2]")
	}
	if sorted_headers[3] != "HTTP/1.1 200" {
		t.Fatal("invalid header[3]")
	}
	if sorted_headers[4] != "Server: Powawa/2.2.4" {
		t.Fatal("invalid header[4]")
	}

}



