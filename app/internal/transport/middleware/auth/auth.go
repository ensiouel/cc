package auth

import (
	"cc/app/internal/apperror"
	"cc/app/internal/domain"
	"cc/app/internal/service"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
)

func Middleware(authService service.AuthService) gin.HandlerFunc {
	const prefix = "Bearer "

	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		if authorization == "" || len(authorization) < len(prefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": apperror.ErrUnauthorized,
			})
			return
		}

		payload := authorization[len(prefix):]

		token, err := authService.ParseToken(payload)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": apperror.ErrUnauthorized.SetMessage("access token has expired"),
				})
				return
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": apperror.ErrUnauthorized.SetMessage("invalid access token"),
			})
			return
		}

		claims, ok := token.Claims.(*domain.Claims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": apperror.ErrUnauthorized,
			})
			return
		}

		c.Set("user_id", claims.UserID)

		c.Next()
	}
}
