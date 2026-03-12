package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"opencw/handlers/v1/common"

	"opencw/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB *gorm.DB
}

func (h UserHandler) GetUserInfo(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	c.JSON(http.StatusOK, common.UserInfoResponse{
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	})
}

func (h UserHandler) UpdateEmail(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	var input common.UpdateEmailInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: err.Error()})
		return
	}
	if user.Email == input.Email {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: "New email must be different from current email"})
		return
	}

	if err := h.DB.Model(user).Update("email", input.Email).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c.JSON(http.StatusConflict, common.ErrorResponse{Error: "Email already in use"})
			return
		}
		slog.Error("Failed to update email", "user_id", user.ID, "new_email", input.Email, "err", err)
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Error: "Failed to update email"})
		return
	}

	c.JSON(http.StatusOK, common.MessageResponse{Message: "Email updated"})
}
