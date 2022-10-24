package routes

import (
	controller "github.com/Jayleonc/go-stage/controllers"
	middleware "github.com/Jayleonc/go-stage/middlewares"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.GET("/users/:user_id", controller.GetUser())
	incomingRoutes.PUT("/users/:user_id", controller.UpdateUser())
	incomingRoutes.DELETE("/users/:user_id", controller.DeleteUser())
}
