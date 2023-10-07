package goo

import (
	"html/template"
	"net/http"
	"path"
)

type HandlerFunc func(*Context)

type RouterGroup struct {
	basePath string
	handlers []HandlerFunc
	engine   *Engine
}

type Engine struct {
	*RouterGroup
	router        *router
	htmlTemplates *template.Template
	funcMap       template.FuncMap
}

func New() *Engine {
	engine := &Engine{router: NewRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	return engine
}

func (group *RouterGroup) Group(path string) *RouterGroup {
	return &RouterGroup{
		basePath: group.basePath + path,
		handlers: group.handlers,
		engine:   group.engine,
	}
}

func (group *RouterGroup) GET(path string, handlers ...HandlerFunc) {
	group.addRouter("GET", path, handlers)
}

func (group *RouterGroup) POST(path string, handlers ...HandlerFunc) {
	group.addRouter("POST", path, handlers)
}

func (group *RouterGroup) addRouter(mathod string, relativePath string, partHandlers []HandlerFunc) {
	path := group.basePath + relativePath
	handlers := append(group.handlers, partHandlers...)
	group.engine.router.addRouter(mathod, path, handlers)
}

func (group *RouterGroup) Use(handlers ...HandlerFunc) {
	group.handlers = append(group.handlers, handlers...)
}

func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.basePath, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("*filepath")
		if _, err := fs.Open(file); err != nil {
			c.Status(404)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPath := path.Join(relativePath, "*filepath")
	group.GET(urlPath, handler)
}

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(path string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(path))
}

func (engine *Engine) Run(post string) (err error) {
	err = http.ListenAndServe(post, engine)
	return
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := NewContext(w, r)
	c.engine = engine
	engine.router.handler(c)
}

func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recover())
	return engine
}
