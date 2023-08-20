package artefak

import (
	"net/http"
	"strings"
)

type HandlerFunc func(*Ctx)

type (
	RouterGroup struct {
		prefix      string
		middlewares []HandlerFunc
		parent      *RouterGroup
		artefak     *Artefak
	}

	Artefak struct {
		*RouterGroup
		router *router
		groups []*RouterGroup
	}
)

func New() *Artefak {
	artefak := &Artefak{router: NewRouter()}
	artefak.RouterGroup = &RouterGroup{artefak: artefak}
	artefak.groups = []*RouterGroup{artefak.RouterGroup}

	return artefak
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	artefak  := group.artefak
	newGroup := &RouterGroup{
		prefix  : group.prefix + prefix,
		parent  : group,
		artefak : artefak,
	}

	artefak.groups = append(artefak.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *RouterGroup) Route(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	group.artefak.router.Route(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.Route("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.Route("POST", pattern, handler)
}

func (artefak *Artefak) Run(addr string) (err error) {
	return http.ListenAndServe(addr, artefak)
}

func (artefak *Artefak) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range artefak.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	c := NewCtx(w, req)
	c.handlers = middlewares
	artefak.router.handle(c)
}
