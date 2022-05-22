package main

import (
	"kabu"
	"net/http"
)

func main() {
	r := kabu.New()                         //返回一个引擎
	r.GET("/index", func(c *kabu.Context) { //继承了方法
		c.HTML(http.StatusOK, "<h1>Index Page</h1>")
	})
	v1 := r.Group("/v1") //
	{

		v1.GET("/", func(c *kabu.Context) {
			c.HTML(http.StatusOK, "<h1>Hello Gee1111</h1>")
		})

		v1.GET("/hello", func(c *kabu.Context) {
			// expect /hello?name=geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}
	v2 := r.Group("/v2")
	{
		v2.GET("/hello/:name", func(c *kabu.Context) {
			// expect /hello/geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v2.POST("/login", func(c *kabu.Context) {
			c.JSON(http.StatusOK, kabu.H{
				"username": c.Postform("username"),
				"password": c.Postform("password"),
			})
		})

	}

	r.Run(":9999")
}
