package ginutils

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetUUID(c *gin.Context, key string) (u uuid.UUID) {
	if val, ok := c.Get(key); ok && val != nil {
		u, _ = val.(uuid.UUID)
	}

	return
}
