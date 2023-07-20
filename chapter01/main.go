package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(20);not null"`
	Email    string `gorm:"type:varchar(20);not null;unique"`
	Password string `json:"-"`
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
	db := InitDB()
	r := gin.Default()

	r.GET("/users", func(c *gin.Context) {
		var users []User
		if err := db.Find(&users).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "查询失败",
			})
			log.Println(err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": users,
		})
	})

	r.POST("/users", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "参数错误",
			})
			log.Println(err)
			return
		}
		if err := db.Create(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  err.Error(),
			})
			log.Println(err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": user,
		})
	})

	r.PUT("/users/:id", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "参数错误",
			})
			return
		}
		id := c.Param("id")
		if err := db.Model(&user).Where("id = ?", id).Updates(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "更新失败",
			})
			log.Println(err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": user,
		})
	})

	r.DELETE("/users/:id", func(c *gin.Context) {
		var user User
		id := c.Param("id")
		if err := db.Where("id = ?", id).Delete(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "删除失败",
			})
			log.Println(err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": "删除成功",
		})
	})

	err := r.Run(":8080")
	if err != nil {
		log.Fatalln(err)
	}
}
