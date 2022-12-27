package auth

import (
	"cc/app/internal/apperror"
	"cc/app/internal/domain"
	"cc/app/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Middleware(authService service.AuthService) gin.HandlerFunc {
	const prefix = "Bearer "

	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		if authorization == "" || len(authorization) < len(prefix) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		payload := authorization[len(prefix):]

		token, err := authService.ParseToken(payload)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*domain.Claims)
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !claims.VerifyExpiresAt(time.Now(), true) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": apperror.ErrUnauthorized,
			})
			return
		}

		c.Set("user_id", claims.UserID)

		c.Next()
	}
}
