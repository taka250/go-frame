package main

import (
	"kabu"
	"net/http"
)

func main() {
	r := kabu.New()
	r.GET("/", func(c *kabu.Context) {
		c.HTML(http.StatusOK, "<h1>Hello kabu</h1>")
	})

	r.GET("/hello", func(c *kabu.Context) {
		// expect /hello?name=geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.GET("/hello/:name", func(c *kabu.Context) {
		// expect /hello/geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	r.GET("/assets/*filepath", func(c *kabu.Context) {
		c.JSON(http.StatusOK, kabu.H{"filepath": c.Param("filepath")})
	})

	r.Run(":9999")
}
