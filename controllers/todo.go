package controllers

import (
	"golang-todos/database"
	"golang-todos/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTodos(context *gin.Context) {
	var todo []models.Todo

	record := database.Instance.Find(&todo)
	if record.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": todo})
}

func InsertTodo(context *gin.Context) {
	var todo models.TodoCreation

	if err := context.ShouldBindJSON(&todo); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	record := database.Instance.Create(&todo)

	if record.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": record.Error.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": todo})
}

func UpdateTodo(context *gin.Context) {
	var todo models.TodoUpdate
	id := context.Param("id")

	if err := context.ShouldBindBodyWithJSON(&todo); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.Instance.Where("id = ?", id).Updates(&todo).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": todo})

}
