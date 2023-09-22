package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	g1 := r.Group("/test1", func(c *gin.Context) {
		fmt.Println(1)
	})
	g1.GET("b//a/", func(c *gin.Context) {
		c.String(200, "1")
	})
	g2 := r.Group("/test1", func(c *gin.Context) {
		fmt.Println(2)
	})
	g2.GET("/1", func(c *gin.Context) {
		c.String(200, "2")
	})
	// 启动Gin服务
	r.Run(":8080")
}
