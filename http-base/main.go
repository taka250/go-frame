package main

import (
	"kabu"
	"net/http"
)

func main() {
	r := kabu.New()
	r.GET("/", func(c *kabu.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})
	r.GET("/hello", func(c *kabu.Context) {
		// expect /hello?name=geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.POST("/login", func(c *kabu.Context) {
		c.JSON(http.StatusOK, kabu.H{
			"username": c.Postform("username"),
			"password": c.Postform("password"),
		})
	})

	r.Run(":9999")
}
