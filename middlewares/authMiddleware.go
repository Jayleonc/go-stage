package middlewares

import (
	"fmt"
	helper "github.com/Jayleonc/go-stage/helpers"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No Authorization header provided")})
			c.Abort()
			return
		}

		claims, err := helper.ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"ValidateToken error": err})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("user_name", claims.UserName)
		c.Set("user_id", claims.UserId)
		c.Set("user_type", claims.UserType)
		c.Next()
	}
}
