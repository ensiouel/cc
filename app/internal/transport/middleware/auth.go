package middleware

import (
	"cc/app/internal/apperror"
	"cc/app/internal/domain"
	"cc/app/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
)

func Auth(authService service.AuthService) gin.HandlerFunc {
	const prefix = "Bearer "

	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		if authorization == "" || len(authorization) < len(prefix) {
			_ = c.Error(apperror.Unauthorized.WithMessage("invalid access token"))
			c.Abort()
			return
		}

		payload := authorization[len(prefix):]

		token, err := authService.ParseToken(payload)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				_ = c.Error(apperror.Unauthorized.WithMessage("access token has expired"))
				c.Abort()
				return
			}

			_ = c.Error(apperror.Unauthorized.WithMessage("invalid access token"))
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*domain.Claims)
		if !ok {
			_ = c.Error(apperror.Unauthorized.WithMessage("invalid access token"))
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)

		c.Next()
	}
}
