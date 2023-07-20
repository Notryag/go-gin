package main

import (
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 中间件
func myHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		// 只要调用了这个中间件都可以拿到这个参数
		context.Set("userSession", "userid-1")
		context.Next()
	}
}

type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(20);not null"`
	Email    string `gorm:"type:varchar(20);not null;unique"`
	Password string `gorm:"size:255"`
}

func InitDB() *gorm.DB {
	host := "localhost"
	port := "3306"
	userName := "root"
	password := "root"
	database := "test"
	charset := "utf8"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", userName, password, host, port, database, charset)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database, err:" + err.Error())
	}

	// 自动迁移
	err = db.AutoMigrate(&User{})
	if err != nil {
		panic("failed to migrate database, err:" + err.Error())
	}

	return db
}

func main() {
	InitDB()
	r := gin.Default()

	r.POST("/api/auth/register", func(ctx *gin.Context) {
		name := ctx.PostForm("name")
		telephone := ctx.PostForm("telephone")
		password := ctx.PostForm("password")

		if len(telephone) != 11 {
			ctx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
				"code": 422,
				"msg":  "手机号必须为11位",
			})
			return
		}

		if len(password) < 6 {
			ctx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
				"code": 422,
				"msg":  "密码不能少于6位",
			})
			return
		}
		if len(name) == 0 {
			name = RandomString(10)
		}
		log.Println(name, telephone, password)

		ctx.JSON(http.StatusOK, gin.H{
			"msg": "success",
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

func RandomString(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var letters = []byte("asdfghjklzxcvbnmqwertyuiopASDFGHJKLZXCVBNMQWERTYUIOP")
	result := make([]byte, n)
	for i := range result {
		result[i] = letters[r.Intn(len(letters))]
	}
	return string(result)
}
