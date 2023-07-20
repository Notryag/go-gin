package posts

import (
	"chapter03/users"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Router(r *gin.Engine) {
	posts := r.Group("/posts")
	posts.Use(users.Logger())

	posts.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "success",
		})
	})
	posts.POST("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "success",
		})
	})
	posts.DELETE("/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "success",
		})
	})
	posts.PUT("/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "success",
		})
	})
}
