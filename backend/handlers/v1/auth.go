package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	common2 "opencw/common"
	"opencw/utils"
	"strings"
	"time"

	"opencw/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB *gorm.DB
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input common2.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, common2.NewErrorResponse(common2.ErrorCodeInvalidRequestBody, "Invalid request body"))
		return
	}

	var user models.User
	err := h.DB.Select("username, email").
		Where("username = ? OR email = ?", input.Username, input.Email).
		First(&user).Error
	if err == nil {
		if user.Username == input.Username && user.Email == input.Email {
			c.JSON(http.StatusConflict, common2.NewErrorResponse(common2.ErrorCodeConflict, "Username and email already exists"))
		} else if user.Username == input.Username {
			c.JSON(http.StatusConflict, common2.NewErrorResponse(common2.ErrorCodeConflict, "Username already exists"))
		} else {
			c.JSON(http.StatusConflict, common2.NewErrorResponse(common2.ErrorCodeConflict, "Email already exists"))
		}
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		slog.Error("Failed to query user", "err", err, "username", input.Username)
		c.JSON(http.StatusInternalServerError, common2.NewErrorResponse(common2.ErrorCodeDatabaseFailure, "Database failure"))
		return
	}

	hash, err := utils.HashPassword(input.Password)
	if err != nil {
		slog.Error("Failed to hash password", "err", err)
		c.JSON(http.StatusInternalServerError, common2.NewErrorResponse(common2.ErrorCodePasswordHashFailed, "Failed to hash password"))
		return
	}

	user = models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: hash,
	}
	if err := h.DB.Create(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c.JSON(http.StatusConflict, common2.NewErrorResponse(common2.ErrorCodeConflict, "Registration conflict, please try again"))
			return
		}

		slog.Error("Failed to create user", "err", err, "username", input.Username)
		c.JSON(http.StatusInternalServerError, common2.NewErrorResponse(common2.ErrorCodeInternalServerError, "Failed to create user"))
		return
	}

	rawToken, accessToken, err := utils.IssueTokenPair(h.DB, user.ID, time.Now())
	if err != nil {
		slog.Error("Failed to issue token pair", "user_id", user.ID, "err", err)
		c.JSON(http.StatusInternalServerError, common2.NewErrorResponse(common2.ErrorCodeTokenIssueFailed, "Failed to issue token, try to login."))
		return
	}

	c.JSON(http.StatusOK, common2.AuthTokenPairResponse{
		RefreshToken: rawToken,
		AccessToken:  accessToken,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input common2.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, common2.NewErrorResponse(common2.ErrorCodeInvalidRequestBody, "Invalid request body"))
		return
	}

	var queryString string
	if strings.Contains(input.Identifier, "@") {
		queryString = "email = ?"
	} else {
		queryString = "username = ?"
	}

	var user models.User
	if err := h.DB.Take(&user, queryString, input.Identifier).Error; err != nil {
		slog.Warn("Login failed", "identifier", input.Identifier, "ip", c.ClientIP())
		c.JSON(http.StatusUnauthorized, common2.NewErrorResponse(common2.ErrorCodeInvalidCredentials, "Invalid credentials"))
		return
	}

	match, err := utils.ComparePasswordAndHash(input.Password, user.Password)
	if err != nil {
		slog.Error("Failed to compare password hash", "user_id", user.ID, "err", err)
		c.JSON(http.StatusInternalServerError, common2.NewErrorResponse(common2.ErrorCodeInternalServerError, "internal error"))
		return
	}
	if !match {
		slog.Warn("Login failed: invalid password", "user_id", user.ID, "ip", c.ClientIP())
		c.JSON(http.StatusUnauthorized, common2.NewErrorResponse(common2.ErrorCodeInvalidCredentials, "Invalid credentials"))
		return
	}

	rawToken, accessToken, err := utils.IssueTokenPair(h.DB, user.ID, time.Now())
	if err != nil {
		slog.Error("Failed to issue token pair on login", "user_id", user.ID, "err", err)
		c.JSON(http.StatusInternalServerError, common2.NewErrorResponse(common2.ErrorCodeTokenIssueFailed, "Failed to issue token"))
		return
	}

	c.JSON(http.StatusOK, common2.AuthTokenPairResponse{
		RefreshToken: rawToken,
		AccessToken:  accessToken,
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var input common2.RefreshInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, common2.NewErrorResponse(common2.ErrorCodeInvalidRequestBody, "Invalid request body"))
		return
	}
	now := time.Now()

	var result common2.HttpErrorResponse
	var newRefreshToken models.RefreshToken
	var newRawToken string

	hashedInput, err := utils.HashStringRefreshToken(input.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, common2.NewErrorResponse(common2.ErrorCodeInvalidToken, "Invalid refresh token"))
		return
	}

	err = h.DB.Transaction(func(tx *gorm.DB) error {
		var refreshToken models.RefreshToken

		if err := tx.Take(&refreshToken, "token = ? AND revoked = false", hashedInput).Error; err != nil {
			result = common2.HttpErrorResponse{Status: http.StatusUnauthorized, Code: common2.ErrorCodeInvalidToken, Err: "Invalid refresh token"}
			return common2.ErrInvalidToken
		}

		if now.After(refreshToken.ExpiresAt) {
			result = common2.HttpErrorResponse{Status: http.StatusUnauthorized, Code: common2.ErrorCodeExpiredToken, Err: "Refresh token expired"}
			return common2.ErrExpiredToken
		}

		if err := tx.Model(&refreshToken).Update("revoked", true).Error; err != nil {
			result = common2.HttpErrorResponse{Status: http.StatusInternalServerError, Code: common2.ErrorCodeInternalServerError, Err: "Failed to revoke refresh token"}
			return err
		}

		rawToken, hashedToken, err := utils.GenerateRefreshToken()
		if err != nil {
			result = common2.HttpErrorResponse{Status: http.StatusInternalServerError, Code: common2.ErrorCodeInternalServerError, Err: "Failed to generate refresh token"}
			return err
		}
		newRawToken = rawToken

		newRefreshToken = models.RefreshToken{
			UserID:    refreshToken.UserID,
			Token:     hashedToken,
			ExpiresAt: now.Add(time.Hour * 24 * 30),
		}
		if err := tx.Create(&newRefreshToken).Error; err != nil {
			result = common2.HttpErrorResponse{Status: http.StatusInternalServerError, Code: common2.ErrorCodeInternalServerError, Err: "Failed to create new refresh token"}
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
		c.JSON(result.Status, common2.NewErrorResponse(result.Code, result.Err))
		return
	}

	accessToken, err := utils.GenerateAccessToken(newRefreshToken.UserID, now)
	if err != nil {
		slog.Error("Failed to generate token", "user_id", newRefreshToken.UserID, "err", err)
		c.JSON(http.StatusInternalServerError, common2.NewErrorResponse(common2.ErrorCodeTokenIssueFailed, "Failed to generate access token"))
		return
	}

	c.JSON(http.StatusOK, common2.AuthTokenPairResponse{
		RefreshToken: newRawToken,
		AccessToken:  accessToken,
	})
}
