package kabu

import (
	"log"
	"net/http"
)

//定义函数类型
type HandlerFunc func(c *Context)

type (
	Engine struct {
		Router       *Router
		*RouterGroup                //匿名结构体
		groups       []*RouterGroup // 存储所有的groups
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
	engine.groups = []*RouterGroup{engine.RouterGroup} //每个都映射到引擎
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
	c := newContext(w, req)
	engine.Router.handle(c)

}
