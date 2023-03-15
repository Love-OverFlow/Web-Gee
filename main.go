package main

import (
	"Web-Gee/gee"
	"fmt"
	"net/http"
)

func main() {
	r := gee.New()
	v1 := r.Group("/v1")

	v1.Use(func(c *gee.Context) {
		fmt.Println("预处理中间件操作")
		c.Next()
		fmt.Println("善后中间件操作")
	})

	v1.GET("/abc", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello World</h1>")
	})

	v1.GET("/login", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	v1.GET("/hello", func(c *gee.Context) {
		// 如 /hello?name=gee
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	v1.GET("/hello/:name", func(c *gee.Context) {
		// 如 /hello/gee
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	r.Run(":9998")
}
