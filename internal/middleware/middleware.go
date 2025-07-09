package middleware

import (
	"net/http"
	"strings"

	"golang/stockLkBack/internal/utils/jwtgen"

	"github.com/gin-gonic/gin"
)

// TokenAuthMiddleware Middleware для проверки JWT токена.
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем токен из заголовка Authorization
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			c.Abort()

			return
		}

		// Удаляем Bearer из токена
		tokenString = strings.ReplaceAll(tokenString, "Bearer ", "")

		// Проверяем и парсим токен
		claims, err := jwtgen.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()

			return
		}

		// Если токен валиден, добавляем пользователя в контекст запроса
		c.Set("login", claims.Login)
		c.Set("Role", claims.Role)

		c.Next()
	}
}
