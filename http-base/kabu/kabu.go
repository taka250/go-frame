package kabu

import (
	"net/http"
)

//定义函数类型
type HandlerFunc func(c *Context)

type Engine struct {
	Router *Router
}

//New函数是构造器
func New() *Engine {
	return &Engine{Router: newRouter()}
}

//以下都为方法
//增加一个路由
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.Router.addRoute(method, pattern, handler)
}

//get
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

//post
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
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
