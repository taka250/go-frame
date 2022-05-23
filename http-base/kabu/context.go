package kabu

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	StatusCode int
	Params     map[string]string //前缀树的增加，模糊匹配导致增加参数

	handlers []HandlerFunc //用于中间件
	index    int
	engine   *Engine //引擎指针
}

//初始化
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

func (c *Context) Next() { //下一个中间件
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

//表单数据
func (c *Context) Postform(key string) string {
	return c.Req.FormValue(key)
}

//查询参数,返回的是第一个值
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

//状态码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

//头消息
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

//fail
func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})

}

//字符串 ...是可变参数
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

//json形式
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	_, err = c.Writer.Write(jsonBytes)
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

//html
func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Fail(500, err.Error())
	}

}

func (c *Context) Param(key string) string { //返回参数
	value := c.Params[key]
	return value
}
