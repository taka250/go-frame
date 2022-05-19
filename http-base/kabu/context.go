package kabu

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	Writer http.ResponseWriter
	Req    *http.Request
	Path   string
	Method string

	StatusCode int
}

//初始化
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
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
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}
