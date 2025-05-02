package middlewares

import (
	"github.com/andro-kes/Blog/utils"
	"github.com/gin-gonic/gin"

	"log"
)

func IsAuthorized() gin.HandlerFunc {
	return func(c *gin.Context) {
		token , err := c.Cookie("token")
		if err != nil {
			log.Println("Middlewares: Не удалось извлечь токен из куки")
			c.AbortWithStatusJSON(500, gin.H{"error": "Пользователь не авторизован"})
			return
		}
		claims, err := utils.ParseToken(token)
		if err != nil {
			log.Println("Middlewares: Не удалось достать Claims")
			c.AbortWithStatusJSON(500, gin.H{"error": "Не удалось извлечь права доступа"})
			return
		}
		c.Set("role", claims.Role)
		c.Next()
	}
}