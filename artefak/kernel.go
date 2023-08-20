package artefak

import (
	"fmt"
	"net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

type Artefak struct {
	router map[string]HandlerFunc
}

func New() *Artefak {
	return &Artefak{router: make(map[string]HandlerFunc)}
}

func (artefak *Artefak) Route(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	artefak.router[key] = handler
}

func (artefak *Artefak) GET(pattern string, handler HandlerFunc) {
	artefak.Route("GET", pattern, handler)
}

func (artefak *Artefak) POST(pattern string, handler HandlerFunc) {
	artefak.Route("POST", pattern, handler)
}

func (artefak *Artefak) Run(addr string) (err error) {
	return http.ListenAndServe(addr, artefak)
}

func (artefak *Artefak) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := artefak.router[key]; ok {
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 Not Found: %s\n", req.URL)
	}
}
