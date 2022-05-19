package kabu

import (
	"log"
	"net/http"
)

// 路由表结构
type router struct {
	method  string
	handler HandlerFunc
}

//引擎
type Router struct {
	router map[string]router
}

//构造新的路由
func newRouter() *Router {
	return &Router{router: make(map[string]router)}
}

//增加路由
func (r *Router) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	r.router[pattern] = router{method, handler}
}

func (r *Router) handle(c *Context) {
	if router, ok := r.router[c.Path]; ok && router.method == c.Method {
		router.handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
