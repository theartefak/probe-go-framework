package artefak

import (
	"net/http"
)

type HandlerFunc func(*Ctx)

type Artefak struct {
	router *router
}

func New() *Artefak {
	return &Artefak{router: NewRouter()}
}

func (artefak *Artefak) Route(method string, pattern string, handler HandlerFunc) {
	artefak.router.Route(method, pattern, handler)
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
	c := NewCtx(w, req)
	artefak.router.handle(c)
}
