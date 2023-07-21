package main

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

type User struct {
	gorm.Model
	ID       int    `gorm:"primary_key" json:"id"`
	UserName string `gorm:"type:varchar(20);not null" json:"userName"`
	Password string `gorm:"type:varchar(220);not null" json:"password"`
}

type Todo struct {
	gorm.Model
	ID      int    `gorm:"primaryKey" json:"id"`
	Title   string `gorm:"type:varchar(20);not null" json:"title"`
	Status  string `gorm:"type:varchar(20);" json:"status"`
	Content string `gorm:"type:varchar(20);" json:"content"`
	UserID  int    `json:"userId"`
}

const secretKey = "secret"

func initDb() *gorm.DB {
	dsn := "root:root@tcp(localhost:3306)/chapter04?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	err = db.AutoMigrate(&User{}, &Todo{})
	if err != nil {
		panic("failed to migrate database")
	}

	// 关闭数据库连接
	//sqlDB, err := db.DB()
	//if err != nil {
	//	panic("failed to get DB instance")
	//}
	//defer sqlDB.Close()
	return db
}

func main() {
	db := initDb()
	r := gin.Default()

	r.POST("/signup", func(c *gin.Context) {
		signUp(c, db)
	})

	r.POST("/signin", func(c *gin.Context) {
		sinIn(c, db)
	})

	authorized := r.Group("/api")
	authorized.Use(authenticationMiddleware())
	{
		authorized.POST("/todos", func(c *gin.Context) {
			createTodo(c, db)
		})

		authorized.GET("/todos", func(c *gin.Context) {
			getTodos(c, db)
		})

		authorized.PUT("/todos/:id", func(c *gin.Context) {
			updateTodo(c, db)
		})

		authorized.DELETE("/todos/:id", func(c *gin.Context) {
			deleteTodoById(c, db)
		})

		authorized.POST("/signout", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"msg": "success"})
		})

	}

	if err := r.Run(":8081"); err != nil {
		panic(err.Error())
	}
}

func signUp(c *gin.Context, db *gorm.DB) {

	var user User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var exitingUser User
	//if err := db.Where("username = ?", user.UserName).First(&exitingUser).Error; err == nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
	//	return
	//}
	if err := db.Where("user_name = ?", user.UserName).First(&exitingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
		return
	}
	log.Println("exitingUser", exitingUser.UserName)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user.Password = string(hashedPassword)
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//db.Create(&User{UserName: user.UserName, Password: user.Password})
	c.JSON(http.StatusOK, user)
}

func sinIn(c *gin.Context, db *gorm.DB) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var exitingUser User
	if err := db.Where("user_name = ?", user.UserName).First(&exitingUser).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username does not exist"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(exitingUser.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       exitingUser.ID,
		"username": exitingUser.UserName,
		"password": exitingUser.Password,
	})

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"token": signedToken})
}

func authenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("token")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header not provider"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("userID", int(claims["id"].(float64)))
	}
}

func createTodo(c *gin.Context, db *gorm.DB) {
	var todo Todo

	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo.UserID = c.GetInt("userID")
	db.Create(&todo)
	c.JSON(http.StatusOK, gin.H{"msg": "success"})
}

func getTodos(c *gin.Context, db *gorm.DB) {
	var todos []Todo
	db.Where("user_id = ?", c.GetInt("userID")).Find(&todos)
	c.JSON(http.StatusOK, gin.H{"data": todos})
}

func updateTodo(c *gin.Context, db *gorm.DB) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updateTodo Todo
	if err := c.ShouldBindJSON(&updateTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 判断是否存在
	var todo Todo
	if err := db.Where("id = ?", id).First(&todo).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Model(&Todo{}).Where("id = ?", id).Updates(updateTodo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updateTodo)

}

func deleteTodoById(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	var todo Todo
	if err := db.Where("id =?", id).First(&todo).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Delete(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "success"})

}

//func deleteTodo(c *gin.Context, db *gorm.DB)  {
//	id := c.Param("id")
//
//
//}
