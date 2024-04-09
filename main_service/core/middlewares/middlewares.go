package middlewares

import (
	jwt_utils "main_service_core/jwt_utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	token := c.Request.Header.Get("X-API-Key")

	id, err := jwt_utils.GetStrIdFromJWT(token)
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{"error": "API key is missing or invalid"})
		return
	}

	c.Set("id", id)
	c.Next()
}
