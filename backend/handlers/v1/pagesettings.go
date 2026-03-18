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

type PageSettingsHandler struct {
	DB *gorm.DB
}

func (h PageSettingsHandler) GetSettings(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	var settings common.PageSettingsResponse
	if err := h.DB.Model(&models.PageSettings{}).Take(&settings, "user_id = ?", user.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, common.FromPageSettingsModel(models.GetDefaultPageSettings()))
			return
		}
		slog.Error("Failed to get settings", "user_id", user.ID, "err", err)
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Error: "Failed to get settings"})
		return
	}

	c.JSON(http.StatusOK, settings)
}

func (h PageSettingsHandler) UpdateSettings(c *gin.Context) {
	var input common.PageSettingsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: "Invalid request body"})
		return
	}

	user := c.MustGet("user").(*models.User)
	settings := models.PageSettings{UserID: user.ID}

	if err := h.DB.Where(&settings).Assign(&input).FirstOrCreate(&settings).Error; err != nil {
		slog.Error("Failed to update settings", "user_id", user.ID, "err", err)
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Error: "Failed to update settings"})
		return
	}

	c.JSON(http.StatusOK, common.MessageResponse{Message: "Settings updated"})
}
