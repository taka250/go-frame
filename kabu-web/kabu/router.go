package kabu

import (
	"net/http"
	"strings"
)

// 路由表结构

//引擎，要结合trie来实现,roots存放的是不同请求
type Router struct {
	roots  map[string]*node
	router map[string]HandlerFunc
}

//构造新的路由
func newRouter() *Router {
	return &Router{
		router: make(map[string]HandlerFunc),
		roots:  make(map[string]*node),
	}
}

// 分析模式,根据‘/’分割了字符串，并在*之前返回字符串数组,最多一个*
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

//增加路由,同时插入对应方法的节点
func (r *Router) addRoute(method string, pattern string, handler HandlerFunc) {

	parts := parsePattern(pattern)
	_, ok := r.roots[method] //查看是否已经存在有关method的路由
	if !ok {
		r.roots[method] = &node{} //不存在就新建一个节点
	}
	r.roots[method].insert(pattern, parts, 0) //开始插入
	r.router[pattern] = handler
}

//
func (r *Router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path) //实际访问的path
	params := make(map[string]string) //返回参数表
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	} //

	n := root.search(searchParts, 0)
	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts { //查找是否有模糊匹配，并赋值给变量返回
			if parts[0] == ":" {
				params[part[1:]] = searchParts[index]
			}
			if parts[0] == "*" && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/") //赋值给变量一个目录值
				break
			}

		}
		return n, params

	}
	return nil, nil

}

//
func (r *Router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params //变量赋值
		c.handlers = append(c.handlers, r.router[n.pattern])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}
