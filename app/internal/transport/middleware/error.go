package middleware

import (
	"cc/app/internal/apperror"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

func Error() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, err := range c.Errors {
			switch {
			case errors.Is(err.Err, apperror.Internal):
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Err})
			case errors.Is(err.Err, apperror.NotExists):
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Err})
			case errors.Is(err.Err, apperror.AlreadyExists):
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Err})
			case errors.Is(err.Err, apperror.InvalidParams):
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Err})
			case errors.Is(err.Err, apperror.InvalidCredentials):
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Err})
			case errors.Is(err.Err, apperror.Unauthorized):
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Err})
			}
		}
	}
}
