package main

import (
	"kabu"
	"net/http"
)

func main() {
	r := kabu.Default()
	r.GET("/", func(c *kabu.Context) {
		c.String(http.StatusOK, "{\"message\":\"Internal Server Error\"}")
	})
	// index out of range for testing Recovery()
	r.GET("/panic", func(c *kabu.Context) {
		names := []string{"gsssssss"}
		c.String(http.StatusOK, names[100])
	})

	r.Run(":9999")
}
