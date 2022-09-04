package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.tmpl")
	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
	})

	group := r.Group("/config")
	group.POST("/create", create)
	group.GET("/list", list)
	group.GET("/get", get)
	// delete config by server name
	group.DELETE("/del/:server_name", del)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
