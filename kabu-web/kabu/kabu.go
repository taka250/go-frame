package kabu

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

//定义函数类型
type HandlerFunc func(c *Context)

type (
	Engine struct {
		Router        *Router
		*RouterGroup                     //匿名结构体 继承方法
		groups        []*RouterGroup     // 存储所有的groups
		htmlTemplates *template.Template //渲染html模板
		funcMap       template.FuncMap
	}

	RouterGroup struct {
		Prefix      string //前缀
		middlewares []HandlerFunc
		parent      *RouterGroup //父类group
		engine      *Engine      //所有的groups 享有一个引擎
	}
)

//New函数是构造器
func New() *Engine {
	engine := &Engine{Router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup} //这里将engine的第一个Routergroup加入到切片中。
	return engine
}

//构建默认engine

func Default() *Engine {
	engine := New()
	engine.Use(Recovery())
	return engine
}

//构造一个新的group  其实是子group
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		Prefix: group.Prefix + prefix,
		engine: engine,
		parent: group,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (engine *Engine) SetFuncMap(funcmap template.FuncMap) {
	engine.funcMap = funcmap
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

//增加一个路由
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.Prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.Router.addRoute(method, pattern, handler)
}

//get
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

//post
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

//开启监听
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

//这里的比较用的是map查询
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.Prefix) { //遍历所有的group
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	c.engine = engine
	engine.Router.handle(c)

}

//中间件代码实现

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

//以下是静态文件
//创建静态handler
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.Prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("fliepath")
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

//serve files
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath") //添加一个含有模糊查找的路径
	//register get handlers
	group.GET(urlPattern, handler)
}
