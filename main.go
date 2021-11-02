package main

import (
	"mygin"
	"net/http"
)

func main() {
	r := mygin.New()
	r.GET("/index", func(c *mygin.Context) {
		c.HTML(http.StatusOK, "<h1>index page </h1>")
	})
	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *mygin.Context) {
			c.HTML(http.StatusOK, "<h1>Hello myGin </h1>")
		})
		v1.GET("/hello", func(c *mygin.Context) {
			c.String(http.StatusOK, "hello %s, you are at %s \n", c.Query("name"), c.Path)
		})
	}
	v2 := r.Group("/v2")
	{
		v2.GET("/hello/:name", func(c *mygin.Context) {
			c.String(http.StatusOK, "hello %s, you are at %s \n", c.Param("name"), c.Path)
		})
		v2.GET("/assets/*filepath", func(c *mygin.Context) {
			c.JSON(http.StatusOK, mygin.H{"filepath": c.Param("filepath")})
		})
	}
	/*	r.POST("/login", func(c *mygin.Context) {
		c.JSON(http.StatusOK, mygin.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})*/
	r.Run(":9999")
}
