package main

import (
	"os"

	routes "github.com/dilshad-cp/golang-jwt-project/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRouter(router)
	routes.UserRouter(router)

	router.GET("api-1", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"success": "Access Granted for api-1"})
	})

	router.GET("api-2", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"success": "Access Granted for api-2"})
	})

	router.Run(":" + port)
}
