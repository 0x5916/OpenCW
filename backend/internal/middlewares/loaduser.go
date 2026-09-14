package middlewares

import (
	"log/slog"
	"net/http"
	"opencw/internal/common"
	"opencw/internal/models"
	"uuid"

	"github.com/gin-gonic/gin"
	googleuuid "github.com/google/uuid"
	"gorm.io/gorm"
)

func LoadUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		gormUserID, err := googleuuid.Parse(userID.String())
		if err != nil {
			c.JSON(http.StatusUnauthorized, common.NewErrorResponse(common.ErrorCodeInvalidToken, "Invalid token"))
			c.Abort()
			return
		}

		var user models.User
		if err := db.Take(&user, gormUserID).Error; err != nil {
			slog.Warn("LoadUser failed: user not found", "user_id", userID, "err", err)
			c.JSON(http.StatusUnauthorized, common.NewErrorResponse(common.ErrorCodeUserNotFound, "User not found"))
			c.Abort()
			return
		}

		c.Set("user", &user)
		c.Next()
	}
}
