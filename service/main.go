package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tryag/common"
	"log"
)

func main()  {
	fmt.Println(common.Add(1, 3))
	r := gin.Default()
	r.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"name": "hello world",
		})
	})
	err := r.Run(":8080")
	if err != nil {
		log.Fatalln(err)
	}
}