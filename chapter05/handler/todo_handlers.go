package handler

import (
	"chapter05/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handlers) CreateTodo(c *gin.Context) {
	var todo model.Todo

	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo.UserID = c.GetInt("userID")
	h.DB.Create(&todo)
	c.JSON(http.StatusOK, gin.H{"msg": "success"})
}

func (h *Handlers) GetTodos(c *gin.Context) {
	var todos []model.Todo
	h.DB.Where("user_id = ?", c.GetInt("userID")).Find(&todos)
	c.JSON(http.StatusOK, gin.H{"data": todos})
}

func (h *Handlers) UpdateTodo(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updateTodo model.Todo
	if err := c.ShouldBindJSON(&updateTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 判断是否存在
	var todo model.Todo
	if err := h.DB.Where("id = ?", id).First(&todo).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.DB.Model(&model.Todo{}).Where("id = ?", id).Updates(updateTodo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updateTodo)

}

func (h *Handlers) DeleteTodoById(c *gin.Context) {
	id := c.Param("id")

	var todo model.Todo
	if err := h.DB.Where("id =?", id).First(&todo).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.DB.Delete(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "success"})

}
