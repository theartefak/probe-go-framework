package artefak

import (
	"html/template"
	"net/http"
	"path"
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
		router        *router
		groups        []*RouterGroup
		htmlTemplates *template.Template
		funcMap       template.FuncMap
	}
)

func Setup() *Artefak {
	artefak := &Artefak{router: NewRouter()}
	artefak.RouterGroup = &RouterGroup{artefak: artefak}
	artefak.groups = []*RouterGroup{artefak.RouterGroup}

	return artefak
}

func New() *Artefak {
	artefak := Setup()
	artefak.Use(Logger(), Recover())

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

func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))

	return func(c *Ctx) {
		file := c.Param("filepath")

		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")

	group.GET(urlPattern, handler)
}

func (artefak *Artefak) SetFuncMap(funcMap template.FuncMap) {
	artefak.funcMap = funcMap
}

func (artefak *Artefak) LoadHTMLGlob(pattern string) {
	artefak.htmlTemplates = template.Must(template.New("").Funcs(artefak.funcMap).ParseGlob(pattern))
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
	c.artefak  = artefak
	artefak.router.handle(c)
}
