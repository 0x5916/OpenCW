package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	common2 "opencw/common"
	"opencw/models"
	"opencw/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB *gorm.DB
}

func (h UserHandler) GetUserInfo(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	c.JSON(http.StatusOK, common2.UserInfoResponse{
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	})
}

func (h UserHandler) GetOtherUserInfo(c *gin.Context) {
	otherUserId := c.Param("id")

	var otherUser models.User
	if err := h.DB.Take(&otherUser, "id = ?", otherUserId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, common2.ErrorResponse{Error: "User not found"})
		} else {
			slog.Error("GetOtherUserInfo DB error", "id", otherUserId, "err", err)
			c.JSON(http.StatusInternalServerError, common2.ErrorResponse{Error: "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, common2.UserInfoResponse{
		Username:  otherUser.Username,
		CreatedAt: otherUser.CreatedAt,
	})
}

func (h UserHandler) UpdateEmail(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	var input common2.UpdateEmailInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, common2.ErrorResponse{Error: "Invalid request body"})
		return
	}
	if user.Email == input.Email {
		c.JSON(http.StatusBadRequest, common2.ErrorResponse{Error: "New email must be different from current email"})
		return
	}

	if err := h.DB.Model(user).Update("email", input.Email).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c.JSON(http.StatusConflict, common2.ErrorResponse{Error: "Email already in use"})
			return
		}
		slog.Error("Failed to update email", "user_id", user.ID, "new_email", input.Email, "err", err)
		c.JSON(http.StatusInternalServerError, common2.ErrorResponse{Error: "Failed to update email"})
		return
	}

	c.JSON(http.StatusOK, common2.MessageResponse{Message: "Email updated"})
}

func (h UserHandler) UpdatePassword(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	var input common2.UpdatePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, common2.ErrorResponse{Error: "Invalid request body"})
		return
	}

	match, err := utils.ComparePasswordAndHash(input.OldPassword, user.Password)
	if err != nil {
		slog.Error("Failed to compare password hash", "user_id", user.ID, "err", err)
		c.JSON(http.StatusInternalServerError, common2.ErrorResponse{Error: "internal error"})
		return
	}
	if !match {
		slog.Warn("Update failed: invalid password", "user_id", user.ID, "ip", c.ClientIP())
		c.JSON(http.StatusUnauthorized, common2.ErrorResponse{Error: "Invalid credentials"})
		return
	}

	hash, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		slog.Error("Failed to hash password", "err", err)
		c.JSON(http.StatusInternalServerError, common2.ErrorResponse{Error: "Failed to hash password"})
		return
	}

	if err := h.DB.Model(user).Update("password", hash).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c.JSON(http.StatusConflict, common2.ErrorResponse{Error: "Email already in use"})
			return
		}
		slog.Error("Failed to update password", "user_id", user.ID, "err", err)
		c.JSON(http.StatusInternalServerError, common2.ErrorResponse{Error: "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, common2.MessageResponse{Message: "Password updated"})
}
