package middlewares

import (
	"log/slog"
	"net/http"

	"opencw/handlers/v1/common"
	"opencw/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func LoadUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID")

		var user models.User
		if err := db.Take(&user, userID).Error; err != nil {
			slog.Warn("LoadUser failed: user not found", "user_id", userID, "err", err)
			c.JSON(http.StatusUnauthorized, common.ErrorResponse{Error: "User not found"})
			c.Abort()
			return
		}

		c.Set("user", &user)
		c.Next()
	}
}
