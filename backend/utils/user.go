package utils

import (
	"opencw/models"

	"github.com/gin-gonic/gin"
)

func MustGetUser(c *gin.Context) *models.User {
	return c.MustGet("user").(*models.User)
}
