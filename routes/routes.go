package routes

import (
	"golang-todos/controllers"
	"golang-todos/middlewares"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	api := router.Group("/api")
	{
		api.POST("/login", controllers.Login)
		api.POST("/refresh-token", controllers.HandleRefresh)
		api.POST("/user/register", controllers.RegisterUser)

		secured := api.Group("/").Use(middlewares.Auth())
		{
			secured.GET("/todos", controllers.GetTodos)
			secured.POST("/todos/create", controllers.InsertTodo)
			secured.PUT("/todos/:id", controllers.UpdateTodo)
		}
	}
	return router
}
