package routes

import (
	"github.com/dilshad-cp/golang-jwt-project/middleware"

	controller "github.com/dilshad-cp/golang-jwt-project/controller"

	"github.com/gin-gonic/gin"
)

func UserRouter(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.GET("/users/:user_id", controller.GetUser())
}
