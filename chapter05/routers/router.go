package routers

import (
	"chapter05/database"
	"chapter05/handler"
	"chapter05/middleware"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Server interface {
	Run(addr ...string) (err error)
}

func setupRouter() *gin.Engine {
	db, err := database.InitDB()
	if err != nil {
		log.Println("failed to connect database")
	}
	h := handler.Handlers{DB: db}

	r := gin.Default()
	r.POST("/signup", h.SignUp)

	r.POST("/signin", h.SinIn)

	authorized := r.Group("/api")
	authorized.Use(middleware.Auth())
	{
		authorized.POST("/todos", h.CreateTodo)

		authorized.GET("/todos", h.GetTodos)

		authorized.PUT("/todos/:id", h.UpdateTodo)

		authorized.DELETE("/todos/:id", h.DeleteTodoById)

		authorized.POST("/signout", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"msg": "success"})
		})

	}
	return r
}

func NewServer() Server {
	return setupRouter()
}
