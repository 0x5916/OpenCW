package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"opencw/handlers/v1/common"
	"opencw/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB *gorm.DB
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input common.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: err.Error()})
		return
	}

	var user models.User
	err := h.DB.Select("username, email").
		Where("username = ? OR email = ?", input.Username, input.Email).
		First(&user).Error
	if err == nil {
		if user.Username == input.Username && user.Email == input.Email {
			c.JSON(http.StatusConflict, common.ErrorResponse{Error: "Username and email already exists"})
		} else if user.Username == input.Username {
			c.JSON(http.StatusConflict, common.ErrorResponse{Error: "Username already exists"})
		} else {
			c.JSON(http.StatusConflict, common.ErrorResponse{Error: "Email already exists"})
		}
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		slog.Error("Failed to query user user", "err", err, "username", input.Username)
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Error: "Database failure"})
		return
	}

	hash, err := common.HashPassword(input.Password)
	if err != nil {
		slog.Error("Failed to hash password", "err", err)
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Error: "Failed to hash password"})
		return
	}

	var result common.HttpErrorResponse
	user = models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: hash,
	}

	err = h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				result = common.HttpErrorResponse{Status: http.StatusConflict, Err: "Registration conflict, please try again"}
			} else {
				result = common.HttpErrorResponse{Status: http.StatusInternalServerError, Err: "Failed to create user"}
			}
			return err
		}
		return nil
	})
	if err != nil {
		if result.Status == http.StatusInternalServerError {
			slog.Error("Registration transaction failed", "err", err, "username", input.Username)
		}
		c.JSON(result.Status, common.ErrorResponse{Error: result.Err})
		return
	}

	rawToken, accessToken, err := common.IssueTokenPair(h.DB, user.ID, time.Now())
	if err != nil {
		slog.Error("Failed to issue token pair", "user_id", user.ID, "err", err)
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Error: "Failed to issue token, try to login."})
		return
	}

	c.JSON(http.StatusOK, common.AuthTokenPairResponse{
		RefreshToken: rawToken,
		AccessToken:  accessToken,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input common.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: err.Error()})
		return
	}

	var user models.User
	if err := h.DB.Take(&user, "username = ?", input.Username).Error; err != nil {
		slog.Warn("Login failed: user not found", "username", input.Username, "ip", c.ClientIP())
		c.JSON(http.StatusUnauthorized, common.ErrorResponse{Error: "Invalid credentials"})
		return
	}

	match, err := common.ComparePasswordAndHash(input.Password, user.Password)
	if err != nil {
		slog.Error("Failed to compare password hash", "user_id", user.ID, "err", err)
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Error: "internal error"})
		return
	}
	if !match {
		slog.Warn("Login failed: invalid password", "user_id", user.ID, "ip", c.ClientIP())
		c.JSON(http.StatusUnauthorized, common.ErrorResponse{Error: "Invalid credentials"})
		return
	}

	rawToken, accessToken, err := common.IssueTokenPair(h.DB, user.ID, time.Now())
	if err != nil {
		slog.Error("Failed to issue token pair on login", "user_id", user.ID, "err", err)
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Error: "Failed to issue token"})
		return
	}

	c.JSON(http.StatusOK, common.AuthTokenPairResponse{
		RefreshToken: rawToken,
		AccessToken:  accessToken,
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var input common.RefreshInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: err.Error()})
		return
	}
	now := time.Now()

	var result common.HttpErrorResponse
	var newRefreshToken models.RefreshToken
	var newRawToken string

	hashedInput, err := common.HashStringRefreshToken(input.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, common.ErrorResponse{Error: "Invalid refresh token"})
		return
	}

	err = h.DB.Transaction(func(tx *gorm.DB) error {
		var refreshToken models.RefreshToken

		if err := tx.Take(&refreshToken, "token = ? AND revoked = false", hashedInput).Error; err != nil {
			result = common.HttpErrorResponse{Status: http.StatusUnauthorized, Err: "Invalid refresh token"}
			return common.ErrInvalidToken
		}

		if now.After(refreshToken.ExpiresAt) {
			result = common.HttpErrorResponse{Status: http.StatusUnauthorized, Err: "Refresh token expired"}
			return common.ErrExpiredToken
		}

		if err := tx.Model(&refreshToken).Update("revoked", true).Error; err != nil {
			result = common.HttpErrorResponse{Status: http.StatusInternalServerError, Err: "Failed to revoke refresh token"}
			return err
		}

		rawToken, hashedToken, err := common.GenerateRefreshToken()
		if err != nil {
			result = common.HttpErrorResponse{Status: http.StatusInternalServerError, Err: "Failed to generate refresh token"}
			return err
		}
		newRawToken = rawToken

		newRefreshToken = models.RefreshToken{
			UserID:    refreshToken.UserID,
			Token:     hashedToken,
			ExpiresAt: now.Add(time.Hour * 24 * 30),
		}
		if err := tx.Create(&newRefreshToken).Error; err != nil {
			result = common.HttpErrorResponse{Status: http.StatusInternalServerError, Err: "Failed to create new refresh token"}
			return err
		}
		return nil
	})

	if err != nil {
		if result.Status == http.StatusInternalServerError {
			slog.Error("Refresh token transaction failed", "err", err)
		} else {
			slog.Warn("Refresh failed", "reason", result.Err, "ip", c.ClientIP())
		}
		c.JSON(result.Status, common.ErrorResponse{Error: result.Err})
		return
	}

	accessToken, err := common.GenerateAccessToken(newRefreshToken.UserID, now)
	if err != nil {
		slog.Error("Failed to generate token", "user_id", newRefreshToken.UserID, "err", err)
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Error: "Failed to generate access token"})
		return
	}

	c.JSON(http.StatusOK, common.AuthTokenPairResponse{
		RefreshToken: newRawToken,
		AccessToken:  accessToken,
	})
}
