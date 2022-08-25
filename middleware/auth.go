package middleware

import (
	"net/http"

	helper "github.com/dilshad-cp/golang-jwt-project/helpers"

	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientToken := ctx.Request.Header.Get("token")

		if clientToken == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Auth token missing"})
			ctx.Abort()
			return
		}
		claims, err := helper.ValidateToken(clientToken)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Auth error"})
			ctx.Abort()
			return
		}
		ctx.Set("emails", claims.Email)
		ctx.Set("first_name", claims.First_name)
		ctx.Set("last_name", claims.Last_name)
		ctx.Set("uid", claims.Uid)
		ctx.Set("user_type", claims.User_type)
		ctx.Set("user_type", claims.User_type)
		ctx.Next()
	}
}
