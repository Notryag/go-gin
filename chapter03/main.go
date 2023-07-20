package chapter03

import (
	"chapter03/posts"
	"chapter03/users"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	r := gin.Default()
	// users 路由
	users.Router(r)
	// posts 路由
	posts.Router(r)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
