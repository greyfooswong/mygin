package main

import (
	"mygin"
	"net/http"
)

func main() {
	r := mygin.New()
	r.GET("/", func(c *mygin.Context) {
		c.HTML(http.StatusOK, "<h1>Hello myGin </h1>")
	})
	r.GET("/hello", func(c *mygin.Context) {
		c.String(http.StatusOK, "hello %s, you are at %s \n", c.Query("name"), c.Path)
	})
	r.POST("/login", func(c *mygin.Context) {
		c.JSON(http.StatusOK, mygin.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})
	r.Run(":9999")
}
