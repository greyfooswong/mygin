package mygin

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

type HandlerFunc func(c *Context)

type Engine struct {
	*RouterGroup
	router        *router
	groups        []*RouterGroup
	htmlTemplates *template.Template
	funcMap       template.FuncMap
}

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc
	parent      *RouterGroup
	engine      *Engine
}

func New() *Engine {
	engine := &Engine{
		router: newRouter(),
	}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method, comp string, handlerFunc HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Router %4s _ %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handlerFunc)
}

func (group *RouterGroup) GET(pattern string, handlerFunc HandlerFunc) {
	group.addRoute("GET", pattern, handlerFunc)
}

func (group *RouterGroup) POST(pattern string, handlerFunc HandlerFunc) {
	group.addRoute("POST", pattern, handlerFunc)
}

func (engine *Engine) Run(add string) (err error) {
	return http.ListenAndServe(add, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.engine = engine
	c.handlers = middlewares
	engine.router.handle(c)
}

func (group *RouterGroup) createStaticHeanlder(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}
func (group *RouterGroup) Static(relative, root string) {
	handler := group.createStaticHeanlder(relative, http.Dir(root))
	urlPattern := path.Join(relative, "/*filepath")
	group.GET(urlPattern, handler)
}
