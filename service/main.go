package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tryag/common"
	"log"
	"net/http"
)

// 中间件
func myHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		// 只要调用了这个中间件都可以拿到这个参数
		context.Set("userSession", "userid-1")
		context.Next()
	}
}

func main() {
	fmt.Println(common.Add(1, 3))
	r := gin.Default()
	r.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"name": "hello world",
		})
	})

	r.GET("/user/info", myHandler(), func(ctx *gin.Context) {
		userSession := ctx.MustGet("userSession").(string)
		log.Println(userSession)
		userId := ctx.Query("userId")
		userName := ctx.Query("userName")
		ctx.JSON(http.StatusOK, gin.H{
			"userId":   userId,
			"userName": userName,
		})
	})

	r.GET("/user/info/:userId/:userName", func(ctx *gin.Context) {
		userId := ctx.Param("userId")
		userName := ctx.Param("userName")
		ctx.JSON(http.StatusOK, gin.H{
			"userId":   userId,
			"userName": userName,
		})
	})

	// 获取json参数
	r.POST("/user/json", func(ctx *gin.Context) {
		data, _ := ctx.GetRawData()
		fmt.Println(data)
		var m map[string]interface{}
		_ = json.Unmarshal(data, &m)
		ctx.JSON(http.StatusOK, m)

	})

	// 获取form参数
	r.POST("/user/form", func(ctx *gin.Context) {
		userName := ctx.PostForm("userName")
		userId := ctx.PostForm("userId")
		ctx.JSON(http.StatusOK, gin.H{
			"msg":      "ok",
			"userName": userName,
			"userId":   userId,
		})
	})

	r.GET("/moreJSON", func(c *gin.Context) {
		// 你也可以使用一个结构体
		var msg struct {
			Name    string `json:"user"`
			Message string
			Number  int
		}
		msg.Name = "Lena"
		msg.Message = "hey"
		msg.Number = 123
		// 注意 msg.Name 在 JSON 中变成了 "user"
		// 将输出：{"user": "Lena", "Message": "hey", "Number": 123}
		c.JSON(http.StatusOK, msg)
	})

	err := r.Run(":8080")
	if err != nil {
		log.Fatalln(err)
	}
}
