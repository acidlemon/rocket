package rocket

import (
	"log"
	"net/http"
	"time"
)

type MiddlewareHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request, http.HandlerFunc)
}

type MiddlewareHandlerFunc func(http.ResponseWriter, *http.Request, http.HandlerFunc)

func (mhf MiddlewareHandlerFunc) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	mhf(w, req, next)
}

type middleware struct {
	handler MiddlewareHandler
	next    *middleware
}

func (m middleware) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	m.handler.ServeHTTP(w, req, m.next.ServeHTTP)
}

func build(handlers []MiddlewareHandler) middleware {
	var next middleware
	if len(handlers) == 0 {
		return voidMiddleware()
	} else if len(handlers) > 1 {
		next = build(handlers[1:])
	} else {
		next = voidMiddleware()
	}
	return middleware{handlers[0], &next}
}

func voidMiddleware() middleware {
	return middleware{
		handler: MiddlewareHandlerFunc(func(http.ResponseWriter, *http.Request, http.HandlerFunc) {}),
		next:    &middleware{},
	}
}

func WrapMiddlewareHandler(handler http.Handler) MiddlewareHandler {
	return MiddlewareHandlerFunc(func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		handler.ServeHTTP(w, req)
		next(w, req)
	})
}

type LoggerMiddleware struct{}

func (lm *LoggerMiddleware) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	start := time.Now()
	log.Printf("Started %s %s", req.Method, req.URL.Path)
	next(w, req)
	log.Printf("Completed in %v", time.Since(start))
}

type RecoverMiddleware struct{}

func (rm *RecoverMiddleware) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	defer func() {
		if err := recover(); err != nil {
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
		}
	}()
	next(w, req)
}
