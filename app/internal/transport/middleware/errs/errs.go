package errs

import (
	"cc/app/internal/apperror"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, err := range c.Errors {
			switch {
			case errors.Is(err.Err, apperror.ErrInternalError):
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Err})
			case errors.Is(err.Err, apperror.ErrNotExists):
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Err})
			case errors.Is(err.Err, apperror.ErrAlreadyExists):
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Err})
			case errors.Is(err.Err, apperror.ErrInvalidParams):
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Err})
			case errors.Is(err.Err, apperror.ErrInvalidCredentials):
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Err})
			case errors.Is(err.Err, apperror.ErrUnauthorized):
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Err})
			default:
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": apperror.ErrUnknownError.SetError(err.Err)})
			}
		}
	}
}
