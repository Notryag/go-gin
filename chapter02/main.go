package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func main() {
	r := gin.Default()

	// 路由分组
	api := r.Group("/api")

	// 定义路由分组中的路由
	api.GET("/users", logger(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "success",
		})
	})

	api.POST("/users", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "success",
		})
	})

	api.PUT("/users", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "success",
		})
	})

	api.DELETE("/users", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "success",
		})
	})
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		log.Printf("[%s] %s\n", t.Format("2006-01-02 15:04:05"), c.Request.URL.Path)
		c.Next()
	}
}
