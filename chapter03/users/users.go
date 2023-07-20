package users

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		//	简单记录日志
		fmt.Println("start")
		c.Next()
	}
}

// Router 需要大写
func Router(r *gin.Engine) {
	users := r.Group("/users")
	users.Use(Logger())

	users.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "success",
		})
	})
	users.POST("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "success",
		})
	})
	users.DELETE("/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "success",
		})
	})
	users.PUT("/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "success",
		})
	})
}
