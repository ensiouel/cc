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
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Err})
			case errors.Is(err.Err, apperror.ErrNotExists):
				c.JSON(http.StatusNotFound, gin.H{"error": err.Err})
			case errors.Is(err.Err, apperror.ErrAlreadyExists):
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Err})
			case errors.Is(err.Err, apperror.ErrInvalidParams):
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Err})
			case errors.Is(err.Err, apperror.ErrInvalidCredentials):
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Err})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": apperror.ErrUnknownError.SetError(err.Err)})
			}
		}
	}
}
