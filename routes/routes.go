package routes

import (
	"golang-todos/controllers"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	api := router.Group("/api")
	{
		api.POST("/login", controllers.Login)
		api.POST("/refresh-token", controllers.HandleRefresh)
		api.POST("/user/register", controllers.RegisterUser)
	}
	return router
}
