package handler

import (
	"chapter05/config"
	"chapter05/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Handlers struct {
	DB *gorm.DB
}

func (h *Handlers) SignUp(c *gin.Context) {

	var user model.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var exitingUser model.User

	if err := h.DB.Where("user_name = ?", user.UserName).First(&exitingUser).Error; err == nil {
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
	if err := h.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//db.Create(&model.User{UserName: user.UserName, Password: user.Password})
	c.JSON(http.StatusOK, user)
}

func (h *Handlers) SinIn(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var exitingUser model.User
	if err := h.DB.Where("user_name = ?", user.UserName).First(&exitingUser).Error; err != nil {
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

	config.Init()

	signedToken, err := token.SignedString([]byte(config.Cfg.SecretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"token": signedToken})
}
